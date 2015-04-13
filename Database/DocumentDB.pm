#!/usr/bin/env perl 
#===============================================================================
#
#         FILE: DocumentDB.pm
#
#  DESCRIPTION: Data Access Layer for Expense Object
#
#      OPTIONS: ---
# REQUIREMENTS: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 08/04/15 17:58
#     REVISION: ---
#===============================================================================

package DocumentDB;
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
use DataTypes::Document;
use Time::Piece;

sub getDocument
{
	my ($self, $documentID) = @_;
	my $dbh = $self->_openDB();
	my $query = 'select r.date, r.filename, r.filesize, r.text, r.textmoddate from Documents r where r.reid = ?';
	my $sth = $dbh->prepare($query);
	$sth->execute($documentID);

	my $row = $sth->fetchrow_arrayref();
	my $document = Document->new(	DocumentID=>$documentID,
								ModDate=>$$row[0],
								Filename=>$$row[1],
								Filesize=>$$row[2],
								Text=>$$row[3],
								TextModDate=>$$row[4],
								);
}

sub _createNewDocument
{
	my ($self, $document) = @_;
	my $dbh = $self->_openDB();
	my $insertString='insert into documents (date, filename, filesize, text, textmoddate) values (?, ?, ?, ?, ?)';
	my $sth = $dbh->prepare($insertString);
	$sth->execute($document->getModDate, $document->getFilename, $document->getFileSize, $document->getText, $document->getTextModDate);
	$sth->finish();

	$sth=$dbh->prepare('select max(reid) from Documents');
	$sth->execute();
	$document->setDocumentID($sth->fetchrow_arrayref()->[0]);
	$sth->finish();
}

sub _updateDocument
{
	my ($self, $document) = @_;
	my $dbh = $self->_openDB();
	my $query = 'update documents set date = ?, filename = ?, filesize = ?, text=?, textmoddate = ? where reid = ?';
	my $sth = $dbh->prepare($query);
	$sth->execute($document->getModDate, $document->getFilename, $document->getFileSize, $document->getText, $document->getTextModDate, $document->getDocumentID);
	$sth->finish();
}

sub saveDocument
{
	my ($self, $document) = @_;
	if ($document->getDocumentID > -1)
	{
		$self->_updateDocument($document);
	}
	else
	{
		$self->_createNewDocument($document);
	}
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

