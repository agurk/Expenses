#!/usr/bin/perl

use strict;
use warnings;

use Cwd qw(abs_path getcwd);
BEGIN
{
    push (@INC, getcwd());
    no if $] >= 5.018, warnings => "experimental";
}

use Net::DBus;
use Net::DBus::Reactor;

use EventSettings;

use Try::Tiny;
use Switch;

use DataTypes::Expense;

use ExpenseData::Loaders::Loader;
use ExpenseData::Loaders::Loader_AMEX;
use ExpenseData::Loaders::Loader_Aqua;
use ExpenseData::Loaders::Loader_Nationwide;

use Database::DAL;
use Database::ExpenseDB;
use Database::ExpensesDB;
use Database::ClassificationsDB;
use Classifier;

my $classificationDB = ClassificationsDB->new();
my $expenseDB = ExpenseDB->new();
my $expensesDB = ExpensesDB->new();

sub handleMessage
{
	my ($message, $args) = @_;
	switch ($message) {
		case 'CONFIRM_CLASSIFICATION' { $expenseDB->confirmClassification($$args{'eid'}) }
		case 'CHANGE_CLASSIFICATION' { $expenseDB->saveClassification($$args{'eid'}, $$args{'cid'}, 1) }
		case 'SAVE_CLASSIFICATION' { _save_classification($args) }
		case 'CHANGE_AMOUNT' { $expenseDB->saveAmount($$args{'eid'}, $$args{'amount'}); }
		case 'TAG_EXPENSE' { $expenseDB->setTagged($$args{'eid'}, $$args{'tag'}) }
		case 'DUPLICATE_EXPENSE' { $expenseDB->duplicateExpense($$args['eid']); }
		case 'LOAD_RAW' { _load_raw_data($args) }
		case 'CLASSIFY' { _classify_data() }
	}

#
#	if (defined $args)
#	{
#		foreach (keys %$args) {print "$_ ->	",$$args{$_},"\n"}
#	}
}

sub _classify_data
{
    print "Classifying new rows\n";
    my $classifier = Classifier->new(expenseDB=>$expenseDB, expensesDB=>$expensesDB);
    $classifier->processUnclassified();
    return 0;
}


sub _loadAccounts
{
    my ($self) = @_;
    my @loaders;
    foreach (@{$expensesDB->getAccounts()})
    {
        push (@loaders, $_->[0]->new(   numbers_store => $expensesDB,
                                        account_name  => $_->[1],
                                        account_id    => $_->[2],
                                        build_string  => $_->[3]));
    }
    return \@loaders;
}


sub _load_raw_data
{
    my ($args) = @_;
	if ((defined $args) && (%$args))
	{
		print "TODO now load just this alid\n";
	}
	else
	{
	    print "Loading Account data...";
	    my $accounts = _loadAccounts();
	    print "done\n";
	    print "loading expenses data...\n";
	    foreach (@$accounts)
	    {
	        print "    Loading: ",$_->account_name(),'...';
	        try { $_->loadRawInput(); }   catch { print "ERROR: ",$_; };
	        print "done.\n";
	    }
	}
    print "done\n";
}

sub _save_classification
{
    my ($args) = @_; 
    #$classification->setClassificationID($$commands[0]);
    #$classification->setDescription($$commands[1]);
    #$classification->setValidFrom($$commands[2]);
    #$classification->setValidTo($$commands[3]);
    #$classification->setExpense($$commands[4]);
    #$classificationDB->saveClassification($classification);
}


sub main
{

	my $bus=Net::DBus->session();
	my $service=$bus->get_service($DBUS_SERVICE_NAME);
	my $object=$service->get_object($SERVICE_OBJECT_NAME, $DBUS_INTERFACE_NAME);
	
	
	$object->connect_to_signal($EVENT_TYPE, \&handleMessage);
	
	my $reactor=Net::DBus::Reactor->main();
	$reactor->run();
}

main();


