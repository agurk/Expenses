#!/usr/bin/perl

use strict;
use warnings;

use WWW::Mechanize;

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
    my $contents = 'b=' . $baseCCY;
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
    my @from = ('2020','01','01');
    my @to   = ('2020','12','31');
    my $ccy =   'GBP';
    my @ccys = ('EUR', 'USD', 'DKK');
    
    $agent->post('http://fx.sauder.ubc.ca/cgi/fxdata', Content=>makeContents(\@to, \@from, $ccy, \@ccys));
    processData($agent->content(), $ccy, \@ccys);

    @ccys = ('EUR', 'USD');
    $ccy =   'DKK';
    $agent->post('http://fx.sauder.ubc.ca/cgi/fxdata', Content=>makeContents(\@to, \@from, $ccy, \@ccys));
    processData($agent->content(), $ccy, \@ccys);
    
}

main();

