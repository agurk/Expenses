#
#===============================================================================
#
#         FILE: Processor_AMEX.pm
#
#  DESCRIPTION: Class to convert raw amex input lines into an Expense class
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 23/12/14 21:27:55
#     REVISION: ---
#===============================================================================

package Processor_AMEX;
use Moose;
extends 'Processor';

use strict;
use warnings;

use constant INPUT_LINE_PARTS_LENGTH => 4;

sub processRawLine
{
	my ($self, $line, $rid, $aid, $ccy) = @_;
	$line =~ s/^([0-9\/]*),//;
	my $date = $1;
    my @lineParts=split(/","/, $line);
    die "wrong line length\n" unless (scalar @lineParts >= INPUT_LINE_PARTS_LENGTH);
    # Value comes in quotes. Ridiculous.
    $lineParts[1]  =~ s/\"//g;
    $lineParts[1]  =~ s/ //g;
    $lineParts[2]  =~ s/\"//g;
    my $expense = Expense->new (	AccountID => $aid,
									ExpenseDate => $self->_setDate($date),
									ExpenseDescription => $lineParts[2],
									ExpenseAmount => $self->_getAmount($lineParts[1]),
									Currency => $ccy,
                        );
	$expense->addRawID($rid);
	$self->_setFX($expense, $lineParts[3]);
	return $expense;
}

sub _setFX
{
	my ($self, $expense, $description) = @_;
	$description =~ s/\"//g;
	if ($description =~ m/^([0-9.,]{1,}) ([A-Z]{3}).* Currency Conversion Rate ([0-9.,]{1,}) Commission Amount ([0-9,.]*[0-9])/)
	{
		$expense->setFXAmount($1);
		$expense->setFXCCY($2);
		$expense->setFXRate($3);
		$expense->setCommission($4);
	} 
}

sub _getAmount
{
	my ($self, $amount) = @_;
	my $returnAmount;
	if ($amount =~ m/-(.*)/)
	{
		$returnAmount = $1
	}
	else
	{
		$returnAmount = "-$amount";
	}
	return $returnAmount;
}

sub _setDate
{
	my ($self, $dateString) = @_;
	my @dateParts = split (/\//, $dateString);
	return "$dateParts[2]-$dateParts[1]-$dateParts[0]";
}
