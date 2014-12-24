#
#===============================================================================
#
#         FILE: Processor_Aqua.pm
#
#  DESCRIPTION: Converts Aqua raw(ish) lines into an Expense object
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

use strict;
use warnings;


use constant DATE_INDEX => 0;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;
 
sub processRawLine
{
	my ($self, $line, $rid, $aid) = @_;
    my $lineParts=$self->_splitLine($line);
    $$lineParts[AMOUNT_INDEX] *= -1 if ($$lineParts[CREDIT_DEBIT_INDEX] =~ m/CR/);
    return Expense->new (   RawID => $rid,
							AccountID => $aid,
                            ExpenseDate => $$lineParts[DATE_INDEX],
                            ExpenseDescription => $$lineParts[DESCRIPTION_INDEX],
                            ExpenseAmount => $$lineParts[AMOUNT_INDEX],
                        )
}
