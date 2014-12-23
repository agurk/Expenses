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
	my ($line) = @_;
    my @lineParts=split(/,/, $$line);
    die "wrong line length\n" unless (scalar @lineParts >= INPUT_LINE_PARTS_LENGTH);
    # Value comes in quotes. Ridiculous.
    $lineParts[2]  =~ s/\"//g;
    return Expense->new ( OriginalLine => $$line,
                          ExpenseDate => $lineParts[0],
                          ExpenseDescription => $lineParts[3],
                          ExpenseAmount => $lineParts[2],
                          AccountName => $self->account_name,
                        )
}
