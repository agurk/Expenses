#!/usr/bin/perl

use strict;
use warnings;

use HTTP::Request::Common;
use CACertOrg::CA;
$ENV{PERL_LWP_SSL_VERIFY_HOSTNAME} = 0;
use JSON;
use LWP::UserAgent;

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

sub truncDate
{
my $date = shift;
$date =~ s/T.*$//g;
return $date;
}

sub main
{
    my ($filename, $user, $account) = @_;
    open(my $file, '<', $filename) or die "Cannot open $filename\n";
    foreach (<$file>)
    {
        my $records = decode_json $_;
        foreach (@{$records->{'expenses'}})
        {
            next unless ($_->{'date'} =~ m/^2020/);
            next if ($_->{'deleted_at'});
            print ($_->{'id'},"\n");
            next unless ($_->{'repayments'}->[0]->{'from'} == $user);
            print $_->{'date'}, "\n";
            my $expense = newExpense($account ,$_->{'currency_code'});
            $expense->{'description'} = $_->{'description'};
            print $_->{'repayments'}->[0]->{'amount'}, "\n";
            my $amount = $_->{'repayments'}->[0]->{'amount'} * -1;
            $expense->{'amount'} = "$amount";
            $expense->{'transactionReference'} = "$_->{'id'}";
            $expense->{'date'}= $_->{'date'};
            #       $expense->{'date'} = truncDate($_->{'effectiveDate'});
            #       if  ($_->{'referenceNbr'}) {
            #           if ($_->{'referenceNbr'} eq "  000000000000000000000")
            #           {
            #               $expense->{'transactionReference'} = 'syntheticRef' . $expense->{'description'} . $expense->{'date'} . $expense->{'amount'};
            #           } else {
            #               $expense->{'transactionReference'} = $_->{'referenceNbr'};
            #           }
            #           $expense->{'processDate'} = truncDate($_->{'postDate'});
            #       } else {
            #           $expense->{'metadata'}{'temporary'} = $JSON::true;
            #       }
            #       if ($_->{'exchangeRate'} > 0) {
            #           $expense->{'fx'}{'rate'} = $_->{'exchangeRate'} + 0;
            #           $expense->{'fx'}{'amount'} = sprintf("%.2f", (($_->{'amount'}) / $_->{'exchangeRate'})) * -1;
            #       }
               print(sendLine($expense),"\n");
            }
    }
    close ($file);
}

main('data/splitwise.json', 1234567, 14);

