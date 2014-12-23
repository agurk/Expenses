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
use Expense;

# abstract method
sub processRawLine{exit 1}

1;

