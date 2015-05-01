#!/usr/bin/perl
#
#===============================================================================
#
#         FILE: Loader_Doxie.pm
#
#  DESCRIPTION: Load documents from a Doxie scanner
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 08/04/15 18:54
#     REVISION: ---
#===============================================================================

package Loader_Doxie;
use Moose;
extends 'Loader';

use strict;
use warnings;

require HTTP::Request;
require LWP::UserAgent;

use Database::DAL;
use Database::DocumentDB;
use Database::ExpensesDB;

use DataTypes::Document;

use JSON;

sub loadDocument
{
	my ($self) = @_;
	my $rdb = DocumentDB->new();

	my $request = HTTP::Request->new(GET => 'http://192.168.1.100:8080/scans.json');
	#print $request->as_string,"\n";

	my $ua = LWP::UserAgent->new;
	my $response = $ua->request($request);

	my $scans = decode_json $response->content;
	foreach (@$scans)
	{
		#print keys(%$_),"\n";
		print $_->{name},"\n";
		print $_->{modified},"\n";
		print $_->{size},"\n\n";
		$_->{name} =~ m/([^\/]*)$/;
		my $name = $1;

		chdir ('data/documents');

		if ($rdb->isNewDocument($name, $_->{size}, $_->{modified}))
		{
			$self->_get_image($name);
			
			my $document = Document->new(Filename=>$name,
										ModDate=>$_->{modified},
										FileSize=>$_->{size},);
			$rdb->saveDocument($document);
			print "saved: $document->getFilename\n";
		}
		else
		{
			print "skipping $_->{name}";
		}
	}
}

sub _get_image
{
	my ($self, $name) = @_;
		my $request2 = HTTP::Request->new(GET => 'http://192.168.1.100:8080/scans/DOXIE/JPEG/' . $name);
		print $request2->as_string,"\n";
		my $ua2 = LWP::UserAgent->new;
		my $response2 = $ua2->request($request2);
		my $filename = "$name";
		print "writing to $filename\n";
		open (my $file, '>', $filename) or die "Cannot open $filename for writing\n";
		print $file $response2->content;
		close ($file);
}

1;

