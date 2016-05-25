#
#===============================================================================
#
#         FILE: Processor_Aqua.pm
#
#  DESCRIPTION: Generic convertor for JSON data from aqua
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

package Processor_Aqua;
use Moose;
extends 'Processor';

use Cpanel::JSON::XS qw(decode_json);

use Database::ExpensesDB;

use strict;
use warnings;

sub processRawLine
{
	my ($self, $json, $rid, $aid, $ccy) = @_;
	my $data = decode_json $json;

	use Data::Dumper;
	print Dumper $data;

	my $amount = $data->{'amount'} * -1;
	my $date = $self->_formatDate($data->{'effectiveDate'});

	my $expense = $self->_findExpense($aid, $date, $data->{'description'}, $amount, $ccy);
	
	unless (defined $expense)
	{
		$expense = Expense->new (
										AccountID => $aid,
										Date => $date,
										Description => $data->{'description'},
										Amount => $amount,
										Currency => $ccy,
								);
	}

	$self->_addFX($expense, $data);
	$expense->setTemporary(1) unless ($data->{'tranRefNo'});
	$expense->addRawID($rid);
#	# for temporary expenses to be updated to the right amount
	$expense->setAmount($amount);
	return $expense;
}

sub reprocess
{
	my ($self, $expense, $json) = @_;
	my $data = decode_json $json;
	my $amount = $data->{'amount'} * -1;
	$expense->setTemporary(1) unless ($data->{'tranRefNo'});
	$expense->setAmount($amount);
	$self->_addFX($expense, $data);
}

sub _formatDate
{
	my ($self, $date) = @_;
	$date =~ m/([0-9]{2})\/([0-9]{2})\/([0-9]{4})/;
	return $3.'-'.$2.'-'.$1
}

sub _addFX
{
	my ($self, $expense, $rawLine) = @_;
	#$expense->setFXAmount($rawLine->getFXAmount) if defined($rawLine->getFXAmount);
	#$expense->setFXCCY($rawLine->getFXCCY) if defined ($rawLine->getFXCCY);
	#$expense->setFXRate($rawLine->getFXRate) if defined ($rawLine->getFXRate);
	#$expense->setCommission($rawLine->getCommission) if defined ($rawLine->getCommission);
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

