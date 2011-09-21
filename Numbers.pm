#!/usr/bin/perl

# Module to represent the data store for the input transactions
#
# Data File takes the following format:
# 0: Input String (acts as key)
# 1: Description of transaction
# 2: Date
# 3: Value
# 4: Classification

package Numbers;
use Moose;

has 'data_list' => ( is => 'rw', isa =>'HashRef', default => sub{{}});
has 'data_file_name' => ( is => 'rw', isa => 'Str', required => 1);

# These are the positions in the data list array for the above values
# They key isn't listed here as that is the key for the HASH of which
# this array is the payload
use constant ITEM_DESCRIPTON => 0;
use constant ITEM_DATE=> 1;
use constant ITEM_AMOUNT => 2;
use constant ITEM_CLASSIFICATION => 3;

sub BUILD
{
    my $self = shift;
    $self->loadData();
}

# Data type
# key is input line
sub loadData
{
    my $self=shift;
    my $DATA = $self->data_list();
    if (open(my $file,"<",$self->data_file_name()))
    {
        foreach(<$file>)
        {
    	chomp();
    	my @lineParts = split (/\|/, $_);
	#$lineParts[3] =~ s/ //g;
    	my $key = shift (@lineParts);
    	$$DATA{$key} = \@lineParts;
        }
        close($file);
    }
    else
    {
	print "No file found, empty store created\n";
    }
}

# Add a new value to the existing numbers store
#
sub addValue
{
    my $self = shift;
    my ($key, $ref) = @_;
    #print "adding $key\n";
    my $DATA = $self->data_list();
    if (exists $$DATA{$key})
    {
        warn "Cannot add, as collision on key: $key\n";
	return 0;
    }
    else {$$DATA{$key} = $ref;}
    return 1;
}

sub save
{
    my $self = shift;
    open(my $file,">",$self->data_file_name());
    my $DATA = $self->data_list();
    foreach(keys %$DATA)
    {	
	print $file "$_|";
	my $array = $$DATA{$_};
	foreach(@$array)
	{
	    print $file "$_|";
	}
	print $file "\n";
    }
    close($file);
}

sub isDupe
{
    my $self = shift;
    my $line=shift;
    #print "Dupe checker: $line\n";
    my $DATA=$self->data_list();
    chomp($line);
    return 1 if (exists $$DATA{$line});
#return 1 if ($$DUPELIST{$line}==1);
    return 0;
}

# Returns a numberical version of the month from the date
# In the header
sub _getItemMonth
{
    my $monthIn = shift;
    my %months  = ('Jan',1,'Feb',2,'Mar',3,'Apr',4,'May',5,'Jun',6,
		   'Jul',7,'Aug',8,'Sep',9,'Oct',10,'Nov',11,'Dec',12);
    return $1 if ($monthIn =~ m/[0-9]{2}\/([0-9]{2})\/[0-9]{4}/);
    my @lineParts = split(/ /, $monthIn);
    return $months{$lineParts[1]};
}

# The results have the classification number as the index
# so 0 is going to be 0, as there is no classification for it
# Takes month in as a number
sub getExpensesByMonth
{
    my ($self, $month) = @_;
    my @results = (0,0,0,0,0,0,0,0,0,0,0);
    my $DATA = $self->data_list();
    foreach (keys %$DATA)
    {
#	print $$DATA{$_}->[ITEM_CLASSIFICATION],'---',$_,"\n";
	next if ($$DATA{$_}->[ITEM_CLASSIFICATION] == -1);
	next unless (_getItemMonth($$DATA{$_}->[ITEM_DATE]) == $month);
        $results[$$DATA{$_}->[ITEM_CLASSIFICATION]] += $$DATA{$_}->[ITEM_AMOUNT];
    }
    return \@results;
}

sub _processExpensesDay
{
    my $day = shift;
    $day =~ s/^0//;
    $day =~ s/ //g;
    return $day;
}

sub getExpensesByDay
{
    my $self = shift;
    my $month= shift;
    my @months = (0,'Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec');
    # 32 as zero indexed (and we want a day 0 at 0 here ;) 
    my @results = (0) x 32;
    my $DATA = $self->data_list();
    foreach my $key (keys %$DATA)
    {
	next if ($$DATA{$key}->[ITEM_CLASSIFICATION] == -1);
        if ($$DATA{$key}->[ITEM_DATE] =~ m/([0-9]{2})\/([0-9]{2})\/[0-9]{4}/)
        {
                next unless ($2 == $month);
		$results[_processExpensesDay($1)] += $$DATA{$key}->[ITEM_AMOUNT]
        } else {
                my @lineParts = split(/ /, $$DATA{$key}->[ITEM_DATE]);
                next unless ($lineParts[ITEM_DATE] eq $months[$month]);
		$results[_processExpensesDay($lineParts[0])] += $$DATA{$key}->[ITEM_AMOUNT];
        }
    }
    # Now lets do the summing so the numbers are cumulative
    for (my $i = 1; $i < 32; $i++)  { $results[$i] += $results[$i-1]; }
    return \@results;
}

sub getExpensesTypeForMonth
{
    my ($self, $month, $classification) = @_;
    my $DATA = $self->data_list();
    foreach (keys %$DATA)
    {
	next unless ($$DATA{$_}->[ITEM_CLASSIFICATION] == $classification);
	next unless (_getItemMonth($$DATA{$_}->[ITEM_DATE]) == $month);
	print $_,"\n";
    }
}

1;

