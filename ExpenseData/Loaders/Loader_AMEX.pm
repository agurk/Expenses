#!/usr/bin/perl

package Loader_AMEX;
use Moose;
extends 'Loader';

use WWW::Mechanize;
use LWP::ConnCache;

has 'AMEX_PASSWORD' => ( is => 'rw', isa=>'Str', writer => 'setAmexPass');
has 'AMEX_USERNAME' => ( is => 'rw', isa=>'Str', writer => 'setAmexUser');
has 'AMEX_CARD_NUMBER' => ( is => 'rw', isa=>'Str', writer => 'setAmexCardNo');
# Index is 0-rated
has 'AMEX_INDEX' => ( is => 'rw', isa=>'Str', writer => 'setAmexIndex');

# build string formats:
# file; filename
# notfile; cardno; user; password; index

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
		$self->setAmexCardNo($buildParts[1]);
		$self->setAmexUser($buildParts[2]);
		$self->setAmexPass($buildParts[3]);
		$self->setAmexIndex($buildParts[4]);
	}
}

sub _ignoreYear
{
        my ($self, $record) = @_;
    #    return 0 unless (defined $self->settings->DATA_YEAR);
    #    $record->getDate =~ m/([0-9]{4}$)/;
    #    return 0 if ($1 eq $self->settings->DATA_YEAR);
        return 0;
}

sub _setAgentHeaders
{
	my ($self, $agent) = @_;
	$agent->add_header(  Connection => 'keep-alive');
	$agent->add_header(  DNT => 1 );
	$agent->add_header( 'Accept-Encoding' => 'gzip, deflate' );
	$agent->add_header(  Accept => 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8');
	$agent->add_header( 'Accept-Language' =>'en-GB,en;q=0.5' );
	$agent->add_header( 'User-Agent' => 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:41.0) Gecko/20100101 Firefox/41.0' );
}

# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number
# selectradio => with the value set to the statement periods we want to download
sub _pullOnlineData
{
    my $self = shift;
    my $agent = WWW::Mechanize->new( cookie_jar => {} );
	$agent->conn_cache(LWP::ConnCache->new);
	$agent->add_handler("request_send", sub { print '-' x 80,"\n"; shift->dump(maxlength => 0); return });
	$agent->add_handler("response_done", sub {print '-' x 80,"\n"; shift->dump(); return });

	$self->_setAgentHeaders($agent);
	$agent->get("https://www.americanexpress.com/");
	$agent->form_id('ssoform');
	$agent->current_form()->action('https://online.americanexpress.com/myca/logon/us/action/LogLogonHandler?request_type=LogLogonHandler&Face=en_US');
	$agent->set_fields( UserID => $self->AMEX_USERNAME );
	$agent->set_fields( Password => $self->AMEX_PASSWORD );
	$agent->set_fields( act => 'soa'  );
	$agent->set_fields( DestPage => 'https://online.americanexpress.com/myca/acctmgmt/us/myaccountsummary.do?request_type=authreg_acctAccountSummary&Face=en_US&omnlogin=us_homepage_myca' );
	$self->_setAgentHeaders($agent);
	$agent->submit();

	$self->_setAgentHeaders($agent);
    $agent->follow_link(text => 'View transactions');

	$self->_setAgentHeaders($agent);
    $agent->follow_link(text => 'Export Statement Data');

    $agent->form_name('DownloadForm');
	$agent->field('Format' => 'CSV' );
	$agent->field('downloadFormat' => 'on' );
	for my $input ($agent->current_form()->inputs)
	{
		$input->check if (('selectradio' eq $input->name) and ($input->id =~ m/radioid0/));
		if (('selectradio' eq $input->name) and ($input->id =~ m/selectCard10/))
		{
			$input->check;
			$input->value($self->AMEX_CARD_NUMBER);
		}
	}
	$self->_setAgentHeaders($agent);
	$agent->add_header( Host => 'global.americanexpress.com');
    $agent->submit();

    # Assume the download has failed if this string is in the results
    if ($agent->content() =~ m/DownloadErrorPage/)
    {
        print " AMEX failed, retrying ";
        return 0;
    }
    my @lines = split ("\n",$agent->content());
    return \@lines;
}

sub _checkNumberOnPage
{
    my $self = shift;
    my $agent = shift;
    my %foundNumbers;
    my $content = $agent->content();
    while ( $content =~ m/([0-9]{10,})/g )
    {
        $foundNumbers{$1} = 1;
    }
    return \%foundNumbers;
}

1;

