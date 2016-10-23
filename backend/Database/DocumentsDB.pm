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

package DocumentsDB;
use Moose;
extends 'DAL';

use strict;
use warnings;
use utf8;

use DBI;
use DataTypes::Expense;
use Time::Piece;

sub getAllDocuments
{
	my ($self) = @_; 
	my $dbh = $self->_openDB();

	my $selectString = 'select d.did, d.date, d.filename, d.filesize, d.text, d.textmoddate, d.deleted from documents d';

	my $sth = $dbh->prepare($selectString);
	$sth->execute();

	my @returnArray;
	while (my @row = $sth->fetchrow_array())
	{
		push (@returnArray, \@row);
	}
	#my @returnArray;
	#while (my $row = $sth->fetchrow_arrayref())
	#{
#		push (@returnArray, $row);
#	}

	$sth->finish();
	$dbh->disconnect();

	return \@returnArray;
}

sub getUnclassifiedDocuments
{
	my ($self) = @_; 
	my $dbh = $self->_openDB();

	my $selectString = 'select d.did from documents d where d.deleted = 0 and (d.text = "" or d.text is null)  and d.did not in (select distinct did from documentexpensemapping)';

	my $sth = $dbh->prepare($selectString);
	$sth->execute();

	my @returnArray;
	while (my @row = $sth->fetchrow_array())
	{
		push (@returnArray, $row[0]);
	}

	$sth->finish();
	$dbh->disconnect();

	return \@returnArray;
}

sub confirmDocEx
{
	my ($self, $dmid) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('update documentexpensemapping set confirmed = 1 where dmid = ?');
	$sth->execute($dmid);
	$sth->finish();
	$dbh->disconnect();
}


sub removeDocEx
{
	my ($self, $dmid) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('delete from documentexpensemapping where dmid = ?');
	$sth->execute($dmid);
	$sth->finish();
	$dbh->disconnect();
}

sub findTaggedDocuments
{ 
	my ($self, $fromDate, $toDate, $tag) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare("select distinct(dem.did) from DocumentExpenseMapping dem join tagged t on dem.eid = t.eid where t.tag = ? and datetime(t.modified,'unixepoch') >= strftime(?) and datetime(t.modified,'unixepoch') < strftime(?);");
	$sth->execute($tag, $fromDate, $toDate);

	my @documents;
	while (my @row = $sth->fetchrow_array())
	{
		push (@documents, $row[0]);
	}
	return \@documents;
}

sub isNewDocument
{
	my ($self, $filename, $filesize, $date) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('select count (*) from documents where filename = ? and filesize = ? and date = ?');
	$sth->execute($filename, $filesize, $date);
	my $count = $sth->fetchrow_arrayref()->[0];
	return 0 if ($count > 0);
	return 1;
}

1;

