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
use NumbersDB;
use SSWriter;
use Loaders::Loader;
use Loaders::Loader_AMEX;
use Loaders::Loader_Nationwide;
use Loaders::Loader_Aqua;
use Classifier;

use Try::Tiny;

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
	foreach (@{$numbersStore->getAccounts()})
	{
		print 'loading: ',$_->[0],"\n";
		push (@loaders, $_->[0]->new(	numbers_store => $numbersStore,
										settings => $settings,
										account_name  => $_->[1],
										account_id	  => $_->[2],
										build_string  => $_->[3]));
	}
	return \@loaders;
}

sub inital_setup
{
	print "Running inital setup of expenses\n";
    my $settings = Settings->new();
    my $foo = NumbersDB->new(settings=>$settings);
	$foo->create_tables();
	print "setup now complete\n";
}

sub main
{
    my $settings = Settings->new();
    my $foo = NumbersDB->new(settings=>$settings);

    print "Loading Account data...";
    my $accounts = loadAccounts($settings, $foo);
    print "done\n";
    print "loading expenses data...\n";
    foreach (@$accounts)
    {
        print "    Loading: ",$_->account_name(),'...';
        try { $_->loadRawInput(); }   catch { print "ERROR: ",$_; };
        print "done.\n";
    }
    print "done\n";

	print "Classifying new rows\n";
	my $classifier = Classifier->new(numbers_store=>$foo,settings=>$settings);
	$classifier->processUnclassified();

#
#    foreach (@$accounts)
#    {
#        $_->loadNewClassifications();
#    }
#    print "Creating Google Docs Data...";
#    my @results;
#    # Cycle through all the months
#    for (my $i=1; $i<13; $i++)
#    {
#        foreach (SSWriter->createRowMonth($i, $foo->getExpensesByMonth($i)))
#        {
#            push(@results, $_);
#        }
#        foreach (SSWriter->createRowDays_HACK(14+$i, $foo->getExpensesByDay($i)))
#        {
#            push(@results, $_);
#        }
#    }
#    print "data created, writing...";
#    writeSheet(\@results, $settings);
#    print "done\n";
#    #$foo->getExpensesTypeForMonth(2,10);
}

main();
#inital_setup();

