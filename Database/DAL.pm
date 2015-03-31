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

package DAL;

use strict;
use warnings;
use utf8;

use Moose;

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';
use constant ACCOUNT_LOADERS_TABLE=>'AccountLoaders';
use constant EXPENSE_RAW_MAPPING_TABLE => 'ExpenseRawMapping';
use constant EXPENSE_TAG_TABLE => 'Tagged';

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
	$dbh->do('CREATE TABLE ' . EXPENSE_TAG_TABLE . '(eid INTEGER, tag STRING)');
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

sub _getCurrentDateTime
{
	my $time = gmtime();
	return $time->datetime;
}

1;

