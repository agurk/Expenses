#!/usr/bin/perl

use utf8;
use Encode;

use Moose;
use strict;
use warnings;

sub _date_in_bounds
{
	my ($year, $month, $day) = @_;
	# TODO set as current year
	return 0 if ($year > 2023 or $year < 2000);
	return 0 if ($month > 12 or $month < 1);
	return 0 if ($day > 31 or $day < 1);
	return 1;
}

sub _add_dates
{
	my ($rawDates, $dates, $pattern) = @_;
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
		next unless _date_in_bounds($year, $month, $day);
		$dates->{"$year-$month-$day"} = 1;
	}
}

sub _find_potential_date_matches
{
	my ($txt) = @_;
	my $year  = '(2?0?[0-9]{2})';
	my $month = '(0?[0-9]|1?[0-2])';
	my $day   = '([12][0-9]|3[01]|0?[0-9])';
	my @rawDates;
	my %dates;
	utf8::decode($txt);
	push (@rawDates, ($txt =~ m/$year[-–—\/\\.]$month[-–—\/\\.]$day/g ));
	_add_dates(\@rawDates, \%dates, 'YMD');
	push (@rawDates, ($txt =~ m/$day[-–—\/\\.]$month[-–—\/\\.]$year/g ));
	_add_dates(\@rawDates, \%dates, 'DMY');
	push (@rawDates, ($txt =~ m/(?=$day[^0-9A-Za-z\n] ?$month[^0-9A-Za-z\n] ?$year)/g ));
	_add_dates(\@rawDates, \%dates, 'DMY');
	push (@rawDates, ($txt =~ m/$month[-–—\/\\.]$day[-–—\/\\.]$year/g ));
	_add_dates(\@rawDates, \%dates, 'MDY');
#	push (@dates, @{$self->_add_dates(\@rawDates, 'DMY')});
#	push (@dates, ($document->getText =~ m//g ));
#	push (@dates, ($document->getText =~ m/[0-9]{2}[-\/.][0-9]{2}[-\/.]2?0?[0-9]{2}/g));
#	push (@dates, ($document->getText =~ m/([0-3]?[0-9][^0-9A-Za-z\n]{1,3}[01]?[0-9][^0-9A-Za-z\n]{1,3}2?0?[0-9]{2})/g));
#	push (@dates, ($document->getText =~ m/([0-9]{4}[^0-9A-Za-z\n][0-9]{2}[^0-9A-Za-z\n][0-9]{2})/g));
#	push (@dates, ($document->getText =~ m/(2014-06—27)/g));
	foreach (keys %dates)
	{
		print "$_\n";
	}
	return \%dates;
}

sub main
{
    foreach my $line ( <STDIN> )
    {
        _find_potential_date_matches($line);
    }
}

main();

