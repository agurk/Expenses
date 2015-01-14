#!/usr/bin/perl

use strict;
use warnings;

use v5.10;

#no warnings 'experimental::smartmatch';
no if $] >= 5.018, warnings => "experimental";

package ExpensesBackend;
use Moose;

has 'numbers' => ( is => 'rw', isa => 'Object');
has 'settings' => ( is => 'rw', isa => 'Object');

use Cwd qw(abs_path getcwd); 
BEGIN
{
    push (@INC, getcwd()); 
no if $] >= 5.018, warnings => "experimental";
}   

# Set STDOUT as hot
$| = 1;

use Settings;
use NumbersDB;
use Loaders::Loader;
use Loaders::Loader_AMEX;
use Loaders::Loader_Nationwide;
use Loaders::Loader_Aqua;
use Classifier;

use IO::Socket;
use Try::Tiny;
use Switch;

use Getopt::Std;
my %OPTIONS;
getopts("an",\%OPTIONS);

sub _loadAccounts
{
    my ($self) = @_;
    my @loaders;
	foreach (@{$self->numbers->getAccounts()})
	{
		push (@loaders, $_->[0]->new(	numbers_store => $self->numbers,
										settings => $self->settings,
										account_name  => $_->[1],
										account_id	  => $_->[2],
										build_string  => $_->[3]));
	}
	return \@loaders;
}

sub inital_setup
{
	my ($self) = @_;
	print "Running inital setup of expenses\n";
    my $settings = Settings->new();
    my $foo = NumbersDB->new(settings=>$settings);
	$foo->create_tables();
	print "setup now complete\n";
}

sub merge_expenses
{
	my ($self, $primaryExpense, $secondaryExpense) = @_;
	return 1 unless (defined $primaryExpense and ! $primaryExpense eq '');
	return 1 unless (defined $secondaryExpense and ! $primaryExpense eq '');

	print "Merging: $secondaryExpense into $primaryExpense\n";
	$self->numbers->mergeExpenses($primaryExpense,$secondaryExpense);
	return 0;
}

sub classify_data
{
	my ($self) = @_;
	print "Classifying new rows\n";
	my $classifier = Classifier->new(numbers_store=>$self->numbers,settings=>$self->settings);
	$classifier->processUnclassified();
	return 0;
}

sub load_raw_data
{
	my ($self) = @_;
    print "Loading Account data...";
    my $accounts = $self->_loadAccounts($self->settings, $self->numbers);
    print "done\n";
    print "loading expenses data...\n";
    foreach (@$accounts)
    {
        print "    Loading: ",$_->account_name(),'...';
        try { $_->loadRawInput(); }   catch { print "ERROR: ",$_; };
        print "done.\n";
    }
    print "done\n";
}

sub main
{
    my $settings = Settings->new();
    my $numbers = NumbersDB->new(settings=>$settings);
	my $expensesBackend = ExpensesBackend->new(settings => $settings, numbers=>$numbers);	

	print "Server started. Opening Connections\n";

	while (1)
	{

		my $sock = new IO::Socket::INET ( LocalHost => '127.0.0.1',
										  LocalPort => '7070',
										  Proto => 'tcp',
										  Listen => 1,
										  Reuse => 1,
									    ); die "Could not create socket: $!\n" unless $sock;

		my $incoming = $sock->accept();

		while(<$incoming>)
		{
			print "Received command >$_< ";
			switch($_)
			{
				case 'load_raw' {$expensesBackend->load_raw_data();}
				case 'classify' {$expensesBackend->classify_data();}
				else			{print "!!unknown command"}
			}
			print "\n";
		}

		close($sock);
	}
}

main();

