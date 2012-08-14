#!/usr/bin/perl

package Loader_AMEX;
use Moose;

use WWW::Mechanize;

extends 'Loader';

sub _makeRecord
{
    my ($self, $lineParts) = @_;
    my @record = ($$lineParts[3],$$lineParts[0],$$lineParts[2]);
    # Value comes in quotes. Ridiculous.
    $record[2] =~ s/\"//g;
    return \@record;
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
    $agent->get("https://www.americanexpress.com/uk/cardmember.shtml") or die "Can't load page\n";
    $agent->form_id("ssoform") or die "Can't get form\n";
    $agent->set_fields('UserID' => $self->settings->AMEX_USERNAME);#; or die "can't fill username\n";
    $agent->set_fields('Password' => $self->settings->AMEX_PASSWORD );
    $agent->submit() or die "can't login\n";
    $agent->follow_link( text_regex => qr/View Latest Transactions for British Airways American Express Credit Card/) or die "1\n";
    $agent->follow_link( text_regex => qr/Download statement data/) or die "1\n";
    $agent->form_name('DownloadForm') or die "patience\n";
    # set the download format
    $agent->set_fields('Format' => 'CSV');# or die "Can't set download format\n";
    # Now we need to set which periods we want
    foreach (split('\n',$agent->content()))
    {
        # we want to find lines that match the following pattern:
        # <input id="radioid03"name="selectradio" type="checkbox"  title="Download Statement for  25 May 11 - 24 Jun 11 " value="20110525~20110624"/>
        # as these contain the value attribute that needs to be selected as part of the form
        if ($_ =~ m/.*selectradio.*value=\"(.*)\".*/)
        {
            $agent->tick('selectradio',$1);# or die "Can't tick $1\n";
        }
    }    
    my $numbersOnPage = $self->_checkNumberOnPage($agent);
    if ($$numbersOnPage{$self->settings->AMEX_CARD_NUMBER})
    {
	$agent->set_fields('selectradio' => $self->settings->AMEX_CARD_NUMBER);
    } else {
	print "**Couldn't find card number ",$self->settings->AMEX_CARD_NUMBER,". It might be:\n";
	foreach (keys %$numbersOnPage)
	{
	    print "    ",$_,"\n";
	}
	return 0;
    }
    $agent->submit();
    my @lines = split ("\n",$agent->content());
    $self->set_input_data(\@lines);
    return 1;
}

sub _checkNumberOnPage
{
    my $self = shift;
    my $agent = shift;
    my %foundNumbers;
    foreach ( $agent->content() =~ m/([0-9]{10,})/ )
    {
	$foundNumbers{$_} = 1;
    }
    return \%foundNumbers;
}

1;

