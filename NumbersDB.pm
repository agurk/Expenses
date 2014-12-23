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

has 'settings' => (is => 'rw', required => 1);

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant CLASSIFIED_DATA_TABLE=>'ClassifiedData';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDefinition';

sub create_tables
{
	my $dbh = shift;
	$dbh->do("DROP TABLE IF EXISTS " . RAW_TABLE);
	$dbh->do("DROP TABLE IF EXISTS " . EXPENSES_TABLE);
	$dbh->do("DROP TABLE IF EXISTS " . CLASSIFICATION_DEFINITION_TABLE);
	$dbh->do("DROP TABLE IF EXISTS " . CLASSIFIED_DATA_TABLE);
	$dbh->do("DROP TABLE IF EXISTS " . ACCOUNT_DEFINITION_TABLE);

	$dbh->do('CREATE TABLE ' . RAW_TABLE . '(rid INTEGER PRIMARY KEY AUTOINCREMENT, rawStr TEXT UNIQUE, importDate DATE, aid INTEGER)');
	$dbh->do('CREATE TABLE ' . EXPENSES_TABLE . '(eid INTEGER PRIMARY KEY, rid INTEGER, description TEXT, amount REAL)');
	$dbh->do('CREATE TABLE ' . CLASSIFICATION_DEFINITION_TABLE . '(cid INTEGER PRIMARY KEY, name TEXT, validFrom DATE, validTo DATE)');
	$dbh->do('CREATE TABLE ' . CLASSIFIED_DATA_TABLE . '(eid INTEGER, cid aoeu)');
	$dbh->do('CREATE TABLE ' . ACCOUNT_DEFINITION_TABLE . '(aid INTEGER PRIMARY KEY, name TEXT)');
}

sub _cleanRawLine
{
	my ($self, $line) = @_;
	$line =~ s/'/''/g;
	return $line;
}

sub addRawExpense
{
	my $dsn = 'dbi:SQLite:dbname=expenses.db';
	my ($self, $rawLine, $account) = @_;
	my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1, HandleError => \&_handleRawError }) or die $DBI::errstr;

	

	my $insertString = 'insert into ' . RAW_TABLE . '(rawStr, importDate, aid) values (\'' 
							. $self->_cleanRawLine($rawLine) . '\',\'' . gmtime() . '\',\'' . $account . '\')' ;

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

sub main
{
	my $dsn = 'dbi:SQLite:dbname=expenses.db';
	my $dbh = DBI->connect($dsn, '', '', { RaiseError => 1 }) or die $DBI::errstr;
	create_tables($dbh);
}

main();

1;



