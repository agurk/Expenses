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

package NumbersDB;

use strict;
use warnings;
use utf8;

use DBI;
use Moose;
use DataTypes::Expense;

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';
use constant ACCOUNT_LOADERS_TABLE=>'AccountLoaders';
use constant EXPENSE_RAW_MAPPING_TABLE => 'ExpenseRawMapping';

use constant DSN => 'dbi:SQLite:dbname=/home/timothy/bin/Expenses/expenses.db';

sub create_tables
{
	my ($self) = @_;
	my $dbh = $self->_openDB();
    $dbh->do("DROP TABLE IF EXISTS " . RAW_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . EXPENSES_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . CLASSIFICATION_DEFINITION_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . CLASSIFIED_DATA_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . ACCOUNT_DEFINITION_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . LOADER_DEFINITION_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . PROCESSOR_DEFINITION_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . ACCOUNT_LOADERS_TABLE);
    $dbh->do("DROP TABLE IF EXISTS " . EXPENSE_RAW_MAPPING_TABLE);

    $dbh->do('CREATE TABLE ' . RAW_TABLE . '(rid INTEGER PRIMARY KEY AUTOINCREMENT, rawStr TEXT UNIQUE, importDate DATE, aid INTEGER)');
    $dbh->do('CREATE TABLE ' . EXPENSES_TABLE . '(eid INTEGER PRIMARY KEY AUTOINCREMENT, aid INTEGER, description TEXT, amount REAL, ccy TEXT, amountFX REAL, ccyFX TEXT, fxRate REAL, commission REAL, date DATE)');
    $dbh->do('CREATE TABLE ' . CLASSIFICATION_DEFINITION_TABLE . '(cid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, validFrom DATE, validTo DATE)');
    $dbh->do('CREATE TABLE ' . CLASSIFIED_DATA_TABLE . '(eid INTEGER PRIMARY KEY, cid INTEGER, confirmed INTEGER)');
    $dbh->do('CREATE TABLE ' . ACCOUNT_DEFINITION_TABLE . '(aid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, pid INTEGER, lid INTEGER, ccy TEXT, isExpense INTEGER)');
    $dbh->do('CREATE TABLE ' . LOADER_DEFINITION_TABLE . '(lid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, loader TEXT)');
    $dbh->do('CREATE TABLE ' . PROCESSOR_DEFINITION_TABLE . '(pid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, processor TEXT)');
    $dbh->do('CREATE TABLE ' . ACCOUNT_LOADERS_TABLE . '(alid INTEGER PRIMARY KEY AUTOINCREMENT, aid INTEGER, buildStr TEXT, enabled INTEGER)');
    $dbh->do('CREATE TABLE ' . EXPENSE_RAW_MAPPING_TABLE . '(MID INTEGER PRIMARY KEY AUTOINCREMENT, EID INTEGER, RID INTEGER)');
    $dbh->disconnect();
}

sub _makeTextQuery
{
    my ($self, $text) = @_;
	return 'NULL' unless (defined $text);
    $text =~ s/'/''/g;
    return $text;
}

sub _openDB
{
	my ($self, $arguments) = @_;
    my $dbh = DBI->connect(DSN, '', '', { RaiseError => 1, HandleError => \&_genericDBErrorHandler}) or die $DBI::errstr;
	return $dbh;
}

sub _genericDBErrorHandler
{
    my $error = shift;
	print 'Error in DB operation: ',$error,"\n";
	return 1;
}

sub addRawExpense
{
    my ($self, $rawLine, $account) = @_;
    #my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1, HandleError => \&_handleRawError }) or die $DBI::errstr;
	my $dbh = $self->_openDB();
	$dbh->{HandleError} = \&_handleRawError;

	my $insertString = 'insert into ' . RAW_TABLE . '(rawStr, importDate, aid) values (?, ?, ?)';
    my $sth = $dbh->prepare($insertString);
	my @bindValues;
	$bindValues[0] = $self->_makeTextQuery($rawLine);
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

sub saveExpense
{
    # just dealing with new expenses so far...
    my ($self, $expense) = @_;

	my $dbh = $self->_openDB();

    my $insertString='insert into '.EXPENSES_TABLE.' (aid, description, amount, ccy, amountFX, ccyFX, fxRate, commission, date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)';
    my $sth = $dbh->prepare($insertString);
    $sth->execute($self->_makeTextQuery($expense->getAccountID()),
				  $self->_makeTextQuery($expense->getExpenseDescription()),
				  $expense->getExpenseAmount(),
				  $expense->getCCY(),
				  $expense->getFXAmount(),
				  $expense->getFXCCY(),
				  $expense->getFXRate(),
				  $expense->getCommission(),
				  $expense->getExpenseDate());
    $sth->finish();

	# TODO: make this a bit safer
    $sth=$dbh->prepare('select max(eid) from expenses');
    $sth->execute();
    $expense->setExpenseID($sth->fetchrow_arrayref()->[0]);
    $sth->finish();

	foreach (@{$expense->getRawIDs()})
	{
		my $insertString='insert into '. EXPENSE_RAW_MAPPING_TABLE .' (eid, rid) values (?, ?)';
		$sth=$dbh->prepare($insertString);
	    $sth->execute($self->_makeTextQuery($expense->getExpenseID(), $self->_makeTextQuery($_)));
		$sth->finish();
	}

    my $insertString2='insert into '.CLASSIFIED_DATA_TABLE.' (eid, cid) values (?, ?)';
    $sth = $dbh->prepare($insertString2);
    $sth->execute($self->_makeTextQuery($expense->getExpenseID()), $self->_makeTextQuery($expense->getExpenseClassification()));
    $sth->finish();

    $dbh->disconnect();
}

sub mergeExpenses
{
	my ($self, $primaryExpense, $secondaryExpense) = @_;
	my $dbh = $self->_openDB();
	$dbh->{AutoCommit} = 0;

	eval
	{
		my $sth=$dbh->prepare('select rid from expenserawmapping where eid = ?');
		$sth->execute($secondaryExpense);
		foreach my $row ( $sth->fetchrow_arrayref())
		{
			my $sth2 = $dbh->prepare('insert into expenserawmapping (eid, rid) values(?,?)');
			$sth2->execute($primaryExpense, $row->[0]);
		}
		$sth = $dbh->prepare('delete from expenses where eid = ?');
		$sth->execute($secondaryExpense);
		$sth = $dbh->prepare('delete from expenserawmapping where eid = ?');
		$sth->execute($secondaryExpense);
		$sth = $dbh->prepare('delete from classifications where eid = ?');
		$sth->execute($secondaryExpense);

		$dbh->commit();

	};

    if($@)
	{
		warn "Error inserting the link and tag: $@\n";
		$dbh->rollback();
	}

}

sub getAccounts
{
	my ($self) = @_;
    my @accounts;
	my $dbh = $self->_openDB();

    my $sth = $dbh->prepare('select ldr.loader, a.name, a.aid, l.buildStr from accountdef a, accountloaders l, loaderdef ldr where a.aid = l.aid and a.lid = ldr.lid and l.enabled;');
    $sth->execute();

    while (my @row = $sth->fetchrow_array)
    {
        push (@accounts, \@row);
    }

    $sth->finish();
    
    return \@accounts;
}

1;

