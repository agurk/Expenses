#!/usr/bin/perl

package Loader_Aqua;
use Moose;
extends 'Loader';

use strict;
use warnings;

use WWW::Mechanize;
use HTTP::Cookies;

use feature "switch";

my $INTERNAL_FIELD_SEPARATOR = '!';
use constant IFR_REGEXP => qr/!/;
use constant EMPTY_LINE => '!x!x!x!0!x!';

has 'input_data' => ( is => 'rw', isa => 'ArrayRef', writer=>'set_input_data', reader=>'get_input_data', default=> sub { my @empty; return \@empty});
has 'USER_NAME' => ( is => 'rw', isa=>'Str', writer=>'setUserName');
has 'SURNAME' => ( is => 'rw', isa=>'Str', writer=>'setSurname' );
has 'SECRET_WORD' => ( is => 'rw', isa=>'Str', writer=>'setSecretWord' );
has 'SECRET_NUMBERS' => ( is => 'rw', isa=>'Str', writer=>'setSecretNo' );

use constant DATE_INDEX => 0;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;

# build string formats:
# file; filename
# notfile; username, surname, secretword, secretnumber

sub BUILD
{
	my ($self) = @_;
	my @buildParts = split (';' ,$self->build_string);
	# if it is a file
	if ($buildParts[0])
	{
		$self->setFileName($buildParts[1]);
	}
	else
	{
		$self->setUserName($buildParts[1]);
		$self->setSurname($buildParts[2]);
		$self->setSecretWord($buildParts[3]);
		$self->setSecretNo($buildParts[4]);
	}
}

sub _splitLine
{
    my ($self, $line) = @_;
    my @splitLine=split(IFR_REGEXP, $line);
    return \@splitLine;
}

# Generated CSV line format is:
# transaction date; processed date; description; amount; debit/credit; fx amount; fx ccy; fx rate

sub _makeLine
{
	my ($self, $line) = @_;
	my $joinedLine;
	foreach (@$line)
	{
		$joinedLine .= $_;
		$joinedLine .= $INTERNAL_FIELD_SEPARATOR;
	}
	return $joinedLine;
}

sub _ignoreYear
{
    my ($self, $record) = @_;
    return 0 unless (defined $self->settings->DATA_YEAR);
    $record->getExpenseDate =~ m/([0-9]{2}$)/;
    my $found_year = $1;
    $self->settings->DATA_YEAR =~ m/([0-9]{2}$)/;
    return 0 if ($found_year eq $1);
    return 1;
}

sub _generateSecretNumbers
{
    my $self = shift;
# start with 0 so we can use 1 for 1 array referencing
# i.e. first number (we care about) is in array posn 1
    my @numbers = (0);
    push (@numbers, split('', $self->SECRET_NUMBERS));
    return \@numbers;
}

sub _getPasscodes
{
    my ($self, $agent) = @_;
    my @returnCodes;
    my $values = $self->_generateSecretNumbers();
    $agent->content() =~ m/([0-9]).. number of your Passcode.*([0-9]).. number of your Passcode/s;
    $returnCodes[0] = $$values[$1];
    $returnCodes[1] = $$values[$2];
    return \@returnCodes;
}

sub _getNextPageLinkName
{
    my ($self, $agent) = @_;
    $agent->content() =~ m/<input type="submit" name="([^"]*)" value=" " title="Next Page" class="rgPageNext" \/>/;
    return $1;
}

# Use the line when contstructing the record
sub _useInputLineInternal
{
    my ($self, $line) = @_;
    return 1 if ($_ =~ m/abbr/);
    return 0;
}

sub _setStatementDate
{
    my ($self, $line) = @_;
    if ($$line[DATE_INDEX] =~ m/([0-9]{2})\/([0-9]{2})\/[0-9]{2}([0-9]{2})/)
    {
	my $month;
	given ($2)
	{   
	    when ('01') { $month = 'Jan'; }
	    when ('02') { $month = 'Feb'; }
	    when ('03') { $month = 'Mar'; }
	    when ('04') { $month = 'Apr'; }
	    when ('05') { $month = 'May'; }
	    when ('06') { $month = 'Jun'; }
	    when ('07') { $month = 'Jul'; }
	    when ('08') { $month = 'Aug'; }
	    when ('09') { $month = 'Sep'; }
	    when ('10') { $month = 'Oct'; }
	    when ('11') { $month = 'Nov'; }
	    when ('12') { $month = 'Dec'; }
	}
	$$line[DATE_INDEX] = $1 . ' ' . $month . ' ' . $3;
    }
}

sub _processInputLine
{
    my ($self, $line ) = @_;
    my $archiveLine = 0;
    $archiveLine = 1 if ($line =~ m/class="thirtytwo"/);
    chomp $line;
    $line =~ s/^[\w\t ]+//g;
    $line =~ s/^<td[^>]*>//g;
    $line =~ s/<\/td><td[^>]*>/$INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<\/td><td class="right">/$INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<abbr title="[^\"]*">/$INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<\/abbr> *<\/td>//g;
    my $newLine = $self->_splitLine($line);

    my $previousLineIn = $self->get_input_data()->[-1];
    $previousLineIn = EMPTY_LINE unless (defined $previousLineIn);
    my $previousLine = $self->_splitLine($previousLineIn);
    $self->_setStatementDate($newLine) if ($archiveLine);
    unshift (@$newLine, $$newLine[DATE_INDEX]) if ($archiveLine);
    $$newLine[AMOUNT_INDEX] =~ s/[^0-9\.]*//g;
    if (
            ($$newLine[DATE_INDEX] eq $$previousLine[DATE_INDEX])
            and ( $$newLine[CREDIT_DEBIT_INDEX] eq $$previousLine[CREDIT_DEBIT_INDEX])
            and $$newLine[AMOUNT_INDEX] =~ m/^0?\.00$/ 
       )
    {
        $$previousLine[DESCRIPTION_INDEX] .= ' ';
        $$previousLine[DESCRIPTION_INDEX] .= $$newLine[DESCRIPTION_INDEX];
        pop @{$self->get_input_data()};
        return $self->_makeLine($previousLine);
    } else {
        return $self->_makeLine($newLine);
    }
}

sub _setOutputData
{
    my ($self, $lines) = @_;
    my $output = $self->get_input_data();
    foreach (@$lines)
    {
        next unless ($self->_useInputLineInternal($_));
        push( @$output, $self->_processInputLine($_) );
    }
}

sub _doPostback
{
    my ($self, $agent, $searchString) = @_;
    $agent->content() =~ m/(doPostBack.*$searchString)/;
    my $string = $1;
    $string =~ m/^doPostBack\('([^']*)','([^']*)'\)/;
    $agent->form_id('mainform');
    $agent->set_fields('__EVENTTARGET' => "$1");
    $agent->set_fields('__EVENTARGUMENT' => "$2");
    $agent->submit();
}

sub _getPageNumber
{
    my ($self, $agent) = @_;
    $agent->content() =~ m/rgCurrentPage[^<]*<span>([^<]*)<\/span>/;
	my $pageNumber = $1;
	$pageNumber = 0 if ($pageNumber eq '' or ! defined $pageNumber);
    return $pageNumber;
}

sub _loadCSVRows2
{
    my ($self) = @_;
	my @allLines;
    open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
    foreach (<$file>)
    {
		chomp;
		push(@allLines, $_);
    }
    close($file);
	$self->_setOutputData(\@allLines);

    return $self->get_input_data();
}

sub _pullOnlineData
{
    my $self = shift;
	my @result;
    my $agent = WWW::Mechanize->new( cookie_jar => {} );
    $agent->get('https://service.aquacard.co.uk/aqua/web_channel/cards/security/logon/logon.aspx');
    $agent->form_id('mainform');
    $agent->set_fields('datasource_3a4651d1-b379-4f77-a6b1-5a1a4855a9fd' => $self->USER_NAME);
    $agent->set_fields('datasource_2e7ba395-e972-4a25-8d19-9364a6f06132' => $self->SURNAME);
    $agent->set_fields( '__EVENTTARGET' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
    $agent->set_fields( '__EVENTARGUMENT' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
    $agent->submit();

    $agent->form_id('mainform');
    $agent->set_fields('datasource_20056b2c-9455-4e42-aeba-9351afe0dbe1' => $self->SECRET_WORD );

    my $secretNumbers = $self->_getPasscodes($agent);
    $agent->set_fields('selectedvalue_dc28fef3-036f-4e48-98a0-586c9a4fbb3c' => $$secretNumbers[0]);
    $agent->set_fields('selectedvalue_f758fdb6-4b2c-4272-b785-cb3989b67901' => $$secretNumbers[1]);
    $agent->set_fields( '__EVENTTARGET' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
    $agent->set_fields( '__EVENTARGUMENT' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
    $agent->submit();

    my $pageNumber = 0;    
    my $linkName = $self->_getNextPageLinkName($agent);

    while ($self->_getPageNumber($agent) > $pageNumber)
    {
        my @lines = split ("\n",$agent->content());
        $self->_setOutputData(\@lines);
        $pageNumber = $self->_getPageNumber($agent);
        $agent->click_button( name => $linkName );
    }
	return $self->get_input_data();

#    $self->_doPostback($agent, 'View statements');
#    $self->_doPostback($agent, 'Transactions');
#
#    $pageNumber = 0;
#    $linkName = $self->_getNextPageLinkName($agent);
#
#    while ($self->_getPageNumber($agent) > $pageNumber)
#    {
#        my @lines = split ("\n",$agent->content());
#        $self->_setOutputData(\@lines);
#        $pageNumber = $self->_getPageNumber($agent);
#        $agent->click_button( name => $linkName );
#    }

}

1;

