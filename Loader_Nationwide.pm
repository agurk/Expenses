#!/usr/bin/perl

package Loader_Nationwide;
use Moose;

extends 'Loader';

# Date _after_ which new style CSV is used
# format of DD MMM YYYY
has 'changeover_date' => ( is=> 'rw', isa => 'Str' );

# The nationwide CSV files have five liens at the top that shouln't be processed
# but we'll do a nice check rather than just ignoring the top five lines!
sub _skipLine
{
    my $self = shift;
    my $line = shift;
    return 1 if ($line eq '');
    return 1 if ($line eq "\n");
    return 1 if ($line eq "\r");
    return 1 if ($line =~ m/^\"Account Name/);
    return 1 if ($line =~ m/^\"Account Balance/);
    return 1 if ($line =~ m/^\"Available Balance/);
    return 1 if ($line =~ m/^\"Date\",\"Transaction type\"/);
    #return 1 if ($line =~ m/^Date,Transactions,Debits/);
    return 0;
}

# Return true if line should be skipped as predates new format
# CSV (and we're assuming it has been loaded alredy)
sub _beforeChangeOver
{
    my $self = shift;
    return 0 unless (defined $self->changeover_date);
    my $date = shift;
    $date =~ s/"//g;
    my @currentDate = split(/ /,$date);
    my @changeoverDate = split(/ /,$self->changeover_date);
    # Don't skip if current year is after changeover year
    return 0 if ($currentDate[2] > $changeoverDate[2]);
    # Skip if year is before (so we know from now onwards in this 
    # that the years will be the same
    return 1 if ($currentDate[2] < $changeoverDate[2]);
    my %months  = ('Jan',1,'Feb',2,'Mar',3,'Apr',4,'May',5,'Jun',6,
                   'Jul',7,'Aug',8,'Sep',9,'Oct',10,'Nov',11,'Dec',12);
    # If in the same year the month is before the change over, then skip
    # otherwise it's good to go
    return 0 if ($months{$currentDate[1]} > $months{$changeoverDate[1]});
    return 1 if ($months{$currentDate[1]} < $months{$changeoverDate[1]});
    # if in the same month, 
    return 0 if ($currentDate[0] > $changeoverDate[0]);
    return 1;
}

# File takes the csv format of:
# OLD_FORMAT: date,name,debit,credit,balance
# NEW_FORMAT: date,transaction type,name,debit,credit,balance
# New style CSV file adds the transaction type field, which we will
# concatenate to our ID field - the date must be used to
#know when to stop using the old format
sub load
{
    my $self = shift;
    my $DATA = $self->numbers_store()->data_list();
    open(my $file,"<",$self->file_name()) or warn "No file exists: ",$self->file_name(),"\n";
    foreach my $line (<$file>)
    {
#	print $line;
	chomp($line);
	# Also removing carriage return, as CSV has windows style
	# line breaks
	$line =~ s/\r//g;
	next if ($self->numbers_store()->isDupe($line));
	next if ($self->_skipLine($line));
	my @lineParts=split(/,/, $line);
	# skip if no debit - this is not an expense!
	next if ($lineParts[3] eq " ");
	next if ($lineParts[3] eq "");
	# could do with a proper date object here...
	next if ($self->_beforeChangeOver($lineParts[0]));
	# Strip leading char - £ sign specifically
	$lineParts[3] =~ s/^[^0123456789\.]*//;
	my $classification = $self->getClassification($line);
	$lineParts[0] =~ s/\"//g;
	$lineParts[3] =~ s/\"//g;
	my @record = ($lineParts[1].$lineParts[2],$lineParts[0],$lineParts[3],$classification);
	#push (@$DATA, \@record);
	#$$DATA{$line} = \@record;
	$self->numbers_store()->addValue($line,\@record);
    }
    close($file);
    $self->numbers_store()->save();
}

# File takes the csv format of:
# date,name,debit,credit,balance
sub _load_old
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

1;

