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

	my $expense = $self->_findExpense($aid, $rawLine->getRefID(), $rawLine->getTransactionDate(), $rawLine->getDescription(), $amount, $ccy, $rawLine->isTemporary());
	
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
	# for temporary expenses update the data to final versions
	$expense->setAmount($amount);
    $expense->setDescription($rawLine->getDescription());
	$expense->setReference($rawLine->getRefID()) if ($rawLine->getRefID());
    $expense->setDetailedDescription( $rawLine->getExtraText() );
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
	$expense->setReference($rawLine->getRefID()) if ($rawLine->getRefID());
	$expense->setTemporary($rawLine->isTemporary());
    $expense->setDetailedDescription( $rawLine->getExtraText() );
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

sub _chooseSimilarExpense
{
    my ($self, $rows, $date, $description, $amount, $temporary) = @_; 
    return unless ($rows);
    my $lastDiff = 10000000;
    my $tempTolerance = 0.01;
    my $confirmedTolerance = 0.05;
    my $eid;
    foreach my $row (@$rows)
    {   
        next if ( $amount * $$row[1] < 0 );
        my $diff = abs(abs($$row[1]) - abs($amount)) / abs($amount);
        next if ($temporary and $diff > $tempTolerance);
        next unless ($diff < $confirmedTolerance);
        next unless ($description =~ m/$$row[2]/);
        $eid = $$row[0] if ($diff < $lastDiff);
    }   
    return $eid;
}


