#!/usr/bin/perl

package Loader_Danske;
use Moose;
extends 'Loader';

use strict;
use warnings;

use WWW::Mechanize::Firefox;
use HTML::Form;

use DataTypes::GenericRawLine;

use Switch;

has 'UserName' => ( is => 'rw', isa=>'Str', writer=>'setUserName', reader=>'getUserName');
has 'Password' => ( is => 'rw', isa=>'Str', writer=>'setPassword', reader=>'getPassword');
has 'AccountName' => ( is => 'rw', isa=>'Str', writer=>'setAccountName', reader=>'getAccountName');

# build string formats:
# file; directory, filename
# notfile; username; password, accountname
sub BUILD
{
    my ($self) = @_;
    my @buildParts = split (';' ,$self->build_string);
    if ($buildParts[0])
    {
        $self->setUserName($buildParts[1]);
        $self->setPassword($buildParts[2]);
        $self->setAccountName($buildParts[3]);
    }
    else
    {
        $self->setDirectory($buildParts[1]);
        $self->setFileName($buildParts[2]);
    }
}

sub _processExpenseLine
{
    my ($self, $agent) = @_;
    my $data = $agent->xpath('/html/body/div/form/table[2]', one=>1)->{innerHTML};
    my $line = GenericRawLine->new();
    my @matches = ($data =~ m/nytabellinie2\("(.*)", ?"(.*)", ?".*"\)/g);
    while (scalar @matches > 1)
    {
        my $key = shift @matches;
        my $value = shift @matches;
        $self->_addValueToLine($line, $key, $value);
    }
    print 'Saving ',$line->toString,"\n";
    return $line;
}

sub _addValueToLine
{
    my ($self, $line, $key, $value) = @_;
    switch ($key)
    {
        case 'Reference number:' { $line->setRefID($value) }
        case 'Reference:' { $line->setRefID($value) }
        case 'Text:' { $line->setDescription($self->_cleanText($value)) }
        case 'Amount:' { $value =~ s/,//g; $line->setAmount($value) }
        case 'Amount in DKK:' { $value =~ s/,//g; $line->setAmount($value) }
        case 'Date:' { $line->setTransactionDate($self->_formatDate($value)) }
        case 'Value date:' { $line->setProcessedDate($self->_formatDate($value)) }
        case 'Currency traded:' {$line->setFXCCY($value) }
        case 'Exchange rate:' { $value =~ s/ //g; $line->setFXRate($value) }
        case 'Amount in foreign currency:' { $value =~ s/ //g; $line->setFXAmount($value) }
        case 'Status:' { $line->setExtraText($line->getExtraText() . "\n" . 'Status: ' . $value)}
        case 'Message:' { $line->setExtraText($line->getExtraText() . $value . "\n") }

        # Fields for a transfer
        case 'Text on account statement:' { $line->setDescription($self->_cleanText($value)) }
        case 'Bank:' { $line->setExtraText($line->getExtraText() . "\n" . 'Bank: ' . $value)}
        case 'Type of payment:' { $line->setExtraText($line->getExtraText() . "\n" . 'Type of payment: ' . $value)}
        case 'From account:' { $line->setExtraText($line->getExtraText() . "\n" . 'From account: ' . $value)}
        case 'Payment reference:' { $line->setExtraText($line->getExtraText() . "\n" . 'Payment reference: ' . $value)}
        case 'To account:' { $line->setExtraText($line->getExtraText() . "\n" . 'To account: ' . $value)}
    }
}

sub _processInPageExpense
{
    my ($self, $agent) = @_;
    my $line = GenericRawLine->new();

    my $message = 0;
    for (my $i = 1; ;$i++)
    {
        my $table = '/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table/tbody/tr[2]/td/div[2]/table/tbody';
        my $type  = $table . "/tr[$i]/td[1]/span";
        my $value = $table . "/tr[$i]/td[2]/div/table/tbody/tr/td/span";

        last unless ( $agent->xpath($value, any=>1) );
        my $actualType = $agent->xpath($type, one=>1)->{innerHTML};
        my $actualValue = $agent->xpath($value, one=>1)->{innerHTML};

        if ($actualType ne '' )
        {
            $message = 0;
            $message = 1 if ($actualType eq 'Message:');

            $self->_addValueToLine($line, $actualType, $actualValue);
        }
        elsif ($message)
        {
            $actualType = 'Message:';
            $actualValue = $agent->xpath($value, one=>1)->{innerHTML};
            $self->_addValueToLine($line, $actualType, $actualValue);
        }
    }

    print 'Saving ',$line->toString,"\n";
    return $line;
}

sub _formatDate
{
    my ($self, $date) = @_;
    $date =~ m/([0-9]{2}).([0-9]{2}).([0-9]{4})/;
    return "$3-$2-$1";
}

sub _cleanText
{
    my ($self, $text) = @_;
    $text =~ s/&amp;/&/g;
    return $text;
}

###############################################################################
## All functions below for navigation of website ##############################
###############################################################################

sub _navToMainPage
{
    my ($self, $agent) = @_;
    $agent->get('https://www.danskebank.dk/en-dk/Personal/Pages/personal.aspx?secsystem=J2');

    ## setup user id
    die "Login button never appeared" 
        unless ($self->_waitForElement($agent, '/html/body/div[3]/div[2]/div/form/div/button[1]'));

    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/form/fieldset/input', $self->getUserName);
    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/form/fieldset/div/div/input', $self->getPassword);

    $agent->xpath('/html/body/div[3]/div[2]/div/form/div/button[1]', one=>1)->click;
    #TODO: check load happened correctly

    ## deal with nemid code..
    die "Nemid form never loaded" 
        unless ($self->_waitForElement($agent, '/html/body/div[3]/div[2]/div/div/form/fieldset/label'));
    my $numb = $agent->xpath('/html/body/div[3]/div[2]/div/div/form/fieldset/label', one=>1)->{textContent};
    print "\nType nemid for security number: $numb\n";
    chomp(my $secret = <>);

    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/div/form/fieldset/input', $secret);
    $agent->xpath('/html/body/div[3]/div[2]/div/div/form/div[8]/button[1]', one=>1)->click;

    die "Account link never loaded" 
        unless ($self->_waitForElement($agent, "//a[(text()=\"" . $self->getAccountName . "\")]"));
    $agent->follow_link(text => $self->getAccountName);

    sleep 30;
}

sub _pullOnlineData
{
    my $self = shift;
    my @result;

    my $agent = WWW::Mechanize::Firefox->new();
    $self->_navToMainPage($agent);

    my @loadedExpenses;

    # starting at 2, as first row is header
    for (my $i = 2; ;$i++)
    {
        # define xpaths
        my $expenseTable = "/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table";
        my $expenseRow = $expenseTable .    "/tbody/tr[$i]";
        my $reconciledBox = $expenseRow .   "/td[12]/div/input";
        my $categorisation = $expenseRow .  "/td[2]/div/div/a";
        my $expenseDetails = $expenseRow .  "/td[5]/div[1]/a";

        my @dataElements;
        $dataElements[0] = '/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table';
        $dataElements[1] = '/html/body/div/form/table[2]';

        # Wait until page fully loaded
        $self->_waitForElement($agent, $expenseTable);

        # if should process
        last unless ($agent->xpath($expenseRow, any=>1));
        next unless ($agent->xpath($expenseDetails, any=>1));
        my $processedEx = 0;
        if ( $agent->xpath($reconciledBox, any=>1) )
        {
            next if ( $agent->xpath($reconciledBox, one=>1)->{checked});
            $processedEx = 1;
        }

        # follow link
        $agent->xpath($expenseDetails, one=>1)->click;

        # process
        my $element = $self->_waitForElements($agent, \@dataElements);
        $agent->back() unless ($element);

        my $line;

        if ($element eq $dataElements[0])
        {
            $line = $self->_processInPageExpense($agent);
            unless ( $processedEx )
            {
                #$line->setAmount( $line->getAmount() * -1 );
                $line->setTemporary(1);
            }
        }
        else
        {
            $line = $self->_processExpenseLine($agent);
        }

        push (@result, $line->toString()) unless ($line->isEmpty());
        $agent->back();

        push (@loadedExpenses, $reconciledBox) if ($processedEx);
    }

    foreach my $reconBox (@loadedExpenses)
    {
        $self->_waitForElement($agent, $reconBox);
        $agent->xpath($reconBox, one=>1)->click;
        sleep 3;
    }

    return \@result;
}

sub _setAllValues
{
    my ($self, $agent, $xpath, $value) = @_;
    foreach ($agent->xpath($xpath))
    {
        $_->{value} = $value;
    }
}

1;

