#!/usr/bin/perl

package Loader_AMEX;
use Moose;

use WWW::Mechanize;

extends 'Loader';

sub _loadCSVLine
{
    my ($self, $line) = @_;
    chomp($line);
    $line =~ s/\r//g;
    return if ($self->numbers_store()->isDupe($line));
    my @lineParts=split(/,/, $line);
    # skip payment, but have to leave negative number in case of refund
    my $classification = $self->getClassification($line);
    my @record = ($lineParts[1],$lineParts[0],$lineParts[2],$classification);
    # Value comes in quotes. Rediculous.
    $record[2] =~ s/\"//g;
    #$$DATA{$line} = \@record;
    $self->numbers_store()->addValue($line,\@record);
}

# This will try and open a file, if a file name has been specified in the
# Loader super object, if not it will try and do the online load
# File takes the csv format of:
# date, reference, amount, name, process date
sub load
{
    my $self = shift;
    my $DATA = $self->numbers_store()->data_list();
    if (defined $self->file_name())
    {
        open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
        foreach (<$file>)
        {
            _loadCSVLine($self, $_);
        }
        close($file);
    }
    else
    {
        foreach(@{$self->get_online_data()})
        {
        	_loadCSVLine($self, $_);
        }
    }
    $self->numbers_store()->save();
}

# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number
# selectradio => with the value set to the statement periods we want to download
sub get_online_data
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
    $agent->set_fields('Format' => 'CSV');# or die "fail\n";
    # Now we need to set which periods we want
    foreach (split('\n',$agent->content()))
    {
        # we want to find lines that match the following pattern:
        # <input id="radioid03"name="selectradio" type="checkbox"  title="Download Statement for  25 May 11 - 24 Jun 11 " value="20110525~20110624"/>
        # as these contain the value attribute that needs to be selected as part of the form
        if ($_ =~ m/.*selectradio.*value=\"(.*)\".*/)
        {
            $agent->tick('selectradio',$1);
        }
    }    
    # Now we set the card type
    $agent->set_fields('selectradio' => $self->settings->AMEX_CARD_NUMBER);
    $agent->submit();
    my @lines = split ("\n",$agent->content());
    return \@lines;
}

