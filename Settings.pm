#!/usr/bin/perl

# Package to process the raw data from an input file ready to be input into the store numbers

package Settings;
use Moose;

#has 'numbers_store' => (is => 'r', isa => 'Numbers', required => 1);
has 'GOOGLE_DOCS_USERNAME' => ( is => 'rw', isa => 'Str');
has 'GOOGLE_DOCS_PASSWORD' => ( is => 'rw', isa => 'Str');
has 'GOOGLE_DOCS_WORKBOOK' => ( is => 'rw', isa => 'Str');
has 'GOOGLE_DOCS_WORKSHEET' => ( is => 'rw', isa => 'Str');

has 'DATAFILE_NAME' => ( is => 'rw', isa=>'Str' );

has 'AMEX_PASSWORD' => ( is => 'rw', isa=>'Str' );
has 'AMEX_USERNAME' => ( is => 'rw', isa=>'Str' );
has 'AMEX_CARD_NUMBER' => ( is => 'rw', isa=>'Str' );

has 'CLASSIFICATIONS' => (is=>'rw', isa=>'HashRef');

sub _loadClassifications
{
    my $fileName = shift;
    my %classifications;
    open (my $fh, "<",$fileName) or die "Cannot open classifications file $fileName\n";
    foreach (<$fh>)
    {
	chomp;
	my @lineParts = split(/,/, $_);
	$classifications{$lineParts[0]} = $lineParts[1];
    }
    close ($fh);
    return \%classifications;
}

sub BUILD
{
    my $self = shift;
    $self->GOOGLE_DOCS_USERNAME('');
    $self->GOOGLE_DOCS_PASSWORD('');
    $self->GOOGLE_DOCS_WORKBOOK('');
    $self->GOOGLE_DOCS_WORKSHEET('');
    $self->DATAFILE_NAME('');
    $self->AMEX_USERNAME('');
    $self->AMEX_PASSWORD('');
    $self->AMEX_CARD_NUMBER('');
    $self->CLASSIFICATIONS(_loadClassifications('CLASSIFICATIONS'));
}

1;


