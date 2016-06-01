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
	my ($date, $reference, $amount, $description, $fxData) = $self->_extractCSVValues($line);

    my $expense = $self->_findExpense($aid, $reference, $date, $description, $amount, $ccy);

    unless (defined $expense)
    {
        $expense = Expense->new (   AccountID => $aid,
                                    Date => $date,
                                    Description => $description,
                                    Amount => $amount,
                                    Currency => $ccy,
                           );
    }

    $expense->setReference($reference);
    $expense->addRawID($rid);
    $self->_setFX($expense, $fxData);
    return $expense;
}

sub _extractCSVValues
{
	my ($self, $line) = @_;
    $line =~ s/^([0-9\/]*),//;
    print $line,"\n";
    my $date = $self->_setDate($1);
    my @lineParts=split(/","/, $line);
    die "wrong line length: $line\n" unless (scalar @lineParts >= INPUT_LINE_PARTS_LENGTH);
    $lineParts[0] =~ m/Reference: ([A-Z0-9]*)/;
    my $reference = $1;
    # Value comes in quotes. Ridiculous.
    $lineParts[1]  =~ s/\"//g;
    $lineParts[1]  =~ s/ //g;
    my $amount = $self->_getAmount($lineParts[1]);

    $lineParts[2]  =~ s/\"//g;
    my $description = $lineParts[2];
	return ($date, $reference, $amount, $description, $lineParts[3]);
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
    elsif ($description =~ m/([0-9.]{1,})  *([A-Z]{3})/)
    {
        $expense->setFXAmount($1);
        $expense->setFXCCY($2);
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

sub reprocess
{
    my ($self, $expense, $line) = @_; 
	my ($date, $reference, $amount, $description, $fxData) = $self->_extractCSVValues($line);

	$expense->setReference($reference);
	$expense->setAmount($amount);
	$expense->setDescription($description);
    $self->_setFX($expense, $fxData);
}

