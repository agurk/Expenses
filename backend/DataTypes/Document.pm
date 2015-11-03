#
#===============================================================================
#
#         FILE: Document.pm
#
#  DESCRIPTION: Individual Representation of an Expense
#
#        FILES: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: YOUR NAME (), 
# ORGANIZATION: 
#      VERSION: 1.0
#      CREATED: 08/04/15 17:54
#     REVISION: ---
#===============================================================================

use strict;
use warnings;
 
package Document;
use Moose;

has DocumentID => (	is=>'ro',
					isa => 'Num',
					required => 1,
					default => -1,
					reader => 'getDocumentID',
					writer => 'setDocumentID',
				 );

has ModDate =>	(	is => 'ro',
						isa => 'Str',
						required => 1,
						reader => 'getModDate',
					);

has Filename => ( is	=> 'ro',
				  isa	=> 'Str',
				  reader => 'getFilename',
				);

has FileSize => ( is  => 'ro',
				  isa => 'Num',
				  reader => 'getFileSize',
				  writer => 'setFileSize',
				);

has Text => ( is	=> 'rw',
			  isa	=> 'Str',
			  reader => 'getText',
			  writer => '_setText',
			  default => '',
			);

sub setText
{
	my ($self, $text) = @_;
	$self->_setText($text);
	$self->setTextModDateNow();
}

has TextModDate => (	is	=> 'rw',
						isa	=> 'Str',
						reader => 'getTextModDate',
						writer => 'setTextModDate',
						default => '',
					);

sub setTextModDateNow
{
	my ($self) = @_;
	$self->setTextModDate(gmtime());
}


has ExpenseIDs => ( is=>'rw',
					isa => 'ArrayRef',
					required =>1,
					reader => 'getExpenseIDs',
					writer => '_setExpenseIDs',
					default=> sub { my @empty; return \@empty},
				  );

sub removeAllExpenseIDs
{
	my ($self) = @_;
	my @empty;
	$self->_setExpenseIDs(\@empty);
}

sub addExpenseID
{
    my ($self, $expenseID) = @_; 
    my $eids = $self->getExpenseIDs;
    foreach (@$eids)
    {   
        return if $_ eq $expenseID;
    }   
    push (@$eids, $expenseID);
}

has Deleted => ( is=>'rw',
				 isa => 'Bool',
				 reader => 'isDeleted',
				 writer => 'setDeleted',
				default => 0,
				);

1;

