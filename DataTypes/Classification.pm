#
#===============================================================================
#
#         FILE: Expense.pm
#
#  DESCRIPTION: Individual Representation of a Classification
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy
# ORGANIZATION: 
#      VERSION: 1.0
#      CREATED: 04/02/15 14:26:27
#     REVISION: ---
#===============================================================================

use strict;
use warnings;
 
package Classification;
use Moose;

has Description => (	is	=>'rw',
						isa => 'Str',
						reader => 'getDescription',
						writer => 'setDescription',
					);


has ClassifcationID => (	is=>'ro',
							isa => 'Num',
							reader => 'getClassificationID',
							writer => 'setClassificationID',
							default => -1,
					   );


has ValidFrom => (	is=>'ro',
					isa => 'Str',
					reader => 'getValidFrom',
					writer => 'setValidFrom',
	    		 );


has ValidTo => (	is=>'ro',
					isa => 'Str',
					reader => 'getValidTo',
					writer => 'setValidTo',
			   );

has Expense => (	is=>'rw',
					isa => 'Bool',
					reader => 'isExpense',
					writer => 'setExpense',
			   );


1;

