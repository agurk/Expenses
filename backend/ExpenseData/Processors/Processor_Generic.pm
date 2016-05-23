#
#===============================================================================
#
#         FILE: Processor_Generic.pm
#
#  DESCRIPTION: Generic convertor for generated CSV lines into an Expense object
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 23/12/14 21:36:22
#     REVISION: ---
#===============================================================================

package Processor_Generic;
use Moose;
extends 'Processor';

use DataTypes::GenericRawLine;
use Database::ExpensesDB;

use strict;
use warnings;

sub processRawLine
{
	my ($self, $line, $rid, $aid, $ccy) = @_;
	my $rawLine = GenericRawLine->new();
	$rawLine->fromString($line);

	my $amount = $rawLine->getAmount();
	$amount *= -1 if ($rawLine->getDebitCredit eq 'DR');

	my $expense = $self->_findExpense($aid, $rawLine->getTransactionDate(), $rawLine->getDescription(), $amount, $ccy);
	
	unless (defined $expense)
	{
		$expense = Expense->new (
										AccountID => $aid,
										Date => $rawLine->getTransactionDate(),
										Description => $rawLine->getDescription(),
										Amount => $amount,
										Currency => $ccy,
								);
	}

	$self->_addFX($expense, $rawLine);
	#$expense-> = $rawLine->getProcessedDate();
	$expense->setTemporary($rawLine->isTemporary());
	$expense->addRawID($rid);
	# for temporary expenses to be updated to the right amount
	$expense->setAmount($amount);
	return $expense;
}

sub reprocess
{
	my ($self, $expense, $line) = @_;
	my $rawLine = GenericRawLine->new();
	$rawLine->fromString($line);
	my $amount = $rawLine->getAmount();
	$amount *= -1 if ($rawLine->getDebitCredit eq 'DR');
	$expense->setAmount($amount);
	# Currently RO, so can't update
	#$expense->setDate()
	#$expense->setDescription($rawLine->getDescription());
	$self->_addFX($expense, $rawLine);
}

sub _addFX
{
	my ($self, $expense, $rawLine) = @_;
	$expense->setFXAmount($rawLine->getFXAmount) if defined($rawLine->getFXAmount);
	$expense->setFXCCY($rawLine->getFXCCY) if defined ($rawLine->getFXCCY);
	$expense->setFXRate($rawLine->getFXRate) if defined ($rawLine->getFXRate);
	$expense->setCommission($rawLine->getCommission) if defined ($rawLine->getCommission);
}

sub _findExpense
{
	my ($self, $aid, $date, $description, $amount, $ccy) = @_;
	my $db = ExpenseDB->new();
	my $expense = $db->findExpense($aid, $date, $description, $amount, $ccy); 
	unless ($expense)
	{
		$expense = $db->findTemporaryExpense($aid, $description, $amount, $ccy);
	}
	return $expense;
}

