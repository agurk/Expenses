#!/usr/bin/perl

package Loader_AMEX;
use Moose;
extends 'Loader';

use WWW::Mechanize;

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

# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number
# selectradio => with the value set to the statement periods we want to download
sub _pullOnlineData
{
    my $self = shift;
    my $agent = WWW::Mechanize->new();

$agent->add_handler("request_send", sub { shift->dump; return });
$agent->add_handler("response_done", sub { shift->dump; return });
    $agent->get("https://www.americanexpress.com/uk/cardmember.shtml") or die "Can't load page\n";
    $agent->form_id("ssoform") or die "Can't get form\n";
    $agent->set_fields('UserID' => $self->AMEX_USERNAME);
    $agent->set_fields('USERID' => $self->AMEX_USERNAME);
    $agent->set_fields('Password' => $self->AMEX_PASSWORD );
    $agent->set_fields('PWD' => $self->AMEX_PASSWORD );
	$agent->set_fields('TARGET' => 'https://global.americanexpress.com/myca/intl/acctsumm/emea/accountSummary.do?request_type=&Face=en_GB&linknav=UK-Home-page-Myca-Login-Large');
	$agent->add_header( Host => 'global.americanexpress.com');
	$agent->add_header( Connection => 'keep-alive');
	$agent->add_header( Accept => 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8');
	$agent->add_header( 'Accept-Language' => 'en-GB,en;q=0.5');
	$agent->add_header( 'User-Agent' => 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:40.0) Gecko/20100101 Firefox/40.0');
    $agent->submit() or die "can't login\n";
#    $agent->follow_link(text => 'View Latest Transactions', n => $self->AMEX_INDEX+1) or die "1\n";
    $agent->follow_link(text => 'Export Statement Data');
#    $agent->follow_link( text_regex => qr/Download statement data/) or die "1\n";
    $agent->form_name('DownloadForm');
    # set the download format
    $agent->set_fields('Format' => 'CSV');# or die "Can't set download format\n";
    # Now we need to set which periods we want
    foreach (split('\n',$agent->content()))
    {
        # we want to find lines that match the following pattern:
        # <input id="radioid03"name="selectradio" type="checkbox"  title="Download Statement for  25 May 11 - 24 Jun 11 " value="20110525~20110624"/>
        # as these contain the value attribute that needs to be selected as part of the form
        if ($_ =~ m/id=\"radioid([0-9]).*selectradio.*value=\"(.*)\".*/)
        {
            $agent->tick('selectradio',$2) if ($1 == $self->AMEX_INDEX);
        }
    }    
    my $numbersOnPage = $self->_checkNumberOnPage($agent);
		open(my $file, '>', 'output.html');
		print $file $agent->content();
		close ($file);
#    if ($$numbersOnPage{$self->AMEX_CARD_NUMBER})
#    {
#        $agent->set_fields('selectradio' => $self->AMEX_CARD_NUMBER);
#    } else {
#        print "**Couldn't find card number ",$self->AMEX_CARD_NUMBER,". It might be:\n";
#        foreach (keys %$numbersOnPage)
#        {
#            print "    ",$_,"\n";
#        }
#        return 0;
#    }
#	$agent->tick('radioid00');
	$agent->field('dowloadFormat' => 'on' );
	$agent->field('Format' => 'CSV' );
	$agent->field('selectCard10' => 'on' );
	$agent->field('radioid00' => 'on' );
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

