#!/usr/bin/perl

use strict;
use warnings;

use Cwd qw(abs_path getcwd); 

BEGIN
{
    # paths needed in INC for google writer
    push (@INC, getcwd().'/../'); 
}   


# Set STDOUT as hot
$| = 1;

use NumbersDB;

sub merge_expenses
{
	my ($primaryExpense, $secondaryExpense) = @_;
	exit 1 unless (defined $primaryExpense and ! $primaryExpense eq '');
	exit 1 unless (defined $secondaryExpense and ! $primaryExpense eq '');

    my $foo = NumbersDB->new();

	print "Merging: $secondaryExpense into $primaryExpense\n";
	$foo->mergeExpenses($primaryExpense,$secondaryExpense);
}

merge_expenses($ARGV[0], $ARGV[1]);

