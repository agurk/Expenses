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

use strict;
use warnings;

use constant DATE_INDEX => 0;
use constant PROCESSED_DATE_INDEX => 1;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;
use constant FX_AMOUNT_INDEX => 5;
use constant FX_CCY_INDEX => 6;
use constant FX_RATE_INDEX => 7;
use constant COMMISSION_INDEX => 8;

# Generated CSV line format is:
# transaction date; processed date; description; amount; debit/credit; fx amount; fx ccy; fx rate; commission
sub processRawLine
{
	my ($self, $line, $rid, $aid, $ccy) = @_;
    my @lineParts=split(/;/, $line);
	$lineParts[AMOUNT_INDEX] *= -1 if ($lineParts[CREDIT_DEBIT_INDEX] =~ m/DR/);
    my $expense = Expense->new (
							AccountID => $aid,
                            ExpenseDate => $lineParts[DATE_INDEX],
                            ExpenseDescription => $lineParts[DESCRIPTION_INDEX],
                            ExpenseAmount => $lineParts[AMOUNT_INDEX],
						    Currency => $ccy,
                        );
	$expense->addRawID($rid);
	$self->_addFX($expense, \@lineParts);
	return $expense;
}

sub _addFX
{
	my ($self, $expense, $lineParts);
	$expense->setFXAmount($$lineParts[FX_AMOUNT_INDEX]) if (defined $$lineParts[FX_AMOUNT_INDEX]);
	$expense->setFXCCY($$lineParts[FX_CCY_INDEX]) if (defined $$lineParts[FX_CCY_INDEX]);
	$expense->setFXRate($$lineParts[FX_RATE_INDEX]) if (defined $$lineParts[FX_RATE_INDEX]);
	$expense->setCommission($$lineParts[COMMISSION_INDEX]) if (defined $$lineParts[COMMISSION_INDEX]);
}

