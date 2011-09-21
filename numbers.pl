#!/usr/bin/perl

my $DATA;
my $CATEGORIES;

sub getCount
{
    $count = 0;
    my ($month, $category) = @_;
    foreach (%$DATA)
    {
	if (FOO)
	$count += FOO
    }
    return $count;
}

sub getGraphData
{
    my $month = $_;
    my @months=("Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug","Sep","Nov","Dec");
    for (my $i=0;$i < $month; $i++)
    {
	$index = 0;
	foreach ($CATEGORY)
	{
	    $month[$index] = getCount($i, $_);
	    $index++;
	}
    }
}
