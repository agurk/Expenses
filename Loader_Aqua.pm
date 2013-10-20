#!/usr/bin/perl

package Loader_Aqua;
use Moose;

use strict;
use warnings;

use WWW::Mechanize;
use HTTP::Cookies;

my $INTERNAL_FIELD_SEPARATOR = ';';
use constant IFR_REGEXP => qr/;/;
use constant EMPTY_LINE => ';x;x;x;0;x;';

extends 'Loader';

has 'USER_NAME' => ( is => 'rw', isa=>'Str' );
has 'SURNAME' => ( is => 'rw', isa=>'Str' );
has 'SECRET_WORD' => ( is => 'rw', isa=>'Str' );
has 'SECRET_NUMBERS' => ( is => 'rw', isa=>'Str' );

use constant DATE_INDEX => 0;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;

sub _splitLine
{
    my ($self, $line) = @_;
    my @splitLine=split(IFR_REGEXP, $line);
    return \@splitLine;
}

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

sub _makeRecord
{
    my ($self, $line) = @_;
    my $lineParts=$self->_splitLine($$line);
    #$$lineParts[AMOUNT_INDEX] =~ s/^[^0-9]*//;
    $$lineParts[AMOUNT_INDEX] *= -1 if ($$lineParts[CREDIT_DEBIT_INDEX] =~ m/CR/);
    return Expense->new (   OriginalLine => $$line,
                            ExpenseDate => $$lineParts[DATE_INDEX],
                            ExpenseDescription => $$lineParts[DESCRIPTION_INDEX],
                            ExpenseAmount => $$lineParts[AMOUNT_INDEX],
                            AccountName => $self->account_name,
                        )
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

sub _useInputLine
{
    my ($self, $line) = @_;
    return 1 if ($_ =~ m/abbr/);
    return 0;
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
    unshift (@$newLine, $$newLine[DATE_INDEX]) if ($archiveLine);
    $$newLine[AMOUNT_INDEX] =~ s/[^0-9\.]*//g;
    if (
            ($$newLine[DATE_INDEX] eq $$previousLine[DATE_INDEX])
            and ( $$newLine[CREDIT_DEBIT_INDEX] eq $$previousLine[CREDIT_DEBIT_INDEX])
            and $$newLine[AMOUNT_INDEX] =~ /0*\.00/ 
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
        next unless ($self->_useInputLine($_));
        push(  @$output, $self->_processInputLine($_) );
    }
#    $self->set_input_data(\@output);
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
	return $1;
}

sub _pullOnlineData
{
    my $self = shift;
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

	while ($self->_getPageNumber($agent) > $pageNumber)
	{
		my @lines = split ("\n",$agent->content());
	    $self->_setOutputData(\@lines);
		$pageNumber = $self->_getPageNumber($agent);
		$agent->click_button( name => 'ph_basicpage_content_0$ph_twocolumn_content_2$ph_twocolumnreport_content_0$ph_divsection01_content_1$74c978e5-b650-4c30-8548-49b4775703f6$ctl00$ctl03$ctl01$ctl12');
	}

    $self->_doPostback($agent, 'View statements');
    $self->_doPostback($agent, 'Transactions');

	$pageNumber = 0;
	while ($self->_getPageNumber($agent) > $pageNumber)
    {
        my @lines = split ("\n",$agent->content());
        $self->_setOutputData(\@lines);
		$pageNumber = $self->_getPageNumber($agent);
        $agent->click_button( name => 'ph_basicpage_content_0$ph_twocolumn_content_2$ph_tabbedsublayout_content_3$b7dbcc20-0f57-40a9-bbd7-df3e54f66d67$ctl00$ctl03$ctl01$ctl16');
    }

    return 1;

}

1;

