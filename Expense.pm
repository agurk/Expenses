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

has OriginalLine =>	(	is => 'ro',
						isa => 'Str',
						required => 1,
						reader => 'getOriginalLine',
					);

has ExpenseDate =>	(	is => 'ro',
						isa => 'Str',
						required => 1,
						reader => 'getExpenseDate',
					);

has ExpenseDescription => ( is => 'ro',
							isa => 'Str',
							required => 1,
							reader => 'getExpenseDescription',
						  );

has ExpenseAmount => (	is => 'rw',
						isa => 'Str',
						required => 1,
						reader => 'getExpenseAmount',
						writer => 'setExpenseAmount',
					 );

has ExpenseClassification => (	is => 'rw',
								isa => 'Str',
								reader => 'getExpenseClassification',
								writer => 'setExpenseClassification',
							 );

has AccountName =>	(	is =>  'ro',
						isa => 'Str',
						required => 1,
						reader => 'getAccountName',
					);

sub isValid
{
	my $self = shift;
	return 0 unless (defined $self->getExpenseClassification);
	return 1;
}

1;

