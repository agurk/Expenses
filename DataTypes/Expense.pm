#
#===============================================================================
#
#         FILE: Expense.pm
#
#  DESCRIPTION: Individual Representation of an Expense
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: YOUR NAME (), 
# ORGANIZATION: 
#      VERSION: 1.0
#      CREATED: 17/01/13 12:24:27
#     REVISION: ---
#===============================================================================

use strict;
use warnings;
 
package Expense;
use Moose;

has ExpenseID => (	is=>'ro',
					isa => 'Num',
					required => 1,
					default => -1,
					reader => 'getExpenseID',
					writer => 'setExpenseID',
				 );

has RawIDs => (	is=>'ro',
				isa => 'ArrayRef',
				required =>1,
				reader => 'getRawIDs',
				default=> sub { my @empty; return \@empty},
			 );

sub addRawID
{
	my ($self, $rawID) = @_;
	my $rids = $self->getRawIDs;
	foreach (@$rids)
	{
		return if $_ eq $rawID;
	}
	push (@$rids, $rawID);
}


has AccountID => (	is=>'ro',
					isa => 'Num',
					required =>1,
					reader => 'getAccountID',
				 );

has Date =>	(	is => 'ro',
						isa => 'Str',
						required => 1,
						reader => 'getDate',
						writer => 'setDate',
					);

has Description => ( is => 'ro',
							isa => 'Str',
							required => 1,
							reader => 'getDescription',
							writer => 'setDescription',
						  );

has Amount => (	is => 'rw',
						isa => 'Str',
						required => 1,
						reader => 'getAmount',
						writer => 'setAmount',
					 );

has Currency => (   is => 'rw',
					isa => 'Str',
					required => 1,
					reader => 'getCCY',
					writer => 'setCCY',
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

has Classification => (	is => 'rw',
								isa => 'Str',
								reader => 'getClassification',
								writer => 'setClassification',
		);

has Confirmed => (	is => 'rw',
								isa => 'Str',
								reader => 'isConfirmed',
								writer => 'setConfirmed',
		);

has Tagged => ( is => 'rw',
				isa => 'Bool',
				reader => 'isTagged',
				writer => 'setTagged',
				default => 0,
			);

sub isValid
{
	my $self = shift;
	return 0 unless (defined $self->getClassification);
	return 1;
}

1;

