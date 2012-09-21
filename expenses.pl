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
use Loader_Aqua;

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

sub loadAccounts
{
    my ($settings, $numbersStore) = @_;
    my @loaders;
    open (my $file, '<', $settings->ACCOUNT_FILE) or die "Can't load accounts file\n";
    foreach(<$file>)
    {
	next if ($_ =~ m/^#/);
	chomp;
	my @lineParts = split (/,/, $_);
	if ($lineParts[0] eq 'aqua')
	{
	    push(@loaders, Loader_Aqua->new(numbers_store => $numbersStore,
					    account_name=>$lineParts[1],
					    file_name=>$lineParts[2],
					    settings=>$settings));
	} elsif ($lineParts[0] eq 'amex')
	{
	    push (@loaders, Loader_AMEX->new(numbers_store => $numbersStore, 
					     account_name=>$lineParts[1],
	                    		     file_name=>$lineParts[2],
			    		      AMEX_CARD_NUMBER=>$lineParts[3],
			    		      AMEX_USERNAME=>$lineParts[4],
			    		      AMEX_PASSWORD=>$lineParts[5],
			    		      AMEX_INDEX=>$lineParts[6],
			    		      settings=>$settings));
	
	} elsif ($lineParts[0] eq 'nationwide')
	{
	    push (@loaders, Loader_Nationwide->new(numbers_store => $numbersStore,
					            account_name=>$lineParts[1],
						    file_name=>$lineParts[2],
						    NATIONWIDE_ACCOUNT_NUMBER=>$lineParts[3],
						    NATIONWIDE_ACCOUNT_NAME=>$lineParts[4],
						    NATIONWIDE_MEMORABLE_DATA=>$lineParts[5],
						    NATIONWIDE_SECRET_NUMBERS=>$lineParts[6],
						    settings=>$settings));
	}	
	
    }
    close($file);
    return \@loaders;
}

sub main
{
    my $settings = Settings->new();
    my $foo = Numbers->new(data_file_name => $settings->DATAFILE_NAME, settings=>$settings);
    print "Loading Account data...";
    my $accounts = loadAccounts($settings, $foo);
    print "done\n";
    print "loading expenses data...";
    foreach (@$accounts)
    {
	$_->loadInput();
	print 'done: ',$_->account_name(),'...';
    }
    print "done\n";
    foreach (@$accounts)
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

