#
#===============================================================================
#
#         FILE: Processor_Nationwide.pm
#
#  DESCRIPTION: Class to process raw nationwide input lines into an Expenses class
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.1
#      CREATED: 23/12/14 21:33:32
#     REVISION: ---
#===============================================================================

package Processor_Nationwide;
use Moose;
extends 'Processor';

use strict;
use warnings;
 
use Switch;

sub processRawLine
{
	my ($self, $line, $rid, $aid, $ccy) = @_;
    # Strip leading char - £ sign specifically
    my @lineParts=split(/","/, $line);
	# deal with the very old format CSVs
	if (scalar @lineParts < 6)
	{
		$line =~ s/,"([^\"]*)",/, ,/;
		@lineParts = split(/,/, $line);
		$lineParts[1] = $1;
		splice (@lineParts, 2, 0, $1);
	}

    $lineParts[0] =~ s/\"//g;
    $lineParts[1] =~ s/\"//g;
    $lineParts[2] =~ s/\"//g;
    $lineParts[3] =~ s/[^0123456789\.]//g;
    $lineParts[3] =~ s/\"//g;
    $lineParts[4] =~ s/[^0123456789\.]//g;
    $lineParts[4] =~ s/\"//g;
    my $expense = Expense->new (
							AccountID => $aid,
                            Date => $self->_setDate($lineParts[0]),
                            Description => $self->_makeDescription($lineParts[1], $lineParts[2]),
                            Amount => $self->_getAmount($lineParts[3], $lineParts[4]),
						    Currency => $ccy,
                        );
	$expense->addRawID($rid);
	$self->_setFX($expense, $lineParts[2]);
	return $expense;
}

sub _setFX
{
	my ($self, $expense, $description) = @_;
	if ( $description =~ m/([0-9,.]{1,}) ?([A-Z]{3}) at ([0-9,.]*[0-9])/ )
	{
		$expense->setFXAmount($1);
		$expense->setFXCCY($2);
		$expense->setFXRate($3);
	}
}

sub _getAmount
{
	my ($self, $debit, $credit) = @_;
	return "-$debit" if ($debit);
	return "$credit" if ($credit);
	warn "Invalid amount: no debit or credit\n";
}

sub _makeDescription
{
	my ($self, $desc1, $desc2) = @_;
	return $desc1 if ($desc1 eq $desc2);
	return $desc1 if ($desc2 eq '');
	return $desc2 if ($desc1 eq '');
	return $desc1 . ' - ' . $desc2;
}

sub _setDate
{
	my ($self, $dateString) = @_;
	my @dateParts = split (/ /, $dateString);
	my $month;
	switch ($dateParts[1])
    {   
        case 'Jan' { $month = '01'; }
        case 'Feb' { $month = '02'; }
        case 'Mar' { $month = '03'; }
        case 'Apr' { $month = '04'; }
        case 'May' { $month = '05'; }
        case 'Jun' { $month = '06'; }
        case 'Jul' { $month = '07'; }
        case 'Aug' { $month = '08'; }
        case 'Sep' { $month = '09'; }
        case 'Oct' { $month = '10'; }
        case 'Nov' { $month = '11'; }
        case 'Dec' { $month = '12'; }
    } 
	return "$dateParts[2]-$month-$dateParts[0]";
}
