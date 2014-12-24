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
use Expense;

has 'settings' => (is => 'rw', required => 1);

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';

sub create_tables
{
	my $dbh = shift;
#	$dbh->do("DROP TABLE IF EXISTS " . RAW_TABLE);
#	$dbh->do("DROP TABLE IF EXISTS " . EXPENSES_TABLE);
	$dbh->do("DROP TABLE IF EXISTS " . CLASSIFICATION_DEFINITION_TABLE);
#	$dbh->do("DROP TABLE IF EXISTS " . CLASSIFIED_DATA_TABLE);
#	$dbh->do("DROP TABLE IF EXISTS " . ACCOUNT_DEFINITION_TABLE);
#	$dbh->do("DROP TABLE IF EXISTS " . LOADER_DEFINITION_TABLE);
#	$dbh->do("DROP TABLE IF EXISTS " . PROCESSOR_DEFINITION_TABLE);

#	$dbh->do('CREATE TABLE ' . RAW_TABLE . '(rid INTEGER PRIMARY KEY AUTOINCREMENT, rawStr TEXT UNIQUE, importDate DATE, aid INTEGER)');
#	$dbh->do('CREATE TABLE ' . EXPENSES_TABLE . '(eid INTEGER PRIMARY KEY AUTOINCREMENT, rid INTEGER, aid INTEGER, description TEXT, amount REAL, date DATE)');
	$dbh->do('CREATE TABLE ' . CLASSIFICATION_DEFINITION_TABLE . '(cid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, validFrom DATE, validTo DATE)');
#	$dbh->do('CREATE TABLE ' . CLASSIFIED_DATA_TABLE . '(eid INTEGER PRIMARY KEY, cid INTEGER)');
#	$dbh->do('CREATE TABLE ' . ACCOUNT_DEFINITION_TABLE . '(aid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, pid INTEGER, lid INTEGER)');
#	$dbh->do('CREATE TABLE ' . LOADER_DEFINITION_TABLE . '(lid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, loader TEXT)');
#	$dbh->do('CREATE TABLE ' . PROCESSOR_DEFINITION_TABLE . '(pid INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, processor TEXT)');
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

	print $insertString,"\n";

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
	my $selectString = 'select processor,rawstr,rid,rawdata.aid  from rawdata,accountdef,processordef where rid not in (select distinct rid from expenses) and rawdata.aid = accountdef.aid and accountdef.pid=processordef.pid';

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
	my $insertString='insert into '.EXPENSES_TABLE.' (rid, aid, description, amount, date) values (';
	$insertString .= $self->_makeTextQuery($expense->getRawID()) . ',';
	$insertString .= $self->_makeTextQuery($expense->getAccountID()) . ',';
	$insertString .= $self->_makeTextQuery($expense->getExpenseDescription()) . ',';
	$insertString .= $self->_makeTextQuery($expense->getExpenseAmount()) . ',';
	$insertString .= $self->_makeTextQuery($expense->getExpenseDate()) . ')';
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

	$sth=$dbh->prepare('select max(eid) from expenses');
    $sth->execute();
	$expense->setExpenseID($sth->fetchrow_arrayref()->[0]);
    $sth->finish();

	$sth = $dbh->prepare($self->_makeSaveNewClassificationQuery($expense));
    $sth->execute();
    $sth->finish();

    $dbh->disconnect();
}





sub main
{
	my $dsn = 'dbi:SQLite:dbname=expenses.db';
	my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1 }) or die $DBI::errstr;
	create_tables($dbh);
}




#main();

1;

