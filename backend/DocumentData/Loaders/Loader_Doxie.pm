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
#extends 'Loader';

use strict;
use warnings;

require HTTP::Request;
require LWP::UserAgent;

use Database::DAL;
use Database::DocumentDB;
use Database::DocumentsDB;
use Database::ExpensesDB;

use DataTypes::Document;

has 'Address' =>	( isa => 'Str',
						is => 'rw',
						reader => 'getAddress',
						writer => 'setAddress',
						default => 'http://192.168.1.100:8080',
					);

has 'Password' =>	( isa => 'Str',
					  is  => 'rw',
					  reader => 'getPassword',
					  writer => 'setPassword',
					  default => '',
					); 

has 'User' =>	( isa => 'Str',
					  is  => 'rw',
					  reader => 'getUser',
					  writer => 'setUser',
					  default => 'doxie',
				); 

has 'DocDir' => ( isa => 'Str', is=>'ro', reader=>'getDocDir', required=>1 );

use JSON;

sub loadDocument
{
	my ($self) = @_;
	print "loading...\n";
	my $rdb = DocumentDB->new();
    my $ddb = DocumentsDB->new();

	my $request = HTTP::Request->new(GET => $self->getAddress .'/scans.json');
	$request->authorization_basic($self->getUser, $self->getPassword);
#	print $request->as_string,"\n";

	my $ua = LWP::UserAgent->new;
	my $response = $ua->request($request);

	if ($response->content =~ m/^Can't connect to/)
	{
		print "Cannot connect to device\n";
		return;
	}
	my $scans = decode_json $response->content;
	foreach (@$scans)
	{
		print "\n\n";
        print $_->{name},"\n";
		print $_->{modified},"\n";
		print $_->{size},"\n";
		$_->{name} =~ m/([^\/]*)$/;
		my $name = $1;

		chdir ($self->getDocDir);

		if ($ddb->isNewDocument($name, $_->{size}, $_->{modified}))
		{
			$self->_get_image($name);
			
			my $document = Document->new( Filename=>$name,
										  ModDate=>$_->{modified},
										  FileSize=>$_->{size},);
			$rdb->saveDocument($document);
			print "---> saved: ",$document->getFilename,"\n";
		}
		else
		{
			print "---> skipping: ",$_->{name},"\n";
		}
	}

    print "\nReloading Failed Documents\n";
    my $failedDocs = $self->getFailedDocuments($ddb);
    foreach my $filename (@$failedDocs)
    {
        unlink $filename or warn "Could not unlink $filename: $!";
        $self->_get_image($filename);
    }
    print "done\n";
}

sub getFailedDocuments
{
	my ($self, $ddb) = @_;
	unless (defined $ddb) { $ddb = DocumentsDB->new(); }
    my $docs = $ddb->getAllDocuments();
    # Allow for small differences, like rotated documents, etc
    my $threshold = 5000;
    my @failedDocs;
    foreach my $document (@$docs)
    {
        my $size = -s $$document[2] ;
        my $diff = $size - $$document[3];
        if (abs($diff) > $threshold)
        {
            print 'Size mismatch on ',$$document[2],' (',$diff,'b)',"\n";
            push(@failedDocs, $$document[2]);
        }
    }
    return \@failedDocs;
}

sub _get_image
{
	my ($self, $name) = @_;
	my $request = HTTP::Request->new(GET => $self->getAddress.'/scans/DOXIE/JPEG/' . $name);
	$request->authorization_basic($self->getUser, $self->getPassword);
	my $ua = LWP::UserAgent->new;
	my $response = $ua->request($request);
	my $filename = "$name";
	print "writing to $filename\n";
	open (my $file, '>', $filename) or die "Cannot open $filename for writing\n";
	print $file $response->content;
	close ($file);
}

1;

