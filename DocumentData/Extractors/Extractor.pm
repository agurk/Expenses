#!/usr/bin/perl
#
#===============================================================================
#
#         FILE: Classifier.pm
#
#  DESCRIPTION: Class to extract documents
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 29/07/15 17:00
#     REVISION: ---
#===============================================================================

package Extractor;
use Moose;

use strict;
use warnings;

use File::Copy;

use Database::DAL;
use Database::DocumentDB;
use Database::DocumentsDB;
use Database::ExpensesDB;

use DataTypes::Document;

sub extractDocument
{
	my ($self, $document, $exportFolder) =@_;
	print "coping ",$document->getFilename,' to ',$exportFolder,"\n";
	copy('data/documents/'.$document->getFilename,$exportFolder);
}

sub extractDocuments
{
	my ($self, $fromDate, $toDate, $tag, $exportFolder) = @_;
	my $ddb = DocumentsDB->new();
	my $docdb = DocumentDB->new();
	foreach (@{$ddb->findTaggedDocuments($fromDate, $toDate, $tag)})
	{
		my $document = $docdb->getDocument($_);
		$self->extractDocument($document, $exportFolder);
	}
}
1;

