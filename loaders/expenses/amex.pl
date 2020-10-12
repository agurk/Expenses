#!/usr/bin/perl

use strict;
use warnings;

use HTTP::Request::Common;
use CACertOrg::CA;
$ENV{PERL_LWP_SSL_VERIFY_HOSTNAME} = 0;
use JSON;
use LWP::UserAgent;

use Spreadsheet::XLSX;
use Text::Iconv;

use Switch;

sub sendLine
{
    my $line = shift;
    my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0, SSL_verify_mode => 0x00, SSL_ca_file      => CACertOrg::CA::SSL_ca_file() });
    my $json = JSON->new->allow_nonref;
    my $header = ['Content-Type' => 'application/json; charset=UTF-8'];
    my $url = 'https://debian.home:8000/expenses/';
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

sub _getAmount
{
    my ($amount) = @_;
    my $returnAmount;
    if ($amount =~ m/-(.*)/)
    {
        $returnAmount = $1
    }
    else
    {
        $returnAmount = "-$amount";
    }
    return $returnAmount;
}

sub main
{
    my ($filename, $account, $ccy) = @_;

    my $converter = Text::Iconv -> new("utf-8", "windows-1251");
    my $spreadsheet = Spreadsheet::XLSX -> new ($filename, $converter);

    my $sheet = @{$spreadsheet -> {Worksheet}}[0];
    foreach my $row (7 .. $sheet -> {MaxRow}) {

        my $expense = newExpense($account, $ccy);
        $sheet -> {Cells} [$row] [0] -> { Val } =~ s/^([0-9]{2})\/([0-9]{2})\/([0-9]{4})//;
        $expense->{'date'} = $3.'-'.$2.'-'.$1;
        $expense->{'description'} = $sheet -> { Cells } [$row] [1] -> { Val };
        $expense->{'detailedDescription'} = $sheet -> { Cells } [$row] [3] -> { Val };
        $expense->{'amount'} = _getAmount($sheet -> { Cells } [$row] [2] -> { Val });
        $expense->{'transactionReference'} = $sheet -> { Cells } [$row] [9] -> { Val };

        my $response = sendLine($expense);
        unless ($response eq "200")
        {
            print "Could not save line ",$expense,"\n";
        }
}

#
#    open(my $file, '<', $filename);
#    foreach my $line (<$file>)
#    {
#        my $expense = newExpense($account, $ccy);
#        # remove first column as not quoted
#        $line =~ s/^([0-9]{2})\/([0-9]{2})\/([0-9]{4}),//;
#        $expense->{'date'} = $3.'-'.$2.'-'.$1;
#
#        my @lineParts=split(/","/, $line);
#
#        $lineParts[0] =~ m/Reference: ([A-Z0-9]*)/;
#        $expense->{'transactionReference'} = $1;
#
#        $lineParts[1]  =~ s/\"//g;
#        $lineParts[1]  =~ s/ //g;
#        $expense->{'amount'} = _getAmount($lineParts[1]) ;
#        
#        $lineParts[2]  =~ s/\"//g;
#        $expense->{'description'} = $lineParts[2];
#
#        $lineParts[3] =~ s/\"//g;
#        if ($lineParts[3]=~ m/^([0-9.,]{1,}) ([A-Z]{3}).* Currency Conversion Rate ([0-9.,]{1,}) Commission Amount ([0-9,.]*[0-9])/)
#        {
#            my ($amt, $ccy, $rate, $com) = ($1, $2, $3, $4);
#            $amt =~ s/,//;
#            $expense->{'fx'}->{'amount'} = $amt +0;
#            $expense->{'fx'}->{'currency'} = $ccy;
#            $expense->{'fx'}->{'rate'} = $rate +0;
#            $expense->{'commission'} = $com ;
#        } 
#        elsif ($lineParts[3] =~ m/([0-9.]{1,})  *([A-Z]{3})/)
#        {
#            $expense->{'fx'}->{'amount'} = $1 +0;
#            $expense->{'fx'}->{'currency'} = $2;
#        }
#
#        $lineParts[3] =~ m/Process Date ([0-9]{2})\/([0-9]{2})\/([0-9]{4})/;
#        $expense->{'processDate'} = $3.'-'.$2.'-'.$1;
#        
#        $expense->{'detailedDescription'} = $lineParts[3];
#
#        my $response = sendLine($expense);
#        unless ($response eq "200")
#        {
#            print "Could not save line ",$expense,"\n";
#        }
#    }
#    close($file)
}

main('data/activity.xlsx', 1, 'GBP');

