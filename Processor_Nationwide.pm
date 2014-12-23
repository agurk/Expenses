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
 

sub processRawLine
{
    my ($line) = @_; 
    # Strip leading char - £ sign specifically
    my @lineParts=split(/,/, $$line);
    $lineParts[3] =~ s/^[^0123456789\.]*//;
    $lineParts[0] =~ s/\"//g;
    $lineParts[3] =~ s/\"//g;
    return Expense->new (    OriginalLine => $$line,
                            ExpenseDate => $lineParts[0],
                            ExpenseDescription => $lineParts[1] .' '. $lineParts[2],
                            ExpenseAmount => $lineParts[3],
                            AccountName => $self->account_name,
                        )
}

