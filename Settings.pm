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

has 'CLASSIFICATIONS' => (is=>'rw', isa=>'HashRef');
has 'CLASSIFICATIONS_COUNT' => (is=>'rw', isa=>'Str');

has 'ACCOUNT_FILE' => (is=>'rw', isa=>'Str');

has 'DATA_YEAR' => (is=>'rw', isa=>'Str');

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
    $self->DATAFILE_NAME('DATAFILE');
    $self->CLASSIFICATIONS(_loadClassifications('CLASSIFICATIONS'));
    $self->CLASSIFICATIONS_COUNT(scalar(keys %{$self->CLASSIFICATIONS}));
    $self->ACCOUNT_FILE('ACCOUNTS');
	$self->DATA_YEAR('');
}

1;


