#!/usr/bin/perl

use strict;
use warnings;

use WWW::Mechanize;

sub get_input_value
{
	my $agent = shift;
	foreach (split('\n',$agent->content()))
	{
		if ($_ =~ m/<input name=\"(txt[0-9]*)\".*/)
		{
			#print $1,"\n";
			#$agent->set_fields($1 => '3017776400');
			return $1;
		}
	}
}

sub get_site
{
	my $agent = shift;
	foreach (split('\n',$agent->content()))
	{
		if ($_ =~ m/.*parent.location.href=\"(.*)\".*/)
		{
			#print $1,"\n";
			#$agent->set_fields($1 => '3017776400');
			chomp($1);
			return $1;
		}
	}	
}
sub main
{
	my $agent = WWW::Mechanize->new();
	#$agent->get("http://www.google.com");
	#$agent->get("http://www.nationwide.co.uk");
	$agent->get("https://olb2.nationet.com/signon/index2.asp");
	$agent->form_name("emptyform");
	$agent->submit();
	$agent->form_name('aspnetForm') or die "no form\n";
	print get_input_value($agent),"\n";
	$agent->set_fields(get_input_value($agent) => '3017776400');
#	$agent->set_visible('3017776400') or die "can't set number\n";
#	$agent->submit();
	$agent->click_button(name=>'ctl00$cphMainContent$btnSubmit');
	print "enter auth code:\n";
	my $foo = readline();
	chomp($foo);
	$agent->set_fields(get_input_value($agent) => $foo);
	$agent->click_button(name=>'ctl00$cphMainContent$btnSubmit');
#$agent->follow_link( text_regex => qr/or sign on with memorable data >>/) or die "can't get next page\n";
	$agent->click_button(value=>'Continue >>');
	open(my $file,'>','NATIONWIDE.html');
#	print "https://olb2.nationet.com/" . get_site($agent) ."\n";
	$agent->get("https://olb2.nationet.com/" . get_site($agent));
#	print $agent->uri();
	print "\n";
	print $file  $agent->content();
	close($file);
	#$agent->form_name("ssoform") or die "Can't get form\n";
}


# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number, hard coded so far....
# selectradio => with the value set to the statement periods we want to download
sub main2
{
	my $agent = WWW::Mechanize->new();
	$agent->get("https://www.americanexpress.com/uk/cardmember.shtml") or die "Can't load page\n";
	$agent->form_id("ssoform") or die "Can't get form\n";
	$agent->set_fields('UserID' => 'timothymollba');#; or die "can't fill username\n";
	$agent->set_fields('Password' => 'bb4pwev4' );
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
	$agent->set_fields('selectradio' => '376469512751002');
	$agent->submit();
	my @lines = split ("\n",$agent->content());
	return \@lines;
}

main();


