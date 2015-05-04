#!/usr/bin/env perl 
#===============================================================================
#
#         FILE: DocumentDB.pm
#
#  DESCRIPTION: Data Access Layer for Document Object
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
	my $query = 'select r.date, r.filename, r.filesize, r.text, r.textmoddate, r.deleted from Documents r where r.did = ?';
	my $sth = $dbh->prepare($query);
	$sth->execute($documentID);

	my $row = $sth->fetchrow_arrayref();
	my $document = Document->new(	DocumentID=>$documentID,
								ModDate=>$$row[0],
								Filename=>$$row[1],
								FileSize=>$$row[2],
								Text=>$$row[3],
								TextModDate=>$$row[4],
								Deleted=>$$row[5],
								);

	$query = 'select eid from DocumentExpenseMapping where did = ?';
	$sth = $dbh->prepare($query);
	$sth->execute($documentID);

	foreach my $row ( $sth->fetchrow_arrayref())
	{   
		$document->addExpenseID($$row[0]) if ($row);
	}  

	return $document;

}

sub _createNewDocument
{
	my ($self, $document) = @_;
	my $dbh = $self->_openDB();
	my $insertString='insert into documents (date, filename, filesize, text, textmoddate, deleted) values (?, ?, ?, ?, ?, ?)';
	my $sth = $dbh->prepare($insertString);
	$sth->execute($document->getModDate, $document->getFilename, $document->getFileSize, $document->getText, $document->getTextModDate, $document->isDeleted);
	$sth->finish();

	$sth=$dbh->prepare('select max(did) from Documents');
	$sth->execute();
	$document->setDocumentID($sth->fetchrow_arrayref()->[0]);
	$sth->finish();
	$self->_setDocumentExpenses($document);
}

sub _updateDocument
{
	my ($self, $document) = @_;
	my $dbh = $self->_openDB();
	my $query = 'update documents set date = ?, filename = ?, filesize = ?, text=?, textmoddate = ?, deleted = ? where did = ?';
	my $sth = $dbh->prepare($query);
	$sth->execute($document->getModDate, $document->getFilename, $document->getFileSize, $document->getText, $document->getTextModDate, $document->isDeleted, $document->getDocumentID);
	$sth->finish();
	$self->_setDocumentExpenses($document);
}

sub _setDocumentExpenses
{
	my ($self, $document) = @_; 
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('select distinct eid from documentexpensemapping where did = ?');
    $sth->execute($document->getDocumentID);

    my %RIDS;
    while ( my $row = $sth->fetchrow_arrayref())
    {   
        $RIDS{$$row[0]} = 0 if (defined $row);
    }   

    foreach (@{$document->getExpenseIDs()})
    {
        if (exists $RIDS{$_})
        {
            delete $RIDS{$_};
            next;
        }
		# todo add confirmed
        my $insertString='insert into documentexpensemapping (did , eid, confirmed) values (?, ?, 0)';
        my $sth=$dbh->prepare($insertString);
        $sth->execute($document->getDocumentID(), $self->_makeTextQuery($_));
        $sth->finish();
    }   

    foreach (keys %RIDS)
    {   
        my $query = 'delete from documentexpensemapping where did = ? and eid = ?';
        my $sth=$dbh->prepare($query);
        $sth->execute($document->getDocumentID(), $self->_makeTextQuery($_));
        $sth->finish();
    } 
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

