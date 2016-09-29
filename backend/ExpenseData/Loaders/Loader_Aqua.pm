#!/usr/bin/perl

package Loader_Aqua;
use Moose;
extends 'Loader';

use strict;
use warnings;

#use JSON::PP;
use Cpanel::JSON::XS qw(decode_json);

use WWW::Mechanize;
use HTTP::Cookies;

use DataTypes::GenericRawLine;

use Switch;

has 'input_data' => ( is => 'rw', isa => 'ArrayRef', writer=>'set_input_data', reader=>'get_input_data', default=> sub { my @empty; return \@empty});
has 'USER_NAME' => ( is => 'rw', isa=>'Str', writer=>'setUserName');
has 'SURNAME' => ( is => 'rw', isa=>'Str', writer=>'setSurname' );
has 'SECRET_WORD' => ( is => 'rw', isa=>'Str', writer=>'setSecretWord' );
has 'SECRET_NUMBERS' => ( is => 'rw', isa=>'Str', writer=>'setSecretNo' );
has 'cutoff_date' => (is => 'rw', isa=>'Str', writer=>'setCutoffDate', reader=>'getCutoffDate');
has 'process_statement' => (is => 'rw', isa=>'Bool', writer=>'setProcessStatement', reader=>'getProcessStatement');

use constant DATE_INDEX => 0;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;

# build string formats:
# file; filename
# notfile; username; surname; secretword; secretnumber; processStatement; (cutoff date)
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
		$self->setProcessStatement($buildParts[5]);
		$self->setCutoffDate($buildParts[6]) if defined ($buildParts[6]);
	}
}

sub _getDate
{
    my ($self, $dateString) = @_; 
	if ($dateString =~ m/([0-9]{2})\/([0-9]{2})\/([0-9]{4})/)
	{
		return "$3-$2-$1";
	}
	elsif ($dateString =~ m/([0-9]{2}) ([A-Za-z]{3}) ([0-9]{2})/)
	{
		my $day = $1;
		my $year = $3;
		my $month;
		switch($2)
		{   
			case 'Jan' { $month = '01'; }
			case 'Feb' { $month = '02'; }
		    case 'Mar' { $month = '03'; }
		    case 'Apr' { $month = '04'; }
		    case 'May' { $month = '05'; }
		    case 'Jun' { $month = '06'; }
		    case 'Jul' { $month = '07'; }
		    case 'Aug' { $month = '08'; }
		    case 'Sep' { $month = '09'; }
		    case 'Oct' { $month = '10'; }
		    case 'Nov' { $month = '11'; }
		    case 'Dec' { $month = '12'; }
		}   
		return "20$year-$month-$day";
	}
}

sub _cleanAmount
{
	my ($self, $amount) = @_;
	$amount =~ m/([0-9,.]+)/;
	return $1;
}

sub _afterCutoffDate
{
	my ($self, $line)  = @_;
	return 1 unless (defined $self->getCutoffDate());
	$self->getCutoffDate() =~ m/([0-9]{4})-([0-9]{2})-([0-9]{2})/;
	my $year = $1;
	my $month = $2;
	my $day = $3;
	$line->getTransactionDate() =~ m/([0-9]{4})-([0-9]{2})-([0-9]{2})/;
	return 1 if ($1 > $year);
	return 1 if ($1 >= $year and $2 > $month);
	return 1 if ($1 >= $year and $2 >= $month and $3 > $day);
	return 0;
}

###############################################################################
## All functions below for navigation of website ##############################
###############################################################################

sub _pullOnlineData
{
	my $self = shift;
	my @result;
	my $cookies = HTTP::Cookies->new;
	my $agent = WWW::Mechanize->new( cookie_jar => $cookies );
#	$agent->add_handler("request_send", sub { print '-' x 80,"\n"; shift->dump(maxlength => 0); return });
#	$agent->add_handler("response_done", sub {print '-' x 80,"\n"; shift->dump(); return });

	my $credentials = '{"userId":"' . $self->USER_NAME . '","passwd":"' . $self->SECRET_WORD . '"}';
	$agent->add_header('Content-Type' => 'application/json');
	$agent->add_header('Connection' => 'keep-alive');
	my $resp = $agent->post('https://portal.aquacard.co.uk/accounts/services/rest/login/preauth?portalName=aqua', Content =>$credentials );

	my $contents= 'j_password='.$self->_getPasscodes($resp).'&portalName=aqua&pageName_TC=termsconditions&pageName_MC=postlogin';
	$agent->add_header('Content-Type' => 'application/x-www-form-urlencoded');
	$agent->add_header('Referer' => 'https://portal.aquacard.co.uk/accounts/aqua/login');
	$agent->add_header('Connection' => 'keep-alive');
	$agent->add_header('Accept' => 'application/json');
	$agent->post('https://portal.aquacard.co.uk/accounts/j_spring_security_check', Content=>$contents);

	$agent->add_header('Content-Type' => 'application/json;charset=utf-8');
	$agent->add_header('Referer' => 'https://portal.aquacard.co.uk/accounts/aqua/account-summary');
	$agent->add_header('Accept' => 'application/json, text/plain, */*');
	$contents = '{"tokenId":null,"cardNumber":null,"actNumber":null,"fromDate":null,"toDate":null,"noOfTransaction":50,"tranNbrMonths":null,"detailFlag":"C","tranStartNum":0,"tranStartDate":null,"tranFileType":null,"org":null,"logo":null,"emblem":null,"contData":null,"pendingFlag":true,"transFlag":false}';
	$agent->post('https://portal.aquacard.co.uk/accounts/services/rest/v1/getAllTransactions', Content=>$contents);

	my $records = decode_json $agent->content();
	my @lines;
    my $js = JSON::PP->new;
    $js->canonical(1);

    foreach (@{$records->{'response'}->{transactionDetails}})
    {  
        push(@lines, $js->encode($_));
	} 

    $self->_getOldTransactions($agent, \@lines);

	return \@lines;
}

sub _getOldTransactions
{
    my ($self, $agent, $lines) = @_;
    #setup
    $agent->add_header('Content-Type' => 'application/json;charset=utf-8');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/accounts/aqua/transactions--statements');
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->post('https://portal.aquacard.co.uk/accounts/services/rest/v1/getStatementDates', Content=>'{}');

    for(my $i=0; $i<6 ; $i++)
    {
        $agent->add_header('Content-Type' => 'application/json;charset=utf-8');
        $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/accounts/aqua/transactions--statements');
        $agent->add_header('Accept' => 'application/json, text/plain, */*');
        $agent->add_header('Accept-Encoding' => 'gzip, deflate, br');
        $agent->add_header('Host' => 'portal.aquacard.co.uk');
        $agent->add_header('Accept-Language' => 'en-GB,en;q=0.5');

        my $tranStartNo = '0';

        my $contents = '{"tokenId":null,"cardNumber":null,"actNumber":"XXXXXXXXXXXXXXX4428","fromDate":null,"toDate":null,"noOfTransaction":50,"tranNbrMonths":'.$i.',"detailFlag":"M","tranStartNum":'.$tranStartNo.',"tranStartDate":null,"tranFileType":null,"org":null,"logo":null,"emblem":null}';
        $agent->post('https://portal.aquacard.co.uk/accounts/services/rest/v1/getTransactions', Content=>$contents);

	    my $records = decode_json $agent->content();
        my $js = JSON::PP->new;
        $js->canonical(1);

        foreach (@{$records->{'response'}->{transactionDetails}})
        {  
            push(@{$lines}, $js->encode($_));
	    }
    }
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
    my ($self, $resp) = @_;
    my $values = $self->_generateSecretNumbers();
	$resp->decoded_content() =~ m/"passcdDigits":\[([0-9]),([0-9]),([0-9])\]/;
	return $$values[$1] . '%7C' . $$values[$2] . '%7C' . $$values[$3];
}

1;

