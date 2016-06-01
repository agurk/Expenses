#
#===============================================================================
#
#         FILE: Processor.pm
#
#  DESCRIPTION: Base class for the group of objects that take a raw row and
#                classify the expense item
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0
#      CREATED: 23/12/14 21:12:34
#     REVISION: ---
#===============================================================================

package Processor;

use strict;
use warnings;

use Moose;
use DataTypes::Expense;

# abstract method
sub processRawLine{exit 1}
# Method takes args: my ($self, $line, $rid, $aid) = @_;

sub _chooseSimilarExpense{return}

sub _findExpense
{
    my ($self, $aid, $reference, $date, $description, $amount, $ccy, $temporary) = @_; 
    my $exesDb = ExpensesDB->new();
    my $exDb = ExpenseDB->new();
    my $expense = $exDb->findExpense($aid, $reference, $date, $description, $amount, $ccy); 
    unless ($expense)
    {   
        my $eid = $self->_chooseSimilarExpense($exesDb->getTemporaryExpenses($aid), $date, $description, $amount, $temporary);
        $expense = $exDb->getExpense($eid) if ($eid);
    }   
    return $expense;
}

1;

