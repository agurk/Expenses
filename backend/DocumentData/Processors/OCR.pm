#!/usr/bin/perl
#
#===============================================================================
#
#         FILE: Classifier.pm
#
#  DESCRIPTION: Class to manage the classification of new expense items
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

use utf8;
use Encode;

package OCR;
use Moose;

use strict;
use warnings;

use File::Copy qw(copy);

use DataTypes::Document;

has WorkingDir => ( is=> 'ro', isa=> 'Str', required=>1,reader=>'getWorkingDir');

sub ocr
{
	my ($self, $path, $document) = @_;
	$self->_prepare_image($path, $document->getFilename());
	$self->_run_ocr();
	$self->_save_text($document);
}

sub _prepare_image
{
	my ($self, $path, $filename) = @_;
	copy($path .'/'. $filename, $self->getWorkingDir.'/input_image');
}

sub _run_ocr
{
	my ($self) = @_;
	chdir ($self->getWorkingDir());
    my @command = ('tesseract', 'input_image', 'text_out');
    system(@command) == 0 or warn "Cannot complete OCR\n";
}

sub _save_text
{
	my ($self, $document) = @_;
    my $text = ''; 
    open (my $file, '<',$self->getWorkingDir() . '/text_out.txt')
		or die "Error loading ocr output text\n";
    foreach (<$file>)
    {   
        $text .= $_; 
    }   
    close ($file);
    $document->setText($text);
}

1;

