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

use strict;
use warnings;

sub processRawLine
{
	my ($self, $line, $rid, $aid, $ccy) = @_;
	my $rawLine = GenericRawLine->new();
	$rawLine->fromString($line);

	my $amount = $rawLine->getAmount();
	$amount *= -1 if ($rawLine->getDebitCredit eq 'DR');

	my $expense = Expense->new (
									AccountID => $aid,
									Date => $rawLine->getTransactionDate(),
									Description => $rawLine->getDescription(),
									Amount => $amount,
									Currency => $ccy,
							   );
	$expense->addRawID($rid);
	$self->_addFX($expense, $rawLine);
	return $expense;
}

sub _addFX
{
	my ($self, $expense, $rawLine) = @_;
	$expense->setFXAmount($rawLine->getFXAmount) if defined($rawLine->getFXAmount);
	$expense->setFXCCY($rawLine->getFXCCY) if defined ($rawLine->getFXCCY);
	$expense->setFXRate($rawLine->getFXRate) if defined ($rawLine->getFXRate);
	$expense->setCommission($rawLine->getCommission) if defined ($rawLine->getCommission);
}

