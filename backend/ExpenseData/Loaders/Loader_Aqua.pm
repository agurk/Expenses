#!/usr/bin/perl

package Loader_Aqua;
use Moose;
extends 'Loader';

use strict;
use warnings;

#use JSON::PP;
use Cpanel::JSON::XS qw(decode_json);
use JSON::Create 'create_json';
use URI::Escape;

use WWW::Mechanize;
use HTTP::Cookies;

use Selenium::Chrome;

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

sub _processJSON
{
    my ($self, $results, $json) = @_;
    my $records = decode_json $json;
    #my $js = JSON::PP->new;
    #$js->canonical(1);
    foreach (@{$records})
    {
        push(@{$results}, create_json($_));
    }

}

###############################################################################
## All functions below for navigation of website ##############################
###############################################################################

sub _pullOnlineData
{
    my $self = shift;
    my @results;

    $self->setFileName('/home/timothy/src/Expenses/data/in/aqua.json');
    foreach (@{$self->_loadCSVRows()})
    {
        $self->_processJSON(\@results, $_);
    }

    return \@results;

}


sub _pullOnlineData_WIP
{
	my $self = shift;
	my @result;
	my $cookies = HTTP::Cookies->new;
	my $agent = WWW::Mechanize->new( cookie_jar => $cookies );
	$agent->add_handler("request_send", sub { print "\n",'-' x 80,"\n"; shift->dump(maxlength => 0); return });
	$agent->add_handler("response_done", sub {print "\n",'-' x 80,"\n"; shift->dump(); return });

    $agent->get('https://portal.aquacard.co.uk/aqua/login');

	my $credentials = '{"username":"' . $self->USER_NAME . '","password":"' . $self->SECRET_WORD . '"}';
	$agent->add_header('Content-Type' => 'application/json');
	$agent->add_header('Connection' => 'keep-alive');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/login');
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $agent->add_header('Request-Id' => '|F0m+d.8bM8Z');
    $agent->add_header('X-XSRF-TOKEN' => $cookies->get_cookies('portal.aquacard.co.uk')->{'XSRF-TOKEN'});
    $cookies->set_cookie(0, 'ai_user', 'O0AXa|2018-05-24T12:51:57.901Z', '/', 'portal.aquacard.co.uk', '443', 0, 0, 86400, 0);
    $cookies->set_cookie(0, 'ai_session', 'XbHN8|1527166318261|1527166318261', '/', 'portal.aquacard.co.uk', '443', 0, 0, 86400, 0);
	my $resp = $agent->post('https://portal.aquacard.co.uk/authentication/getPasscodeChallenge', Content =>$credentials );

    my $challengeResp = $self->_buildChallengeResp($resp);
    my $authToken = decode_json uri_unescape $cookies->get_cookies('portal.aquacard.co.uk')->{'Identity.Cookie'};

    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/login');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $agent->add_header('Request-Id' => '|F0m+d.8bM8Z');
    $agent->add_header('X-XSRF-TOKEN' => $cookies->get_cookies('portal.aquacard.co.uk')->{'XSRF-TOKEN'});
    $agent->add_header('Authorization' => $authToken->{'PartAuthToken'}->{'TokenType'} . ' ' . $authToken->{'PartAuthToken'}->{'Token'});
	$agent->add_header('Content-Type' => 'application/json');
	$agent->add_header('Connection' => 'keep-alive');
    $resp = $agent->post('https://portal.aquacard.co.uk/authentication/submitPasscodeChallenge', content => create_json $challengeResp);

    sleep 10;

    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $agent->add_header('Accept-Encoding' => 'gzip, deflate, br');
    $agent->add_header('Accept-Language' => 'en-GB,en;q=0.5');
    $agent->add_header('Authorization' => $authToken->{'PartAuthToken'}->{'TokenType'} . ' ' . $authToken->{'PartAuthToken'}->{'Token'});
    $agent->add_header('Cache-Control' => undef );
    $agent->add_header('Connection' => 'keep-alive');
    $agent->add_header('Content-Type' => undef);
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Pragma' => undef );
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/login');
    $agent->add_header('Request-Id' => '|V/45Z.7EBeq');
    $agent->add_header('User-Agent' => 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0');
    $agent->add_header('X-XSRF-TOKEN' => undef);
    $agent->add_header('request-context' => 'appId=cid-v1:4727a0e4-4266-4122-8f09-386691c84e3c');
    $resp = $agent->get('https://portal.aquacard.co.uk/api/AccountSummary/AccountSummary');



    $agent->add_header('Authorization' => $authToken->{'PartAuthToken'}->{'TokenType'} . ' ' . $authToken->{'PartAuthToken'}->{'Token'});
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/account-summary');
    $agent->add_header('Request-Id' => '|F0m+d.8bM8Z');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $resp = $agent->get('https://portal.aquacard.co.uk/api/AccountSummary/AccountSummary?bypassCache=true');




    $agent->add_header('Authorization' => $authToken->{'PartAuthToken'}->{'TokenType'} . ' ' . $authToken->{'PartAuthToken'}->{'Token'});
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/account-summary');
    $agent->add_header('Request-Id' => '|F0m+d.8bM8Z');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $resp = $agent->get('https://portal.aquacard.co.uk/api/statements/statementswaitingtobeviewed');

    $agent->add_header('Authorization' => $authToken->{'PartAuthToken'}->{'TokenType'} . ' ' . $authToken->{'PartAuthToken'}->{'Token'});
    $agent->add_header('Host' => 'portal.aquacard.co.uk');
    $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/aqua/account-summary');
    $agent->add_header('Request-Id' => '|F0m+d.8bM8Z');
    $agent->add_header('Accept' => 'application/json, text/plain, */*');
    $resp = $agent->get('https://portal.aquacard.co.uk/api/transactions/unstatemented');

	my $records = decode_json $agent->content();
	my @lines;
    #my $js = JSON::PP->new;
    #$js->canonical(1);

    #my $recordNum = $records->{'response'}->{'tranStartNum'};

    foreach (@{$records})
    { 
        push(@lines, create_json ($_));
	} 

    #$self->_getOldTransactions($agent, \@lines);

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

    for(my $i=0; $i<5 ; $i++)
    {
        $agent->add_header('Content-Type' => 'application/json;charset=utf-8');
        $agent->add_header('Referer' => 'https://portal.aquacard.co.uk/accounts/aqua/transactions--statements');
        $agent->add_header('Accept' => 'application/json, text/plain, */*');
        $agent->add_header('Accept-Encoding' => 'gzip, deflate, br');
        $agent->add_header('Host' => 'portal.aquacard.co.uk');
        $agent->add_header('Accept-Language' => 'en-GB,en;q=0.5');

        my $tranStartNo = '0';
        my $tranStartNoResp = '-1';
        my $tranStartDate = 'null';
        my $tranFileType = 'null';

        while ($tranStartNoResp ne '0')
        {
            my $contents = '{"tokenId":null,"cardNumber":null,"actNumber":"XXXXXXXXXXXXXXX4428","fromDate":null,"toDate":null,"noOfTransaction":50,"tranNbrMonths":'.$i.',"detailFlag":"M","tranStartNum":'.$tranStartNo.',"tranStartDate":'.$tranStartDate.',"tranFileType":'.$tranFileType.',"org":null,"logo":null,"emblem":null}';

            $agent->post('https://portal.aquacard.co.uk/accounts/services/rest/v1/getTransactions', Content=>$contents);

    	    my $records = decode_json $agent->content();
            my $js;# = JSON::PP->new;
            $js->canonical(1);
    
            $tranStartNoResp = $records->{'response'}->{'tranStartNum'};
            $tranStartNo = $tranStartNoResp;
            $tranStartDate = '"' . $records->{'response'}->{'tranStartDate'} . '"' if ($records->{'response'}->{'tranStartDate'});
            $tranFileType = '"S"';

            print "Zero transaction response for $tranStartDate\n" unless (scalar @{$records->{'response'}->{transactionDetails}});
    
            foreach (@{$records->{'response'}->{transactionDetails}})
            {
                push(@{$lines}, $js->encode($_));
    	    }
        }
    }
}

sub _generateSecretNumbers
{
    my $self = shift;
    my @numbers = ();
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

sub _buildChallengeResp
{
    my ($self, $resp) = @_;
    my $values = $self->_generateSecretNumbers();
    my $challenge = decode_json $resp->decoded_content();
    for (my $i = 0; $i < 3; $i++)
    {
        $challenge->{'challenge'}->{'passcodeDigits'}->[$i]->{'digit'} = $$values[$challenge->{'challenge'}->{'passcodeDigits'}->[$i]->{'position'}];
    }

    my %challengeResp;
    $challengeResp{'challengeVerificationToken'} = $challenge->{'challengeVerificationToken'};
    $challengeResp{'passcodeDigits'} = $challenge->{'challenge'}->{'passcodeDigits'};
    $challengeResp{'username'} = $self->USER_NAME;
    $challengeResp{'password'} = $self->SECRET_WORD;
    $challengeResp{'rememberMe'} = 'false';

    return \%challengeResp;
}

1;

