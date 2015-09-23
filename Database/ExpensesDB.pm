#!/usr/bin/env perl 
#===============================================================================
#
#         FILE: NumbersDB.pm
#
#  DESCRIPTION: Data Access Layer between DB and program
#
#      OPTIONS: ---
# REQUIREMENTS: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.2
#      CREATED: 23/12/14 11:19:12
#     REVISION: ---
#===============================================================================

package ExpensesDB;
use Moose;
extends 'DAL';

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';
use constant ACCOUNT_LOADERS_TABLE=>'AccountLoaders';
use constant EXPENSE_RAW_MAPPING_TABLE => 'ExpenseRawMapping';

use strict;
use warnings;
use utf8;

use DBI;
use DataTypes::Expense;
use Time::Piece;

sub addRawExpense
{
	my ($self, $rawLine, $account) = @_;
	my $dbh = $self->_openDB();
	$dbh->{HandleError} = \&_handleRawError;

	my $insertString = 'insert into ' . RAW_TABLE . '(rawStr, importDate, aid) values (?, ?, ?)';
	my $sth = $dbh->prepare($insertString);
	my @bindValues;
	$bindValues[0] = $rawLine;
	$bindValues[1] = gmtime();
	$bindValues[2] = $account;
	$sth->execute(@bindValues);

	$dbh->disconnect();
}

sub _handleRawError
{
	my $error = shift;
	unless ($error =~ m/UNIQUE constraint failed: RawData.rawStr/)
	{
		print 'Error performing raw insert: ',$error,"\n";
	}
	return 1;
}

sub getUnclassifiedLines
{
	my ($self, $rawLine, $account) = @_; 
	my $dbh = $self->_openDB();

# TODO: this what if there is no matching account?
	my $selectString = 'select processor,rawstr,rid,rawdata.aid,ccy  from rawdata,accountdef,processordef where rid not in (select distinct rid from expenserawmapping) and rawdata.aid = accountdef.aid and accountdef.pid=processordef.pid';

	my $sth = $dbh->prepare($selectString);
	$sth->execute();

	my @returnArray;
	while (my @row = $sth->fetchrow_array())
	{
		push (@returnArray, \@row);
	}

	$sth->finish();
	$dbh->disconnect();

	return \@returnArray;
}

sub getCurrentClassifications
{
	my ($self) = @_;
	my %classifications;
	my $dbh = $self->_openDB();

	my $sth = $dbh->prepare('select cid,name from ClassificationDef');
	$sth->execute();


	while (my $row = $sth->fetchrow_arrayref)
	{
		$classifications{$$row[0]} = $$row[1];
	}

	$sth->finish();

	return \%classifications;
}


sub getValidClassifications
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare("select cid from classificationdef where date(validfrom) <= date(?) and (validto = '' or date(validto) >= date(?))");
    $sth->execute($expense->getDate(), $expense->getDate());

	my @results;
	while (my $row = $sth->fetchrow_arrayref) {push (@results, $$row[0])}
	$sth->finish();
	return \@results;
}

sub getClassificationStats
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare("select cid, count (*) from expenses e, classifications c where date(e.date) > date( ?, 'start of month','-12 months') and e.eid = c.eid group by cid");
    $sth->execute($expense->getDate());

	my @results;
	while (my @row = $sth->fetchrow_array) {push (@results, \@row)}
	$sth->finish();
	return \@results;
}

sub getExactMatches
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('select cid, count (*) from expenses e, classifications c where e.description = ? and e.eid = c.eid group by cid');
    $sth->execute($expense->getDescription());

	my @results;
	while (my @row = $sth->fetchrow_array) {push (@results, \@row)}
	$sth->finish();
	return \@results;
}

sub getAccounts
{
	my ($self, $alid) = @_;
    my @accounts;
	my $dbh = $self->_openDB();
	my $sth;

	if ($alid)
	{	
		$sth = $dbh->prepare('select ldr.loader, a.name, a.aid, l.buildStr from accountdef a, accountloaders l, loaderdef ldr where a.aid = l.aid and a.lid = ldr.lid and alid=?;');
		$sth->execute($alid);
	}
	else
	{	
		$sth = $dbh->prepare('select ldr.loader, a.name, a.aid, l.buildStr from accountdef a, accountloaders l, loaderdef ldr where a.aid = l.aid and a.lid = ldr.lid and l.enabled;');
		$sth->execute();
	}

    while (my @row = $sth->fetchrow_array)
    {
        push (@accounts, \@row);
    }

    $sth->finish();
    
    return \@accounts;
}

sub getDateMatches
{
	my ($self, $date) = @_;
    my @matches;
	my $dbh = $self->_openDB();

    my $sth = $dbh->prepare('select e.eid, e.description, e.amount, e.amountfx from expenses e where e.date = ?');
    $sth->execute($date);

    while (my @row = $sth->fetchrow_array)
    {
        push (@matches, \@row);
    }

    $sth->finish();
    
    return \@matches;
}

sub getRawLine
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	
    my $sth = $dbh->prepare('select rawstr from rawdata where rid = ?');
	my $rawID = $expense->getRawIDs->[0];
    $sth->execute($rawID);

	my $row = $sth->fetchrow_arrayref();
	return $$row[0];
}

sub getNWCashFees
{
	my ($self) = @_;
	my $dbh = $self->_openDB();
	
    my $sth = $dbh->prepare('select eid from expenses where description = \'Non-Sterling cash fee\'');
	#my $rawID = $expense->getRawIDs->[0];
    #$sth->execute($rawID);
}

sub saveExpenseDocumentMappings
{
	my ($self, $eid, $did) = @_;
	my $dbh = $self->_openDB();
    my $sth = $dbh->prepare('select dmid, did from documentexpensemapping where eid = ?');
    $sth->execute($eid);

    while (my @row = $sth->fetchrow_array)
    {
		if (exists $did->{$row[1]})
		{
			# confirm ?
			delete $did->{$row[1]}
		} else {
			$sth = $dbh->prepare('delete from documentexpensemapping where dmid = ?');
			$sth->execute($row[0]);
		}
    }

	foreach (keys %$did)
	{
		$sth = $dbh->prepare('insert into documentexpensemapping(did, eid, confirmed) values(?, ?, 1)');
		$sth->execute($_, $eid);
	}

    $sth->finish();
    
}

sub saveAccount
{
	my ($self, $aid, $name, $ccy, $lid, $pid) = @_;
	if ($aid eq 'NEW')
	{
		my $dbh = $self->_openDB();
		my $sth = $dbh->prepare('insert into accountdef (name, ccy, lid, pid) values (?, ?, ?, ?)');
		$sth->execute($name, $ccy, $lid, $pid);
		$sth->finish();
	}
}


1;

