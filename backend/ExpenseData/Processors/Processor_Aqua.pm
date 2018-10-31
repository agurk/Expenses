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
	my $date = $self->_formatDate2($data->{'effectiveDate'});

    my $description = $data->{'description'};
    $description =~ s/  */ /g;
	my $expense = $self->_findExpense($aid, $data->{'referenceNbr'}, $date, $description, $amount, $ccy);
	
	unless (defined $expense)
	{
		$expense = Expense->new (
										AccountID => $aid,
										Date => $date,
										Description => $description,
										Amount => $amount,
										Currency => $ccy,
								);
	}

	$self->_addFX($expense, $data);
    $self->_setTemporary($expense, $data);
	$expense->addRawID($rid);
#	# for temporary expenses to be updated to the right amount
	$expense->setAmount($amount);
	return $expense;
}

sub processRawLineOld
{
	my ($self, $json, $rid, $aid, $ccy) = @_;
	my $data = decode_json $json;

	use Data::Dumper;
	print Dumper $data;

	my $amount = $data->{'amount'} * -1;
	my $date = $self->_formatDate($data->{'effectiveDate'});

	my $expense = $self->_findExpense($aid, $data->{'tranRefNo'}, $date, $data->{'description'}, $amount, $ccy);
	
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
    $self->_setTemporary($expense, $data);
	$expense->addRawID($rid);
#	# for temporary expenses to be updated to the right amount
	$expense->setAmount($amount);
	return $expense;
}

sub _setTemporary
{
    my ($self, $expense, $data) = @_;
    if ($data->{'tranRefNo'})
    {
	    $expense->setTemporary(0);
        $expense->setReference($data->{'tranRefNo'});
    }
    if ($data->{'referenceNbr'})
    {
	    $expense->setTemporary(0);
        $expense->setReference($data->{'referenceNbr'});
    }
    else
    {
	    $expense->setTemporary(1);
    }
}

sub reprocess
{
	my ($self, $expense, $json) = @_;
	my $data = decode_json $json;
	my $amount = $data->{'amount'} * -1;
	$expense->setAmount($amount);
    $self->_setTemporary($expense, $data);
	$self->_addFX($expense, $data);
}

sub _formatDate
{
	my ($self, $date) = @_;
	$date =~ m/([0-9]{2})\/([0-9]{2})\/([0-9]{4})/;
	return $3.'-'.$2.'-'.$1
}

sub _formatDate2
{
	my ($self, $date) = @_;
	$date =~ m/([0-9]{4}-[0-9]{2}-[0-9]{2})/;
	return $1;
}

sub _addFX
{
	my ($self, $expense, $data) = @_;
	$expense->setFXAmount($data->{'foreignTxnAmnt'}) if defined($data->{'foreignTxnAmnt'});
	$expense->setFXCCY($data->{'foreignTxnCurrency'}) if defined ($data->{'foreignTxnCurrency'});
	$expense->setFXRate($data->{'foreignExchangeRate'}) if defined ($data->{'foreignExchangeRate'});
	$expense->setFXRate($data->{'exchangeRate'}) if ((defined ($data->{'exchangeRate'})) && $data->{'exchangeRate'} ne '0');
	#$expense->setCommission($data->{''}) if defined ($data->{''});
}

sub _chooseSimilarExpense
{
    my ($self, $rows, $date, $description, $amount) = @_;
    return unless ($rows);
    $description =~ s/ //g;
    $description = uc($description);
    my $lastDiff = 10000000;
    my $eid;
    foreach my $row (@$rows)
    {
        next if ( $amount * $$row[1] < 0 );
        my $diff = abs(abs($$row[1]) - abs($amount)) / abs($amount);
        next unless ($diff <= 0.02);
        $$row[2] =~ s/ //g;
        $$row[2] = uc($$row[2]);
        next unless ($description eq $$row[2]);
        $eid = $$row[0] if ($diff < $lastDiff);
    }
    return $eid;
}

