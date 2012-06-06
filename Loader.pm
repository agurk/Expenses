#!/usr/bin/perl

use strict;
use warnings;

# Package to process the raw data from an input file ready to be input into the store numbers

package Loader;
use Moose;

has 'numbers_store' => (is => 'rw', isa => 'Numbers', required => 1);
has 'file_name' => ( is => 'rw', isa => 'Str' );
has 'settings' => ( is => 'rw', required => 1);
has 'input_data' => ( is => 'rw', isa => 'ArrayRef', writer=>'set_input_data', reader=>'get_input_data');
has 'account_name' => (is =>'rw', isa=>'Str');

# Returns 1 if the passed classification is a valid one
sub _checkClassificationValidity
{
    return 1;
    my ($value, $classifications) = shift;
    print "checking Value: $value\n";
    foreach (keys %$classifications) {print $_,"\n"}
    return 1 if (exists $$classifications{$value});
    return 0;
}

sub validateClassification
{
    my ($self, $value) = @_;
    my $classifications = $self->settings->CLASSIFICATIONS();
    return 0 if ($value eq "");
    return 1 if (exists $$classifications{$value});
    return 0;
}

sub getClassification
{
    my ($self, $lineParts) = @_;
    while(1)
    {
	print "Enter classification for: ";
    	print ($lineParts);
    	print "\n    > ";
    	my $value =<>;
    	chomp ($value);
    	if ($self->validateClassification($value))
	{
	    print "    Classified as: ",$self->settings->CLASSIFICATIONS->{$value},"\n";
	    return $value;
	} else {
	    print " **** Invalid classification: $value\n";
	}
    }
}

sub loadInput
{
    my $self = shift;
    if (defined $self->file_name())
    {
	my @input_data;
        open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
        foreach (<$file>)
        {
            push(@input_data, $_);
        }
        close($file);
	$self->input_data = \@input_data;
    }
    else
    {
        $self->set_input_data($self->_pullOnlineData());
    }
}

sub _loadCSVLine
{
    print "If you're seeing this, something went wrong...\n";
}

sub loadNewClassifications
{
    my $self = shift;
    foreach (@{$self->get_input_data})
    {
        $self->_loadCSVLine($_);
    }
    $self->numbers_store()->save();
}

1;


