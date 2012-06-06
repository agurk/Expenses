#!/usr/bin/perl

use strict;
use warnings;

use Cwd qw(abs_path getcwd); 

BEGIN
{
    # paths needed in INC for google writer
    push (@INC, '/home/timothy/perl');
    push (@INC, getcwd()); 
}   

# Set STDOUT as hot
$| = 1;

use Settings;
use Numbers;
use SSWriter;
use Loader;
use Loader_AMEX;
use Loader_Nationwide;

use Getopt::Std;
my %OPTIONS;
getopts("an",\%OPTIONS);

sub writeSheet
{
    my ($results, $settings) = @_;
    my $writer = SSWriter->new(
	user_name => $settings->GOOGLE_DOCS_USERNAME,
	password=> $settings->GOOGLE_DOCS_PASSWORD,
        workbook => $settings->GOOGLE_DOCS_WORKBOOK,
	worksheet => $settings->GOOGLE_DOCS_WORKSHEET);
    $writer->write_to_sheet($results);
}

sub main
{
    my $settings = Settings->new();
    my %classifications;
    my $foo = Numbers->new(data_file_name => $settings->DATAFILE_NAME, settings=>$settings);
    print "Loading Account data...";
    my @accounts;
    # no need to save these as these methods do a save after loading
    if ($OPTIONS{'a'})
    {
        push (@accounts, Loader_AMEX->new(numbers_store => $foo, 
#                  file_name=>'in/amex.csv',
	                  settings=>$settings,
		          classifications=>\%classifications,
			    account_name=>'AMEX'));
    }
    if ($OPTIONS{'n'})
    {
	push (@accounts, Loader_Nationwide->new(numbers_store => $foo,
#			       file_name=>'in/debit.csv',
			       settings=>$settings,
			       classifications=>\%classifications,
				account_name=>'Nationwide'));
    }
    print "done\n";
    print "loading expenses data...";
    foreach (@accounts)
    {
	$_->loadInput();
	print 'done: ',$_->account_name(),'...';
    }
    print "done\n";
    foreach (@accounts)
    {
	$_->loadNewClassifications();
    }
    print "Creating Google Docs Data...";
    my @results;
    # Cycle through all the months
    for (my $i=1; $i<13; $i++)
    {
	foreach (SSWriter->createRowMonth($i, $foo->getExpensesByMonth($i)))
        {
	    push(@results, $_);
        }
        foreach (SSWriter->createRowDays_HACK(14+$i, $foo->getExpensesByDay($i)))
	{
            push(@results, $_);
	}
    }
    print "data created, writing...";
    writeSheet(\@results, $settings);
    print "done\n";
    #$foo->getExpensesTypeForMonth(9,10);
}

main();

