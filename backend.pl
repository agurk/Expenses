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
use EventGenerator;

use Try::Tiny;
use Switch;

use DataTypes::Expense;

use ExpenseData::Loaders::Loader;
use ExpenseData::Loaders::Loader_AMEX;
use ExpenseData::Loaders::Loader_Aqua;
use ExpenseData::Loaders::Loader_Nationwide;

use Engines::ExpensesEngine;

use Database::DAL;
use Database::ExpenseDB;
use Database::ExpensesDB;
use Database::ClassificationsDB;
use Classifier;

use POSIX qw/strftime/;

my $ExpensesEngine = ExpensesEngine->new();


sub printMessage 
{
	my ($message, $args) = @_;
	my $now = strftime "%Y-%m-%d %H:%M:%S", localtime;
	print '-' x 100, "\n";
	print $now,", received: $message";
	if (keys %$args)
	{
		print ":\n";
		foreach (keys %$args) {print ' ' x 32,"$_ ->	",$$args{$_},"\n"}
	}
	else
	{
		print "\n";
	}
	print '=' x 100, "\n";
}

sub handleMessage
{
	my ($message, $args) = @_;
	printMessage($message, $args);
	switch ($message) {
		case 'CHANGE_AMOUNT' { $ExpensesEngine->change_amount($args) }
		case 'CHANGE_CLASSIFICATION' { $ExpensesEngine->change_classification($args) }
		case 'CLASSIFY' { $ExpensesEngine->classify_data() }
		case 'CONFIRM_CLASSIFICATION' { $ExpensesEngine->confirm_classification($args) }
		case 'DUPLICATE_EXPENSE' { $ExpensesEngine->duplicate_expense($args) }
		case 'LOAD_RAW' { $ExpensesEngine->load_raw_data($args) }
		case 'MERGE_EXPENSE' { $ExpensesEngine->merge_expense($args) }
		case 'MERGE_EXPENSE_COMMISSION' { $ExpensesEngine->merge_expense_commission($args) }
		case 'REPROCESS_EXPENSE' { $ExpensesEngine->reprocess_expense($args)  }
		case 'SAVE_CLASSIFICATION' { $ExpensesEngine->save_classification($args) }
		case 'SAVE_EXPENSE' { $ExpensesEngine->save_expense($args) }
		case 'TAG_EXPENSE' { $ExpensesEngine->tag_expense($args) }
		case 'SAVE_ACCOUNT' { $ExpensesEngine->save_account($args) }
	}
}

sub main
{
	my $pid = fork();
	EventGenerator::runGenerator(0) if ($pid == 0);
	sleep 2;
	unless ($pid == 0)
	{
		my $bus=Net::DBus->session();
		my $service=$bus->get_service($DBUS_SERVICE_NAME);
		my $object=$service->get_object($SERVICE_OBJECT_NAME, $DBUS_INTERFACE_NAME);

		$object->connect_to_signal($EVENT_TYPE, \&handleMessage);
		
		my $reactor=Net::DBus::Reactor->main();
		$reactor->run();
	}
}

main();


