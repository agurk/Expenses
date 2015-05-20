#!/usr/bin/perl
#
#===============================================================================
#
#         FILE: Loader_File.pm
#
#  DESCRIPTION: Load documents from a file or directory
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 20/05/15 21:10
#     REVISION: ---
#===============================================================================

package Loader_File;
use Moose;
extends 'Loader';

use strict;
use warnings;

use File::Basename;
use File::Copy;
use Cwd;

use Database::DAL;
use Database::DocumentDB;
use Database::ExpensesDB;

sub loadDocument
{
	my ($self, $path) = @_;
	print "loading documents from $path...\n";
	my $rdb = DocumentDB->new();

	if ( -f $path )
	{
		$self->_newDocument($path, $rdb);
	}
	elsif ( -d $path )
	{
		print "found dir\n";
		chdir ($path);
		opendir (my $dh, $path) or die "Cannot open $path\n";
		while (readdir $dh)
		{
			$self->_newDocument($_, $rdb) if ( -f $_ );
		}
	}
	else
	{
		print "Cannot process $path as neither directory or file\n";
	}
}

sub _newDocument
{
	my ($self, $documentPath, $rdb) = @_;
	print "importing $documentPath\n";
	my ($filename, $path, $suffix) = basename ($documentPath);
	$path .= '';
	chdir ($path) unless ($path eq '');
    my  ($dev,$ino,$mode,$nlink,$uid,$gid,$rdev,$size,
           $atime,$mtime,$ctime,$blksize,$blocks)
               = stat($filename);
	my $document = Document->new( Filename=>$documentPath,
								  ModDate=>$mtime,
								  FileSize=>$size,);
	$rdb->saveDocument($document);
	copy ($filename, $self->getDocumentDir);
}

