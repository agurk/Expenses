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
					reader => 'getExpenseID',
					writer => 'setExpenseID',
				 );

has RawID => (	is=>'ro',
				isa => 'Num',
				required =>1,
				reader => 'getRawID',
			 );


has AccountID => (	is=>'ro',
					isa => 'Num',
					required =>1,
					reader => 'getAccountID',
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

sub isValid
{
	my $self = shift;
	return 0 unless (defined $self->getExpenseClassification);
	return 1;
}

1;

