#!/usr/bin/perl

use strict;
use warnings;

package ExpensesEngine;
use Moose;

use Cwd qw(abs_path getcwd);
BEGIN
{
    push (@INC, getcwd());
    no if $] >= 5.018, warnings => "experimental";
}

use Try::Tiny;
use Switch;

use DataTypes::Expense;

use ExpenseData::Loaders::Loader;
use ExpenseData::Loaders::Loader_AMEX;
use ExpenseData::Loaders::Loader_Aqua;
use ExpenseData::Loaders::Loader_Danske;
use ExpenseData::Loaders::Loader_Nationwide;

use Database::DAL;
use Database::ExpenseDB;
use Database::ExpensesDB;
use Database::ClassificationsDB;
use Classifier;

my $classificationDB = ClassificationsDB->new();
my $expenseDB = ExpenseDB->new();
my $expensesDB = ExpensesDB->new();

sub merge_expense
{
	my ($self, $args) = @_;
    my $mainEx = $$args{'eid'};
    my $subEx = $$args{'eid_merged'};
	if ($mainEx and $subEx and !($mainEx eq $subEx))
	{
		print "Merging $subEx into $mainEx\n";
		$expenseDB->mergeExpenses($mainEx, $subEx);
	}
}

sub merge_expense_commission
{
	my ($self, $args) = @_;
    my $mainEx = $$args{'eid'};
    my $subEx = $$args{'eid_merged'};
	if ($mainEx and $subEx and !($mainEx eq $subEx))
	{
		print "Merging $subEx into $mainEx as commission\n";
		$expenseDB->mergeExpenseAsCommission($mainEx, $subEx);
	}
}

sub reprocess_expense
{
	my ($self, $args) = @_;
    my $eid = $$args{'eid'};

    my $expense = $expenseDB->getExpense($eid);
    my $rawLines = $expensesDB->getRawLines($eid);

    foreach my $line (@$rawLines)
    {
        $line->[0]->reprocess($expense, $$line[1]);
    }

	$expenseDB->saveExpense($expense);
}

sub classify_data
{
    print "Classifying new rows\n";
    my $classifier = Classifier->new(expenseDB=>$expenseDB, expensesDB=>$expensesDB);
    $classifier->processUnclassified();
    return 0;
}

sub loadAccounts
	{
    my ($self, $args) = @_;
    my @loaders;
	my $alid = '';
	if ((defined $args) && (%$args) && ($$args{'alid'}))
	{
		$alid = $$args{'alid'};
	}
	foreach (@{$expensesDB->getAccounts($alid)})
	{
		push (@loaders, $_->[0]->new(   numbers_store => $expensesDB,
			                            account_name  => $_->[1],
				                        account_id    => $_->[2],
					                    build_string  => $_->[3]));
	}
    return \@loaders;
}


sub load_raw_data
{
    my ($self, $args) = @_;
	print "Loading Account data...";
	my $accounts = $self->loadAccounts($args);
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

sub save_expense
{
    my ($self, $args) = @_; 
	my ($eid, $amount, $description, $date, $classification, $fxAmount, $fxCCY, $fxRate, $commission, $rawDids, $aid, $ccy) =
	($$args{'eid'}, $$args{'amount'}, $$args{'description'}, $$args{'date'}, $$args{'classification'}, $$args{'fxAmount'}, $$args{'fxCCY'}, $$args{'fxRate'}, $$args{'commission'}, $$args{'documents'}, $$args{'aid'}, $$args{'ccy'});
	$fxAmount='' if ($fxAmount eq 'None');
	$fxCCY='' if ($fxCCY eq 'None');
	$fxRate='' if ($fxRate eq 'None');
	$commission='' if ($commission eq 'None');
	my $expense;
	if ($eid eq 'NEW')
	{
		$expense = Expense->new	(
									AccountID => $aid,
									Amount => $amount,
									Description => $description,
									Date => $date,
									Currency => $ccy,
								);
	} else {
		$expense = $expenseDB->getExpense($eid);
		$expense->setAmount($amount);
		$expense->setDescription($description);
		$expense->setDate($date);
        $expense->setAccountID($aid);
        $expense->setCCY($ccy);
	}

	$expense->setClassification($classification);
	$expense->setFXAmount($fxAmount);
	$expense->setFXCCY($fxCCY);
	$expense->setFXRate($fxRate);
	$expense->setCommission($commission);
	$expense->setConfirmed(1);
	
	$expenseDB->saveExpense($expense);
	$eid = $expense->getExpenseID();
	my %dids;
	foreach (split /;/, $rawDids)
	{
		$dids{$_} = 1;
	}
	$expensesDB->saveExpenseDocumentMappings($eid, \%dids);
}


sub save_classification
{
    my ($self, $args) = @_; 
    #$classification->setClassificationID($$commands[0]);
    #$classification->setDescription($$commands[1]);
    #$classification->setValidFrom($$commands[2]);
    #$classification->setValidTo($$commands[3]);
    #$classification->setExpense($$commands[4]);
    #$classificationDB->saveClassification($classification);
}

sub save_account
{
    my ($self, $args) = @_; 
	my ($aid, $name, $ccy, $lid, $pid) = ($$args{'aid'}, $$args{'name'}, $$args{'ccy'}, $$args{'lid'}, $$args{'pid'});
	if ($aid eq 'NEW')
	{
		print "Saving new account $name\n";
		$expensesDB->saveAccount($aid, $name, $ccy, $lid, $pid);
	}
}

sub change_amount
{
    my ($self, $args) = @_; 
	my ($eid, $amount) = ($$args{'eid'}, $$args{'amount'});
	$expenseDB->saveAmount($eid, $amount);
}

sub change_classification
{
    my ($self, $args) = @_; 
	my ($eid, $cid) = ($$args{'eid'}, $$args{'cid'});
	$expenseDB->saveClassification($eid, $cid, 1);
}

sub confirm_classification
{
    my ($self, $args) = @_; 
	my ($eid) = ($$args{'eid'});
	$expenseDB->confirmClassification($eid)
}

sub duplicate_expense
{
    my ($self, $args) = @_; 
	my ($eid) = ($$args{'eid'});
	$expenseDB->duplicateExpense($eid);
}

sub tag_expense
{
    my ($self, $args) = @_; 
	my ($eid, $tag) = ($$args{'eid'}, $$args{'tag'});
	$expenseDB->setTagged($eid, $tag);
}
1;

