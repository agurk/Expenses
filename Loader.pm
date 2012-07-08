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
    my ($self, $record) = @_;
    while(1)
    {
	print "Enter classification for: \n";
    	print '  -- ',$$record[0],"\n  -- ",$$record[1],"  --  £",$$record[2];
    	print "\n  > ";
    	my $value =<>;
    	chomp ($value);
    	if ($self->validateClassification($value))
	{
	    print "Classified as: ",$self->settings->CLASSIFICATIONS->{$value},"\n\n";
	    push (@$record, $value);
	    return 1;
	} elsif ($value eq 'CHANGE VALUE') {
	    my $continue = 1;
	    while($continue)
	    {
		print "\nEnter new amount:\n  > ";
		$value =<>;
		chomp $value;
		if ($value =~ m/^[0-9.]*$/)
		{
		    $$record[2] = $value;
		    print "\n\n";
		    $continue = 0;
		} else {
		    print "**** >$value< is an invalid amount\n";
		}
	    }
	} else {
	    print "**** Invalid classification: $value\n\n";
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

sub _skipLine
{
    return 0;
}

sub _makeRecord{}

sub _loadCSVLine
{
    my ($self, $line) = @_;
    chomp($line);
    $line =~ s/\r//g;
    return if ($self->numbers_store()->isDupe($line));
    return if ($self->_skipLine($line));
    my @lineParts=split(/,/, $line);
    my $record = $self->_makeRecord(\@lineParts);
    $self->getClassification($record);
    $self->numbers_store()->addValue($line,$record);
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


