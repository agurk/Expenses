#!/usr/bin/perl

use strict;
use warnings;

use WWW::Mechanize;
use JSON;
use LWP::UserAgent;
use LWP::Protocol::https;

sub insertLine
{
	my ($date, $ccy1, $ccy2, $amount) = @_;
    #print 'insert into _FXRates(date, ccy1, ccy2, rate) values (\'',$date,'\',\'',$ccy1,'\',\'',$ccy2,'\',',$amount,");\n";
    my $insertString = "insert into _FXRates(date, ccy1, ccy2, rate) values ('$date','$ccy1','$ccy2',$amount)";
    print $insertString,"\n";
    `sqlite3 /home/timothy/src/Expenses/expenses.db "$insertString"`
}

sub formatDate
{
	my ($date) = @_;
	$date =~ s/"//g;
	my @dateParts = split(/\//, $date);
	return $dateParts[0] . '-' . $dateParts[1] . '-' . $dateParts[2]
}

sub processData
{
    my ($input, $baseCCY, $CCYs) = @_;
    open (my $file, '<', \$input);
    <$file>;
    <$file>;
    foreach (<$file>)
    {
        next if ($_ =~ m/^"\(C\)/);
    	chomp;
    	my @lineParts = split (/,/, $_);
    	my $date = formatDate($lineParts[1]);
        my $i = 3;
        foreach (@$CCYs)
        {
    	    insertLine($date, $baseCCY, $_, $lineParts[$i]);
            $i++;
        }
    }
    close($file);
}

sub makeContents
{
    my ($to, $from, $baseCCY, $ccys) = @_;
    my $contents = '?b=' . $baseCCY;
    foreach (@$ccys)
    {
        $contents .= '&c=';
        $contents .= $_;
    }
    $contents .= '&rd=&fd='.$$from[2].'&fm='.$$from[1].'&fy='.$$from[0].'&ld='.$$to[2].'&lm='.$$to[1].'&ly='.$$to[0].'&y=daily&q=volume&f=csv&o=';
    return $contents;
}

sub main
{
    my $agent = WWW::Mechanize->new();
    my @from = ('2022','01','01');
    my @to   = ('2022','12','31');
    my $ccy =   'GBP';
    my @ccys = ('EUR', 'USD', 'DKK');
    
    $agent->get('http://fx.sauder.ubc.ca/cgi/fxdata' . makeContents(\@to, \@from, $ccy, \@ccys));
    processData($agent->content(), $ccy, \@ccys);

    @ccys = ('EUR', 'USD');
    $ccy =   'DKK';
    $agent->get('http://fx.sauder.ubc.ca/cgi/fxdata' . makeContents(\@to, \@from, $ccy, \@ccys));
    processData($agent->content(), $ccy, \@ccys);
    
    @ccys = ('USD');
    $ccy =   'EUR';
    $agent->get('http://fx.sauder.ubc.ca/cgi/fxdata' . makeContents(\@to, \@from, $ccy, \@ccys));
    processData($agent->content(), $ccy, \@ccys);
}

sub crypto
{
    my ($crypto, $ccy) = @_;
    my $agent = WWW::Mechanize->new();
    my $ua = LWP::UserAgent->new;
    $ua->protocols_allowed(['https']);
    my $key = 'CXZ78N3U6VY01PFF';
    my $series = "Time Series (Digital Currency Daily)";
    my $valueTag = "4a. close (USD)";
    
    $agent->get("https://www.alphavantage.co/query?function=DIGITAL_CURRENCY_DAILY&symbol=$crypto&market=$ccy&apikey=$key");
    my $data = decode_json $agent->content();
    foreach (keys %{$data->{$series}})
    {
        print($_,' - ', $data->{$series}->{$_}->{$valueTag}, "\n");
        my ($ticker, $date, $price, $currency) = @_;
        insertLine($_, $crypto, $ccy, $data->{$series}->{$_}->{$valueTag});
    #    insertLine($name, $_, $data->{$series}->{$_}->{$valueTag}, "USD");
    }
}

main();
#crypto('BTC', 'DKK');
#crypto('BTC', 'GBP');
crypto('BTC', 'USD');
