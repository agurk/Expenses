#
#===============================================================================
#
#         FILE: GenericRawLine.pm
#
#  DESCRIPTION: Individual Representation of an Expense
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 1.0
#      CREATED: 07/01/15 17:34:27
#     REVISION: ---
#===============================================================================

use strict;
use warnings;
 
package GenericRawLine;
use Moose;



has ProcessedDate =>	(	is => 'rw',
						isa => 'Str',
						reader => 'getProcessedDate',
						writer => 'setProcessedDate',
					);

has TransactionDate =>	(	is => 'rw',
							isa => 'Str',
							reader => 'getTransactionDate',
							writer => 'setTransactionDate',
						);

has Description => ( is => 'rw',
							isa => 'Str',
							reader => 'getDescription',
							writer => 'setDescription',
						  );

has Amount => (	is => 'rw',
						isa => 'Str',
						reader => 'getAmount',
						writer => 'setAmount',
					 );

has DebitCredit => (	is => 'rw',
						isa => 'Str',
						reader => 'getDebitCredit',
						writer => 'setDebitCredit',
					 );

has FXAmount => (	is => 'rw',
					isa => 'Str',
					reader => 'getFXAmount',
					writer => 'setFXAmount',
				);

has FXCCY	=> (	is => 'rw',
					isa => 'Str',
					reader => 'getFXCCY',
					writer => 'setFXCCY',
			   );

has FXRate => (	is => 'rw',
				isa => 'Str',
				reader => 'getFXRate',
				writer => 'setFXRate',
			  );

has Commission	=> (	is => 'rw',
							isa => 'Str',
							reader => 'getCommission',
							writer => 'setCommission',
						);

has referenceID =>  (	is => 'rw',
						isa => 'Str',
						reader => 'getRefID',
						writer => 'setRefID',
                        default => '',
					);

has Temporary => (
					is  => 'rw',
					isa => 'Bool',
					reader => 'isTemporary',
					writer => 'setTemporary',
					default => 0,
				);

has ExtraText => (
                    is => 'rw',
                    isa => 'Str',
                    reader => 'getExtraText',
                    writer => 'setExtraText',
                    default => '',
                );

# Generated CSV line format is:
# transaction date; processed date; description; amount; debit/credit; fx amount; fx ccy; fx rate; commission; referenceID

sub fromString
{
	my ($self, $string) = @_;
	my @creationParts = split (/;/, $string);
	$self->setTransactionDate($creationParts[0]) if defined ($creationParts[0]);
	$self->setProcessedDate($creationParts[1]) if defined ($creationParts[1]);
	$self->setDescription($creationParts[2]) if defined ($creationParts[2]);
	$self->setAmount($creationParts[3]) if defined ($creationParts[3]);
	$self->setDebitCredit($creationParts[4]) if defined ($creationParts[4]);
	$self->setFXAmount($creationParts[5]) if defined ($creationParts[5]);
	$self->setFXCCY($creationParts[6]) if defined ($creationParts[6]);
	$self->setFXRate($creationParts[7]) if defined ($creationParts[7]);
	$self->setCommission($creationParts[8]) if defined ($creationParts[8]);
	$self->setRefID($creationParts[9]) if defined ($creationParts[9]);
	$self->setTemporary($creationParts[10]) if defined ($creationParts[10]);
	$self->setExtraText($creationParts[11]) if defined ($creationParts[11]);
}

sub toString
{
	my ($self) = @_;
	my @output = ('') x 11;
	$output[0] = $self->getTransactionDate();
	$output[1] = $self->getProcessedDate();
	$output[2] = $self->getDescription();
	$output[3] = $self->getAmount();
	$output[4] = $self->getDebitCredit();
	$output[5] = $self->getFXAmount();
	$output[6] = $self->getFXCCY();
	$output[7] = $self->getFXRate();
	$output[8] = $self->getCommission();
	$output[9] = $self->getRefID();
	$output[10] = $self->isTemporary();
	$output[11] = $self->getExtraText();

	my $returnStr;
	my $first = 1;
	foreach (@output)
	{
		$returnStr .= ';' unless ($first);
		$returnStr .= $_ if (defined $_);
		$first = 0;
	}
	return $returnStr;
}

sub isEmpty
{
	my ($self) = @_;
	return 1 if ($self->toString() eq ';' x 11);
	return 0;
}

sub isValid
{
	my ($self) = @_;
    return 0 if ($self->isEmpty());
    return 0 unless ($self->getTransactionDate());
    # I assume we can't have a 0 amount transaction?
    return 0 unless ($self->getAmount());
    return 1;
}

1;

