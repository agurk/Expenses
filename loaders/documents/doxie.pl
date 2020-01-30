#!/usr/bin/perl

use strict;
use warnings;

use IO::Socket::Multicast;
use Net::DBus;

require HTTP::Request;
require LWP::UserAgent;

use HTTP::Request::Common;
use CACertOrg::CA;
$ENV{PERL_LWP_SSL_VERIFY_HOSTNAME} = 0;
use JSON;

my $MULTICAST_PORT = '1900';
my $MULTICAST_ADDR = '239.255.255.250';
my $MESSAGE_FREQUENCY= '60';

my $DOC_DIR = '/home/timothy/src/Expenses/f2/dist/resources/documents';

sub newDocument
{
    my %document;
    $document{'id'} = 0;
    $document{'filename'} = '';
    $document{'date'} = '';
    $document{'text'} = '';
    $document{'filesize'} = 0;
    return \%document;
}

sub sendLine
{
    my $line = shift;
    my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0, SSL_verify_mode => 0x00, SSL_ca_file      => CACertOrg::CA::SSL_ca_file() });
    my $json = JSON->new->allow_nonref;
    my $header = ['Content-Type' => 'application/json; charset=UTF-8'];
    my $url = 'https://localhost:8000/documents/';
    my $encoded_data = $json->encode($line);
    my $request = HTTP::Request->new('POST', $url, $header, $encoded_data);
    my $response = $ua->request($request);
    print ("Saving: $encoded_data\n");
    print ("Response: ", $response->code,"\n");
    #print ($response->message,"\n");
    return $response->code;
}

sub _get_image
{
	my ($name, $address, $user, $pass, $ua) = @_;
	my $request = HTTP::Request->new(GET => $address.'/scans/DOXIE/JPEG/' . $name);
	$request->authorization_basic($user, $pass);
    #my $ua = LWP::UserAgent->new;
	my $response = $ua->request($request);
	my $filename = "$name";
	print "writing to $filename\n";
	open (my $file, '>', $filename) or die "Cannot open $filename for writing\n";
	print $file $response->content;
    `convert $filename -resize 128x128 thumbs/$filename`;
	close ($file);
}

sub loadDocuments
{
    my ($address, $user, $pass) = @_;
	my $request = HTTP::Request->new(GET => $address .'/scans.json');
	$request->authorization_basic($user, $pass);

	my $ua = LWP::UserAgent->new();
	my $response = $ua->request($request);

	if ($response->content =~ m/^Can't connect to/)
	{
		print "Cannot connect to device\n";
		return;
	}
    print $response->content,"\n";
	my $scans = decode_json $response->content;
	foreach (@$scans)
	{
        $ua->timeout(2);
		print "\n\n";
        print $_->{name},"\n";
		print $_->{modified},"\n";
		print $_->{size},"\n";
		$_->{name} =~ m/([^\/]*)$/;
		my $name = $1;
        $name =~ m/([0-9]{4})/; 
        next if ($1 < 2800);

		chdir ($DOC_DIR);

            my $document =  newDocument();
            $document->{'filename'} = $name;
            $document->{'date'} = $_->{modified};
            $document->{'filesize'} = $_->{size};

        if ( -e $name ) {
			print "---> skipping: ",$_->{name},"\n";
        } else {
			_get_image($name, $address, $user, $pass, $ua);
            sendLine($document);
			print "---> saved: ",$document->{'filename'},"\n";
        }
	}
    
    # todo: reload failed docs
}

sub _get_device
{
	my ($data, $doxies) = @_;
	foreach (@{$doxies})
	{
		if ($data =~ m/$_->{id}/)
		{
			#print "Matched!\n";
			return $_;
		}
	}
	return 0;
}

sub listenMulticast
{ 
    my ($doxies) = @_;
	my $socket = IO::Socket::Multicast->new(
					LocalPort=>$MULTICAST_PORT,
					ReuseAddr=>1,
					ReusePort=>1,
	) or die $!;

	$socket->mcast_add($MULTICAST_ADDR) or die $!;

	my $msg;
	my $lastSentTime=0;
	while (1)
	{
		$socket->recv($msg, 4096);
		print "received:\n";
        print $msg;
		my $device = _get_device($msg, $doxies);
		if ($device and time >= ($lastSentTime + $MESSAGE_FREQUENCY))
		{
			print "found doxie!\n";
			$lastSentTime = time;
            $msg =~ m/Application-URL: (http:\/\/[0-9\.]*:8080)\/scans\//;
            loadDocuments($1, 'user', 'password');
		}
	}
}

my %d1 = (id => '123456789ABC', password => 'password');
my %d2 = (id => 'DEF123456789', password => 'password');
my @ds = (\%d1, \%d2);

#listenMulticast(\@ds);

loadDocuments("http://192.168.1.100:8080", 'user', 'password');
