#!/usr/bin/perl

use strict;
use warnings;

use XML::LibXML;
use HTTP::Request::Common;
use CACertOrg::CA;
$ENV{PERL_LWP_SSL_VERIFY_HOSTNAME} = 0;
use JSON;
use LWP::UserAgent;

use Switch;

sub sendLine
{
    my $line = shift;
    my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0, SSL_verify_mode => 0x00, SSL_ca_file      => CACertOrg::CA::SSL_ca_file() });
    my $json = JSON->new->allow_nonref;
    my $header = ['Content-Type' => 'application/json; charset=UTF-8'];
    my $url = 'https://localhost:8000/expenses/';
    my $encoded_data = $json->encode($line);
    my $request = HTTP::Request->new('POST', $url, $header, $encoded_data);
    my $response = $ua->request($request);
    print ("Saving: $encoded_data\n");
    print ("Response: ", $response->code,"\n");
    #print ($response->message,"\n");
    return $response->code;
}

sub newExpense
{
    my ($accountId, $currency) = @_;
    my %expense;
    $expense{'id'} = 0;
    $expense{'transactionReference'} = '';
    $expense{'description'} = '';
    $expense{'detailedDescription'} = '';
    $expense{'accountId'} = $accountId;
    $expense{'date'} = '';
    $expense{'processDate'} = '';
    $expense{'currency'} = $currency;
    $expense{'commission'} = '0';
    $expense{'metadata'} = {};
    $expense{'metadata'}{'temporary'} = $JSON::false;
    $expense{'fx'} = {};
    $expense{'fx'}{'amount'} = 0;
    $expense{'fx'}{'currency'} = '';
    $expense{'fx'}{'rate'} = 0;
    return \%expense;
}

sub formatDate
{
    my ($date) = @_;
    $date =~ m/([0-9]{4})([0-9]{2})([0-9]{2})/;
    return "$1-$2-$3";
}
sub main
{
    my ($filename, $account, $ccy) = @_;

    my $parser = XML::LibXML->new();
    my $xmldoc = $parser->parse_file($filename);

    for my $sample ($xmldoc->findnodes('/OFX/BANKMSGSRSV1/STMTTRNRS/STMTRS/BANKTRANLIST/STMTTRN')) {
        my $expense = newExpense($account, $ccy);
        for my $property ($sample->findnodes('./*')) {
            switch ($property->nodeName())
            {
                case 'TRNAMT' { $expense->{'amount'} = $property->textContent()  }
                case 'NAME'   { $expense->{'description'} = $property->textContent() }
                case 'FITID'   { $expense->{'transactionReference'} = $property->textContent() }
                case 'DTPOSTED'   { $expense->{'date'} = formatDate($property->textContent()) }
                case 'TRNTYPE'   { $expense->{'detailedDescription'} = 'Transaction type: ' . $property->textContent() }
            }
            $expense->{'processDate'} = $expense->{'date'};
        }
        my $response = sendLine($expense);
        unless ($response eq "200")
        {
            print "Could not save line ",$expense,"\n";
        }
    }
}

main('data/nw.ofx', 3, 'GBP');

