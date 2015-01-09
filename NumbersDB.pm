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
#      VERSION: 0.1
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

has 'settings' => (is => 'rw', required => 1);

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';
use constant ACCOUNT_LOADERS_TABLE=>'AccountLoaders';
use constant EXPENSE_RAW_MAPPING_TABLE => 'ExpenseRawMapping';

sub create_tables
{
    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1 }) or die $DBI::errstr;
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
    $dbh->do('CREATE TABLE ' . CLASSIFIED_DATA_TABLE . '(eid INTEGER PRIMARY KEY, cid INTEGER)');
    $dbh->do('CREATE TABLE ' . ACCOUNT_DEFINITION_TABLE . '(aid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, pid INTEGER, lid INTEGER, ccy TEXT, isExpense INTEGER)');
    $dbh->do('CREATE TABLE ' . LOADER_DEFINITION_TABLE . '(lid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, loader TEXT)');
    $dbh->do('CREATE TABLE ' . PROCESSOR_DEFINITION_TABLE . '(pid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, processor TEXT)');
    $dbh->do('CREATE TABLE ' . ACCOUNT_LOADERS_TABLE . '(alid INTEGER PRIMARY KEY AUTOINCREMENT, aid INTEGER, buildStr TEXT, enabled INTEGER)');
    $dbh->do('CREATE TABLE ' . EXPENSE_RAW_MAPPING_TABLE . '(MID INTEGER PRIMARY KEY AUTOINCREMENT, EID INTEGER, RID INTEGER)');
    $dbh->disconnect();
}

sub _cleanQueryLine
{
    my ($self, $line) = @_;
    $line =~ s/'/''/g;
    return $line;
}

sub _makeTextQuery
{
    my ($self, $text) = @_;
	return 'NULL' unless (defined $text);
    $text = $self->_cleanQueryLine($text);
    return '\'' . $text . '\'';
}

sub addRawExpense
{
    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my ($self, $rawLine, $account) = @_;
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1, HandleError => \&_handleRawError }) or die $DBI::errstr;

    

    my $insertString = 'insert into ' . RAW_TABLE . '(rawStr, importDate, aid) values (\'' 
                            . $self->_cleanQueryLine($rawLine) . '\',\'' . gmtime() . '\',\'' . $account . '\')' ;

    my $sth = $dbh->prepare($insertString);
    $sth->execute();

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
    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my ($self, $rawLine, $account) = @_; 
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1}) or die $DBI::errstr;

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
    my %classifications;
    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1}) or die $DBI::errstr;

    my $sth = $dbh->prepare('select cid,name from ClassificationDef');
    $sth->execute();


    while (my $row = $sth->fetchrow_arrayref)
    {
        $classifications{$$row[0]} = $$row[1];
    }

    $sth->finish();
    

#$classifications{'1'} = 'ONE';
    return \%classifications;
}

sub _makeSaveNewExpenseQuery
{
    my ($self, $expense) =@_;
    my $insertString='insert into '.EXPENSES_TABLE.' (aid, description, amount, ccy, amountFX, ccyFX, fxRate, commission, date) values (';
    $insertString .= $self->_makeTextQuery($expense->getAccountID());
    $insertString .= ',' . $self->_makeTextQuery($expense->getExpenseDescription());
    $insertString .= ',' . $self->_makeTextQuery($expense->getExpenseAmount());
    $insertString .= ',' . $self->_makeTextQuery($expense->getCCY());
    $insertString .= ',' . $self->_makeTextQuery($expense->getFXAmount());
    $insertString .= ',' . $self->_makeTextQuery($expense->getFXCCY());
    $insertString .= ',' . $self->_makeTextQuery($expense->getFXRate());
    $insertString .= ',' . $self->_makeTextQuery($expense->getCommission());
    $insertString .= ',' . $self->_makeTextQuery($expense->getExpenseDate());
	$insertString .= ')';
    return $insertString;
}

sub _makeSaveNewRawProcessedMappingsQuery
{
	my ($self, $expense, $rid) = @_;
    my $insertString='insert into '. EXPENSE_RAW_MAPPING_TABLE .' (eid, rid) values (';
    $insertString .= $self->_makeTextQuery($expense->getExpenseID());
	$insertString .= ',' . $self->_makeTextQuery($rid);
	$insertString .= ')';
    return $insertString;
}

sub _makeSaveNewClassificationQuery
{
    my ($self, $expense) =@_;
    my $insertString='insert into '.CLASSIFIED_DATA_TABLE.' (eid, cid) values (';
    $insertString .= $self->_makeTextQuery($expense->getExpenseID()) . ',';
    $insertString .= $self->_makeTextQuery($expense->getExpenseClassification()) . ')';
    return $insertString;
}

sub saveExpense
{
    # just dealing with new expenses so far...
    my ($self, $expense) = @_;

    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1}) or die $DBI::errstr;

    my $sth = $dbh->prepare($self->_makeSaveNewExpenseQuery($expense));
    $sth->execute();
    $sth->finish();

	# TODO: make this a bit safer
    $sth=$dbh->prepare('select max(eid) from expenses');
    $sth->execute();
    $expense->setExpenseID($sth->fetchrow_arrayref()->[0]);
    $sth->finish();

	foreach (@{$expense->getRawIDs()})
	{
		$sth=$dbh->prepare($self->_makeSaveNewRawProcessedMappingsQuery($expense, $_));
	    $sth->execute();
		$sth->finish();
	}

    $sth = $dbh->prepare($self->_makeSaveNewClassificationQuery($expense));
    $sth->execute();
    $sth->finish();

    $dbh->disconnect();
}

sub getAccounts
{
    my @accounts;
    my $dsn = 'dbi:SQLite:dbname=expenses.db';
    my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1}) or die $DBI::errstr;

    my $sth = $dbh->prepare('select ldr.loader, a.name, a.aid, l.buildStr from accountdef a, accountloaders l, loaderdef ldr where a.aid = l.aid and a.lid = ldr.lid and l.enabled <> 0;');
    $sth->execute();


    while (my @row = $sth->fetchrow_array)
    {
        push (@accounts, \@row);
    }

    $sth->finish();
    
    return \@accounts;
    
}

1;

