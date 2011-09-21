#!/usr/bin/perl

package Loader_Nationwide;
use Moose;

extends 'Loader';

# The nationwide CSV files have five liens at the top that shouln't be processed
# but we'll do a nice check rather than just ignoring the top five lines!
sub _skipLine
{
    my $self = shift;
    my $line = shift;
    return 1 if ($line eq '');
    return 1 if ($line =~ m/^Account name/);
    return 1 if ($line =~ m/^Account balance/);
    return 1 if ($line =~ m/^Available balance/);
    return 1 if ($line =~ m/^Date,Transactions,Debits/);
    return 0;
}

# File takes the csv format of:
# date,name,debit,credit,balance
sub load
{
    my $self = shift;
    my $DATA = $self->numbers_store()->data_list();
    open(my $file,"<",$self->file_name()) or warn "No file exists: ",$self->file_name(),"\n";
    foreach(<$file>)
    {
	chomp();
	next if ($self->numbers_store()->isDupe($_));
	next if ($self->_skipLine($_));
	my @lineParts=split(/,/, $_);
	# skip if no debit - this is not an expense!
	next if ($lineParts[2] eq "");
	# Strip leading char - £ sign specifically
	$lineParts[2] =~ s/^[^0123456789\.]*//;
	my $classification = $self->getClassification($_);
	my @record = ($lineParts[1],$lineParts[0],$lineParts[2],$classification);
	#push (@$DATA, \@record);
	#$$DATA{$_} = \@record;
	$self->numbers_store()->addValue($_,\@record);
    }
    close($file);
    $self->numbers_store()->save();
}
