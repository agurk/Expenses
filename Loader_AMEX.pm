#!/usr/bin/perl

package Loader_AMEX;
use Moose;

use WWW::Mechanize;

extends 'Loader';

has 'AMEX_PASSWORD' => ( is => 'rw', isa=>'Str', required => 1 );
has 'AMEX_USERNAME' => ( is => 'rw', isa=>'Str', required => 1);
has 'AMEX_CARD_NUMBER' => ( is => 'rw', isa=>'Str', required => 1 );
# Index is 0-rated
has 'AMEX_INDEX' => ( is => 'rw', isa=>'Str', required => 1);

use constant INPUT_LINE_PARTS_LENGTH => 4;

sub _makeRecord
{
    my ($self, $line) = @_;
    my @lineParts=split(/,/, $$line);
	die "wrong line length\n" unless (scalar @lineParts >= INPUT_LINE_PARTS_LENGTH);
    # Value comes in quotes. Ridiculous.
    $lineParts[2]  =~ s/\"//g;
    return Expense->new ( OriginalLine => $$line,
                          ExpenseDate => $lineParts[0],
                          ExpenseDescription => $lineParts[3],
                          ExpenseAmount => $lineParts[2],
                          AccountName => $self->account_name,
                        )
}


sub _ignoreYear
{
        my ($self, $record) = @_;
        return 0 unless (defined $self->settings->DATA_YEAR);
        $record->getExpenseDate =~ m/([0-9]{4}$)/;
        return 0 if ($1 eq $self->settings->DATA_YEAR);
        return 1;
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
    $agent->set_fields('UserID' => $self->AMEX_USERNAME);#; or die "can't fill username\n";
    $agent->set_fields('Password' => $self->AMEX_PASSWORD );
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
    if ($$numbersOnPage{$self->AMEX_CARD_NUMBER})
    {
        $agent->set_fields('selectradio' => $self->AMEX_CARD_NUMBER);
    } else {
        print "**Couldn't find card number ",$self->AMEX_CARD_NUMBER,". It might be:\n";
        foreach (keys %$numbersOnPage)
        {
            print "    ",$_,"\n";
        }
        return 0;
    }
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

