#!/usr/bin/perl

use strict;
use warnings;

# Package to process the raw data from an input file ready to be input into the store numbers
#
# Basic Structure is to:
# 1) Create a new instance of the loader (of a specific type)
# 2) Run loadInput() to get the frest data
# 3) run loadNewClassifications() to find and classify any new records

package Loader;
use Moose;
use Expense;

has 'numbers_store' => (is => 'rw', required => 1);
has 'file_name' => ( is => 'rw', isa => 'Str', writer => 'setFileName' );
#has 'settings' => ( is => 'rw', required => 1);
has 'account_name' => (is =>'rw', isa=>'Str');
has 'account_id' => (is =>'rw', isa=>'Str');
#has 'data_year' => (is => 'ro', isa=>'Str');
has 'build_string' => (is => 'ro', isa=>'Str', required => 1);

use constant LOAD_ATTEMPT_LIMIT => 3;
use constant INTERNAL_FIELD_SEPARATOR => '|';

# Methods to be overwritten for subclasses
# To create the standard array to pass into the numbers store
sub _makeRecord{print "NULL make record\n"}
sub _useInputLine { return 1; }
sub _ignoreYear { return 0; }
sub _pullOnlineData{ return 0; }
sub _skipLine{ return 0; }
sub _processInputLine { my ($self, $line) = @_; return $line; }

sub _loadCSVRows
{
	my ($self) = @_;
	my @lines;
    open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
    foreach (<$file>)
	{
		chomp;
		push(@lines, $_);
    }
    close($file);
	return \@lines;
}

sub loadRawInput
{
	my $self = shift;
	my @lines;
    if (defined $self->file_name())
    {
		my $results = $self->_loadCSVRows();
		@lines = @$results;
	}
    else
    {
        my $results = $self->_pullOnlineData();
		@lines = @$results;
    }

	foreach (@lines)
	{
		chomp;chop;
		$self->numbers_store()->addRawExpense($_,$self->account_id()) if ($self->_useInputLine($_));
	}
}

sub loadInput
{
    my $self = shift;
    if (defined $self->file_name())
    {
        my @input_data;
		$self->set_input_data(\@input_data);
        open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
        foreach (<$file>)
        {
            push(@input_data, $self->_processInputLine($_)) if ($self->_useInputLine($_));
        }
        close($file);
    }
    else
    {
        my $attempts = 0;
        my $success = 0;
        while ($attempts < LOAD_ATTEMPT_LIMIT)
        {
            # pullOnlineData to return 0 if it fails, as standard
            if ($self->_pullOnlineData())
            {
                # bump up the attempt count to break the loop
                $attempts = LOAD_ATTEMPT_LIMIT;
                $success = 1;
            }
            $attempts++;
        }

        unless ($success)
        {
            print " couldn't load: ",$self->account_name(),' ';
            # Empty array, so if we call loadNewClassifications we won't try and 
            # do things on an empty array
            my @emptyArray;
            $self->set_input_data(\@emptyArray);
        }
    }
}

# Shouldn't need to change this per loader -- FINAl
sub _loadCSVLine
{
    my ($self, $line) = @_;
    chomp($line);
    $line =~ s/\r//g;
    return if ($self->numbers_store()->isDupe($line));
    return if ($self->_skipLine($line));
    my $expenseRecord = $self->_makeRecord(\$line);
#    print    $expenseRecord->getDescription, ' - ',
#            $expenseRecord->getDate, ' - ',
#            $expenseRecord->getAmount, ' - ',
#            $expenseRecord->getOriginalLine, "\n";
    return if ($self->_ignoreYear($expenseRecord));
    $self->getClassification($expenseRecord);
    $self->numbers_store()->addValue($expenseRecord);
}

#sub loadNewClassifications
#{
#    my $self = shift;
#    foreach (@{$self->get_input_data})
#    {
#        $self->_loadCSVLine($_);
#    }
#    $self->numbers_store()->save();
#}

1;

