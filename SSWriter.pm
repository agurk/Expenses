#!/usr/bin/perl

package SSWriter;
use Moose;

has 'user_name' => ( is => 'rw', isa => 'Str', required => 1);
has 'password' => ( is => 'rw', isa => 'Str', required => 1);
has 'workbook' => ( is => 'rw', isa => 'Str', required => 1);
has 'worksheet' => ( is => 'rw', isa => 'Str', required => 1);

use Net::Google::Spreadsheets;

sub write_to_sheet
{
    my $self = shift;
    my $values = shift;
    my $service = Net::Google::Spreadsheets->new(
	username => $self->user_name(),
	password => $self->password());
    my @spreadsheets = $service->spreadsheets();

    my $spreadsheet = $service->spreadsheet({ title => $self->workbook()});
    my $worksheet = $spreadsheet->worksheet({ title => $self->worksheet()});

  # update cell by batch request
    $worksheet->batchupdate_cell(@$values);

}

sub addDataSets
{
    my $self = shift;
    my @arrayRefs = @_;
    my @arrayOut;
    foreach (@arrayRefs)
    {
	foreach (@$_)
	{
	    push(@arrayOut, $_);
	}
    }
    return \@arrayOut;
}

sub createRowDays_HACK
{
    # starts at row 2, so row 1 can have month names in it 
    my $self = shift;
    my ($rowNum, $arrayIn) = @_;
    return (
        {row => $rowNum, col => 2, input_value => $$arrayIn[1]},
        {row => $rowNum, col => 3, input_value => $$arrayIn[2]},
        {row => $rowNum, col => 4, input_value => $$arrayIn[3]},
        {row => $rowNum, col => 5, input_value => $$arrayIn[4]},
        {row => $rowNum, col => 6, input_value => $$arrayIn[5]},
        {row => $rowNum, col => 7, input_value => $$arrayIn[6]},
        {row => $rowNum, col => 8, input_value => $$arrayIn[7]},
        {row => $rowNum, col => 9, input_value => $$arrayIn[8]},
        {row => $rowNum, col => 10, input_value => $$arrayIn[9]},
        {row => $rowNum, col => 11, input_value => $$arrayIn[10]},
        {row => $rowNum, col => 12, input_value => $$arrayIn[11]},
        {row => $rowNum, col => 13, input_value => $$arrayIn[12]},
        {row => $rowNum, col => 14, input_value => $$arrayIn[13]},
        {row => $rowNum, col => 15, input_value => $$arrayIn[14]},
        {row => $rowNum, col => 16, input_value => $$arrayIn[15]},
        {row => $rowNum, col => 17, input_value => $$arrayIn[16]},
        {row => $rowNum, col => 18, input_value => $$arrayIn[17]},
        {row => $rowNum, col => 19, input_value => $$arrayIn[18]},
        {row => $rowNum, col => 20, input_value => $$arrayIn[19]},
        {row => $rowNum, col => 21, input_value => $$arrayIn[20]},
        {row => $rowNum, col => 22, input_value => $$arrayIn[21]},
        {row => $rowNum, col => 23, input_value => $$arrayIn[22]},
        {row => $rowNum, col => 24, input_value => $$arrayIn[23]},
        {row => $rowNum, col => 25, input_value => $$arrayIn[24]},
        {row => $rowNum, col => 26, input_value => $$arrayIn[25]},
        {row => $rowNum, col => 27, input_value => $$arrayIn[26]},
        {row => $rowNum, col => 28, input_value => $$arrayIn[27]},
        {row => $rowNum, col => 29, input_value => $$arrayIn[28]},
        {row => $rowNum, col => 30, input_value => $$arrayIn[29]},
        {row => $rowNum, col => 31, input_value => $$arrayIn[30]},
        {row => $rowNum, col => 32, input_value => $$arrayIn[31]},
    );
}

sub createRowDays
{
    my $self = shift;
    my ($rowNum, $arrayIn) = @_;
    my @results;
    my $value =0;
    for (my $i=0; $i<32;$i++)
    {
	my $columnNum = 1+$i;
	$value += $$arrayIn[$i];
	$results[$i]= {row => $rowNum, col => $columnNum, input_value => $value};
    }
    return (\@results);
}

sub createRowMonth
{
    my $self = shift;
    my ($rowNum, $arrayIn) = @_;
    return (
    {row => $rowNum, col => 1, input_value => $$arrayIn[5]},
    {row => $rowNum, col => 2, input_value => $$arrayIn[6]},
    {row => $rowNum, col => 3, input_value => $$arrayIn[1]},
    {row => $rowNum, col => 4, input_value => $$arrayIn[2]},
    {row => $rowNum, col => 5, input_value => $$arrayIn[3]},
    {row => $rowNum, col => 6, input_value => $$arrayIn[4]},
    {row => $rowNum, col => 7, input_value => $$arrayIn[7]},
    {row => $rowNum, col => 8, input_value => $$arrayIn[8]},
    {row => $rowNum, col => 9, input_value => $$arrayIn[9]},
    {row => $rowNum, col => 10, input_value => $$arrayIn[10]},
    );
}


1;

