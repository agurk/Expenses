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

package Processor;
use Moose;

use strict;
use warnings;

require HTTP::Request;
require LWP::UserAgent;

use Database::DAL;
use Database::DocumentDB;
use Database::ExpensesDB;

use DataTypes::Document;

sub processDocument
{
	my ($self, $did) = @_;
	my $rdb = DocumentDB->new();
	my $document = $rdb->getDocument($did);
	$self->_ocr_image('data/documents', $document);
	$self->_classify_document($document);
	$rdb->saveDocument($document);
}

sub reclassifyDocument
{
	my ($self, $did) = @_;
	my $rdb = DocumentDB->new();
	my $document = $rdb->getDocument($did);
	$self->_classify_document($document);
	$rdb->saveDocument($document);
}

sub _classify_document
{
	my ($self, $document) = @_;
	$document->removeAllExpenseIDs();
	foreach ( @{$self->_find_matches($document)} )
	{
		$document->addExpenseID($_);
	}
}

sub _ocr_image
{
	my ($self, $path, $document) = @_;
	chdir $path;
	#chdir ('data/documents');
	my @command = ('tesseract', $document->getFilename, $document->getFilename);
	system(@command) == 0 or warn "Cannot complete OCR for " . $document->getFilename . "\n";
	my $text = '';
	open (my $file, '<',$document->getFilename . '.txt');
	foreach (<$file>)
	{
		$text .= $_;
	}
	close ($file);
	$document->setText($text);
}

sub _date_in_bounds
{
	my ($self, $year, $month, $day) = @_;
	return 0 if ($year > 2015 or $year < 2000);
	return 0 if ($month > 12 or $month < 1);
	return 0 if ($day > 31 or $day < 1);
	return 1;
}

sub _add_dates
{
	my ($self, $rawDates, $dates, $pattern) = @_;
	my @pArray = split(//, $pattern);
	my ($dayPos, $monthPos, $yearPos);
	for (my $i = 0; $i<3; $i++)
	{
		$dayPos = $i if ($pArray[$i] eq 'D');
		$monthPos = $i if ($pArray[$i] eq 'M');
		$yearPos = $i if ($pArray[$i] eq 'Y');
	}
	my $max = scalar (@$rawDates);
	for (my $i = 0; $i < $max; $i += 3)
	{
		my $day = $$rawDates[$dayPos];
		my $month = $$rawDates[$monthPos];
		my $year = $$rawDates[$yearPos];
		shift @$rawDates; shift @$rawDates; shift @$rawDates;
		next if $day > 31;
		next if $month > 12;
		$year = '20' . $year if (length $year == 2);
		$month = '0' . $month if (length $month == 1);
		$day = '0' . $day if (length $day == 1);
		next unless $self->_date_in_bounds($year, $month, $day);
		$dates->{"$year-$month-$day"} = 1;
	}
}

sub _find_potential_date_matches
{
	my ($self, $document) = @_;
	my $year  = '(2?0?[0-9]{2})';
	my $month = '(0?[0-9]|1?[0-2])';
	my $day   = '([12][0-9]|3[01]|0?[0-9])';
	
	my @rawDates;
	my %dates;
	my $txt = $document->getText;
	utf8::decode($txt);
	push (@rawDates, ($txt =~ m/$year[-–—\/\\.]$month[-–—\/\\.]$day/g ));
	$self->_add_dates(\@rawDates, \%dates, 'YMD');
	push (@rawDates, ($txt =~ m/$day[-–—\/\\.]$month[-–—\/\\.]$year/g ));
	$self->_add_dates(\@rawDates, \%dates, 'DMY');
	push (@rawDates, ($txt =~ m/$day[^0-9A-Za-z\n]$month[^0-9A-Za-z\n]$year/g ));
	$self->_add_dates(\@rawDates, \%dates, 'DMY');
#	@rawDates=();
	push (@rawDates, ($txt =~ m/$month[-–—\/\\.]$day[-–—\/\\.]$year/g ));
	$self->_add_dates(\@rawDates, \%dates, 'MDY');
#	push (@dates, @{$self->_add_dates(\@rawDates, 'DMY')});
#	push (@dates, ($document->getText =~ m//g ));
#	push (@dates, ($document->getText =~ m/[0-9]{2}[-\/.][0-9]{2}[-\/.]2?0?[0-9]{2}/g));
#	push (@dates, ($document->getText =~ m/([0-3]?[0-9][^0-9A-Za-z\n]{1,3}[01]?[0-9][^0-9A-Za-z\n]{1,3}2?0?[0-9]{2})/g));
#	push (@dates, ($document->getText =~ m/([0-9]{4}[^0-9A-Za-z\n][0-9]{2}[^0-9A-Za-z\n][0-9]{2})/g));
#	push (@dates, ($document->getText =~ m/(2014-06—27)/g));
	foreach (keys %dates)
	{
		print " --> date match $_\n";
	}
	return \%dates;
}

sub _clean_amount
{
	my ($self, $amount, $fxamount) = @_;
	$amount = $fxamount if ($fxamount);
	$amount =~ s/^-//;
	$amount .= '.00' if ($amount =~ m/^[0-9]*$/);
	my $placeholder = '@@_@@_@';
	$amount =~ s/[\.,]/$placeholder/g;
	$amount =~ s/$placeholder/[,.]?/g;
	return $amount;
}

sub _find_matches
{
	my ($self, $document) = @_;
	my $edb = ExpensesDB->new();
	my $dates = $self->_find_potential_date_matches($document);

	my %matches;
	my $maxScore = 0;
	foreach (keys %$dates)
	{
		foreach (@{$edb->getDateMatches($_)})
		{
			my ($eid, $description, $amount, $fxamount) = @$_;
			my $score = 0;
			$amount = $self->_clean_amount($amount, $fxamount);
			print "$description = ";
			$description =~ m/^([^ ]*) *([^ ]*)/;
			my ($one, $two) = ($1, $2);
			$score++ if ($document->getText =~ m/\Q$one/i);
			print "$score ";
			$score++ if ($document->getText =~ m/\Q$two/i);
			print "$score ";
			$score++ if ($document->getText =~ m/$amount/);
			print "$score ";

			$matches{$score} .= ",$eid" if ($score);
			$maxScore = $score if ($score > $maxScore);
			print "$one - $two = $amount -> $score\n";
		}
	}

	my @results;	

	if ($maxScore)
	{
		push (@results, split (/,/, $matches{$maxScore}));
		# we added a , to every ID
		shift @results;
	}

	return \@results;
}

1;

