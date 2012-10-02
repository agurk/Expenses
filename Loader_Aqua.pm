#!/usr/bin/perl

package Loader_Aqua;
use Moose;

use strict;
use warnings;

use WWW::Mechanize;
use HTTP::Cookies;

extends 'Loader';

sub loadInput
{
    my $self = shift;
    my @input_data;
    open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
    foreach (<$file>)
    {
	# All the data come on one line, so we're going to only load
	# those lines into the array
	#push (@input_data, $_) if ($_ =~ m/RECENT TRANSACTIONS/);
	if (($_ =~ m/RECENT TRANSACTIONS/) || ($_ =~ m/STATEMENT DATE/))
	{
	    push(  @input_data, @{$self->_getTransactionsFromLine($_)}  );
	}
    }
    close($file);
    $self->set_input_data(\@input_data);
}

sub _getTransactionsFromLine
{
    my ($self, $line) = @_;
    $line =~ s/<\/tr>/\n/g;
    my @lines = split("\n",$line);
    my @returnLines;
    my $count = 0;
    foreach(@lines)
    {
	$count+=1;
	next if ($count < 4 );
	next if ($_ =~ m/<\/table>/);
	# get rid of any commas before we make a csv file
	$_ =~ s/,/-/g;
	$_ =~ s/.*<tr><td class="date">//;
	$_ =~ s/<\/td><td class="date">/,/;
	$_ =~ s/<\/td><td class="description">/,/;
	$_ =~ s/<\/td><td class="description">/,/;
	$_ =~ s/<\/td><td class="amount">/,/;
	$_ =~ s/<\/td>//;
	# this sometimes gets entered into the description
	$_ =~ s/<br \/>/ /g;
	push (@returnLines, $_);
    }
    return \@returnLines;
}

sub _makeRecord
{
    my ($self, $lineParts) = @_;
    my @record = ($$lineParts[3],$$lineParts[0],$$lineParts[4]);
    $record[2] =~ s/£//g;
    if ($record[2] =~ m/DR/)
    {
	$record[2] =~ s/DR//;
    } else {
        # ASSUME CR!!
	$record[2] =~ s/CR//;
	$record[2] * -1;
    }
    return \@record;
}


# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number
# selectradio => with the value set to the statement periods we want to download
#sub _pullOnlineData
#{
#    my $self = shift;
#    my $agent = WWW::Mechanize->new(autocheck => 1);
#    $agent->agent_alias( 'Linux Mozilla' );
#    $agent->cookie_jar(HTTP::Cookies->new());
#   # $agent->get("https://195.171.220.59") or die "Can't load page\n";
#    $agent->get("http://www.aquacard.co.uk/") or die "Can't load page\n";
#    $agent->follow_link( n => 10) or die "1\n";
#    $agent->form_id("mainform") or die "Can't get form\n";
#    $agent->set_fields('datasource_7a325a0b-a613-4017-9f2b-abe99c1959db' => $self->settings->AQUA_USERNAME);
#    $agent->set_fields('datasource_0bde4679-4621-4b88-ab45-ebcc631fe471' => $self->settings->AQUA_PASSWORD );
#    $agent->submit() or die "can't login\n";
#    $agent->form_id('mainform');
#    $agent->set_fields('answerdatasource_8935509c-d0ae-4a4c-835d-0637283152b6' => $self->settings->AQUA_QUESTION);
#    $agent->submit();
#    
#    print $agent->content();

#
#
#$agent->follow_link( text_regex => qr/View Latest Transactions for British Airways American Express Credit Card/) or die "1\n";
#    $agent->follow_link( text_regex => qr/Download statement data/) or die "1\n";
#    $agent->form_name('DownloadForm') or die "patience\n";
#    # set the download format
#    $agent->set_fields('Format' => 'CSV');# or die "Can't set download format\n";
#    # Now we need to set which periods we want
#    foreach (split('\n',$agent->content()))
#    {
#        # we want to find lines that match the following pattern:
#        # <input id="radioid03"name="selectradio" type="checkbox"  title="Download Statement for  25 May 11 - 24 Jun 11 " value="20110525~20110624"/>
#        # as these contain the value attribute that needs to be selected as part of the form
#        if ($_ =~ m/.*selectradio.*value=\"(.*)\".*/)
#        {
#            $agent->tick('selectradio',$1);# or die "Can't tick $1\n";
#        }
#    }    
#    my $numbersOnPage = $self->_checkNumberOnPage($agent);
#    if ($$numbersOnPage{$self->settings->AMEX_CARD_NUMBER})
#    {
#	$agent->set_fields('selectradio' => $self->settings->AMEX_CARD_NUMBER);
#    } else {
#	print "**Couldn't find card number ",$self->settings->AMEX_CARD_NUMBER,". It might be:\n";
#	foreach (keys %$numbersOnPage)
#	{
#	    print "    ",$_,"\n";
#	}
#	return 0;
#    }
#    $agent->submit();
#    # Assume the download has failed if this string is in the results
#    if ($agent->content() =~ m/DownloadErrorPage/)
#    {
#	print " AMEX failed, retrying ";
#	return 0;
#    }
#    my @lines = split ("\n",$agent->content());
#    $self->set_input_data(\@lines);
#    return 1;
#}
#
#sub _checkNumberOnPage
#{
#    my $self = shift;
#    my $agent = shift;
#    my %foundNumbers;
#    foreach ( $agent->content() =~ m/([0-9]{10,})/ )
#    {
#	$foundNumbers{$_} = 1;
#    }
#    return \%foundNumbers;
#}

1;

