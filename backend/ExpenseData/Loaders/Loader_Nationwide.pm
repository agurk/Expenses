#!/usr/bin/perl

package Loader_Nationwide;
use Moose;
extends 'Loader';

use WWW::Mechanize;

# Date _after_ which new style CSV is used
# format of DD MMM YYYY
has 'changeover_date' => ( is=> 'rw', isa => 'Str' );

has 'NATIONWIDE_ACCOUNT_NUMBER' => ( is => 'rw', isa=>'Str', writer=>'setAccountNo' );
has 'NATIONWIDE_ACCOUNT_NAME' => ( is => 'rw', isa=>'Str', writer=>'setAccountName' );
has 'NATIONWIDE_MEMORABLE_DATA' => ( is => 'rw', isa=>'Str', writer=>'setMemorableData' );
has 'NATIONWIDE_SECRET_NUMBERS' => ( is => 'rw', isa=>'Str', writer=>'setSecretNo' );

# build string formats:
# file; filename
# notfile; accountno; accountName; memorabledata; secretnumbers

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
		$self->setAccountNo($buildParts[1]);
		$self->setAccountName($buildParts[2]);
		$self->setMemorableData($buildParts[3]);
		$self->setSecretNo($buildParts[4]);
	}
}


# The nationwide CSV files have five liens at the top that shouln't be processed
# but we'll do a nice check rather than just ignoring the top five lines!
sub _useInputLine
{
    my $self = shift;
    my $line = shift;
    return 0 if ($line eq '');
    return 0 if ($line eq "\n");
    return 0 if ($line eq "\r");
    return 0 if ($line =~ m/^\"Account Name/);
    return 0 if ($line =~ m/^\"Account Balance/);
    return 0 if ($line =~ m/^\"Available Balance/);
    return 0 if ($line =~ m/^\"Date\",\"Transaction type\"/);
    my @lineParts=split(/,/, $line);
    # skip if no debit - this is not an expense!
    #return 0 if ($lineParts[3] eq " ");
    #return 0 if ($lineParts[3] eq "");
    # could do with a proper date object here...
    #return 0 if ($self->_beforeChangeOver($lineParts[0]));
    return 1;
}

sub _ignoreYear
{
    my ($self, $record) = @_;
    return 0 unless (defined $self->settings->DATA_YEAR);
    $record->getDate =~ m/([0-9]{4}$)/;
    return 0 if ($1 eq $self->settings->DATA_YEAR);
    return 1;
}

# Return true if line should be skipped as predates new format
# CSV (and we're assuming it has been loaded alredy)
sub _beforeChangeOver
{
    my $self = shift;
    return 0 unless (defined $self->changeover_date);
    my $date = shift;
    $date =~ s/"//g;
    my @currentDate = split(/ /,$date);
    my @changeoverDate = split(/ /,$self->changeover_date);
    # Don't skip if current year is after changeover year
    return 0 if ($currentDate[2] > $changeoverDate[2]);
    # Skip if year is before (so we know from now onwards in this 
    # that the years will be the same
    return 1 if ($currentDate[2] < $changeoverDate[2]);
    my %months  = ('Jan',1,'Feb',2,'Mar',3,'Apr',4,'May',5,'Jun',6,
                   'Jul',7,'Aug',8,'Sep',9,'Oct',10,'Nov',11,'Dec',12);
    # If in the same year the month is before the change over, then skip
    # otherwise it's good to go
    return 0 if ($months{$currentDate[1]} > $months{$changeoverDate[1]});
    return 1 if ($months{$currentDate[1]} < $months{$changeoverDate[1]});
    # if in the same month, 
    return 0 if ($currentDate[0] > $changeoverDate[0]);
    return 1;
}

sub _generateSecretNumbers
{
    my $self = shift;
    # start with 0 so we can use 1 for 1 array referencing
    # i.e. first number (we care about) is in array posn 1
    my @numbers = (0);
    push (@numbers, split("", $self->NATIONWIDE_SECRET_NUMBERS));
    return \@numbers;
}

sub _getPasscodes
{
    my ($self, $resp) = @_;
    my @values = @{$self->_generateSecretNumbers()};
    my @returnValues;
    $resp->decoded_content() =~ m/([0-9]).*digit.*([0-9]).*digit.*([0-9])/;
    $returnValues[0] = $values[$1];
    $returnValues[1] = $values[$2];
    $returnValues[2] = $values[$3];
    return \@returnValues;
}

sub _getSecretToken
{
    my ($self, $agent) = @_;
    $agent->content() =~ m/input id="accessmanagementloginloginviamemorabledataandpassnumber" name="__token" type="hidden" value="([^"]*)"/;
    my $value = '' . $1;
    return $value;
}

sub _pullOnlineData
{
    my $self = shift;
    my $agent = WWW::Mechanize->new();
    $agent->get("https://onlinebanking.nationwide.co.uk/AccessManagement/Login") or die "Can't load page\n";

    #my $resp = $agent->post('https://onlinebanking.nationwide.co.uk/AccessManagement/Login/PreparePrompts',
    #                        Content =>'customerNumber=' . $self->NATIONWIDE_ACCOUNT_NUMBER );

    my $resp = $agent->post('https://onlinebanking.nationwide.co.uk/AccessManagement/Login/GetPassnumberDigitPositions',
                            Content =>'customerNumber=' . $self->NATIONWIDE_ACCOUNT_NUMBER );

    my $secretNos = $self->_getPasscodes($resp);

    my $contents = '__token=' . $self->_getSecretToken($agent) 
                  . '&PersistentCookiesEnabled=true'
                  . '&ScreenResolution=1920+x+1080'
                  . '&TimeZoneOffset=1'
                  #. '&LocalTime=' . getLocalTime?
                  . '&PassnumberDigitsLoaded=true'
                  . '&CustomerNumber=' . $self->NATIONWIDE_ACCOUNT_NUMBER
                  . '&RememberCustomerNumber=false'
                  . '&MemorableData=' . $self->NATIONWIDE_MEMORABLE_DATA
                  . '&FirstPassnumberValue=' . $$secretNos[0]
                  . '&SecondPassnumberValue=' . $$secretNos[1]
                  . '&ThirdPassnumberValue=' . $$secretNos[2] ;

    $resp = $agent->post('https://onlinebanking.nationwide.co.uk/AccessManagement/Login/LoginViaMemorableDataAndPassnumber', Content =>$contents );

    my $account_name = $self->NATIONWIDE_ACCOUNT_NAME;
    $agent->follow_link(text_regex => qr/$account_name/);
    my @formFields;
    push(@formFields, 'downloadType');
    $agent->form_with_fields(@formFields);
    $agent->set_fields('downloadType' => 'Csv');
    $agent->submit();
    my @lines = split ("\n",$agent->content());
	return \@lines;
}

1;

