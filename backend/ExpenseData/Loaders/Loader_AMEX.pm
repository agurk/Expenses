#!/usr/bin/perl

package Loader_AMEX;
use Moose;
extends 'Loader';

use WWW::Mechanize::Firefox;
use LWP::ConnCache;

has 'AMEX_PASSWORD' => ( is => 'rw', isa=>'Str', writer => 'setAmexPass');
has 'AMEX_USERNAME' => ( is => 'rw', isa=>'Str', writer => 'setAmexUser');
has 'AMEX_CARD_NUMBER' => ( is => 'rw', isa=>'Str', writer => 'setAmexCardNo');
# Index is 0-rated
has 'AMEX_INDEX' => ( is => 'rw', isa=>'Str', writer => 'setAmexIndex');

use constant csvFile => '/home/timothy/Downloads/ofx.csv';

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
    unlink csvFile;
    my $agent = WWW::Mechanize::Firefox->new();

    $agent->get("https://www.americanexpress.com/");

    $agent->form_id('ssoform');
    $agent->set_fields(UserID => $self->AMEX_USERNAME);
    $agent->set_fields(Password => $self->AMEX_PASSWORD);
    $agent->follow_link(text => 'Log In');

    $self->_waitForElement($agent, '//*[@id="lilo_userName"]');
    $agent->form_id('lilo_formLogon');
    $agent->set_fields(UserID => $self->AMEX_USERNAME);
    $agent->set_fields(Password => $self->AMEX_PASSWORD);
    $agent->submit();
    
    $self->_waitForElement($agent, '//*[@id="estatement-link"]');
    $agent->follow_link(text => 'View transactions');

    # should replace with better xpath
    $self->_waitForElement($agent, '/html/body/div[2]/div[2]/div[3]/div[2]/div[1]/div[1]/div[2]/div[2]/ul/li[1]/a/span[2]');
    $agent->follow_link(text => 'Export Statement Data');

    $self->_waitForElement($agent, '//*[@id="CSV"]');
    $agent->form_name('DownloadForm');
    $agent->field('Format' => 'CSV' );
    $agent->click({xpath => '//*[@id="CSV"]', synchronize => 0});
    sleep 1;
    $agent->click({xpath => '//*[@id="selectCard10"]', synchronize => 0});
    sleep 1;
    $agent->click({xpath => '//*[@id="radioid00"]', synchronize => 0});
    $agent->click({xpath => '//*[@id="radioid01"]', synchronize => 0});
    $agent->click({xpath => '//*[@id="radioid02"]', synchronize => 0});
    $agent->click({xpath => '//*[@id="radioid03"]', synchronize => 0});
     $agent->click({xpath => '//*[@id="myBlueButton1"]', synchronize => 0});
    sleep 15;

    # TODO: catch download properly
    $self->setFileName(csvFile);
    return $self->_loadCSVRows();

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

