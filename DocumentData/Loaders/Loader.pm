#
#===============================================================================
#
#         FILE: Loader.pm
#
#  DESCRIPTION: Base class for the group of objects that load documents
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0
#      CREATED: 01/05/15 17:41
#     REVISION: ---
#===============================================================================

package Loader;

use strict;
use warnings;

use Moose;

has 'DocumentDir' => ( isa => 'Str',
					   is => 'rw',
				   	   default => '/home/timothy/bin/Expenses/data/documents',
				   	   reader => 'getDocumentDir',
				   	   writer => 'setDocumentDir',
					);

# abstract method
sub processDocument{exit 1}

1

