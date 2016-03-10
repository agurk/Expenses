#!/usr/bin/perl

package UDPListener;
use Moose;

use strict;
use warnings;

use IO::Socket::Multicast;
use Net::DBus;

use EventGenerator;
use EventSettings;

# array of hashrefs, each one representing a device. contains fields:
# id - unique id of doxie
# password
has 'Doxies' => ( 
		is => 'ro',
		isa => 'ArrayRef',
		reader=>'getDoxies',
		default=> sub { my @empty; return \@empty},
);

my $MULTICAST_PORT = '1900';
my $MULTICAST_ADDR = '239.255.255.250';

my $MESSAGE_FREQUENCY= '60';

sub _get_device
{
	my ($self, $data) = @_;
	foreach (@{$self->getDoxies})
	{
		if ($data =~ m/$_->{id}/)
		{
			#print "Matched!\n";
			return $_;
		}
	}
	return 0;
}

sub _document_args
{
	my ($self, $device, $data) = @_;
	my %args;
	$data =~ m/Application-URL: (http:\/\/[0-9\.]*:8080)\/scans\//;
	$args{'uri'} = $1 if ($1);
	$args{'password'} = $device->{'password'};
	$args{'time'} = time;
	return \%args;
}

sub listen
{ 
	my ($self) = @_;
	my $socket = IO::Socket::Multicast->new(
					LocalPort=>$MULTICAST_PORT,
					ReuseAddr=>1,
					ReusePort=>1,
	) or die $!;

	$socket->mcast_add($MULTICAST_ADDR) or die $!;

    my $bus=Net::DBus->session();
    my $service=$bus->get_service($DBUS_SERVICE_NAME);
    my $object=$service->get_object($SERVICE_OBJECT_NAME, $DBUS_INTERFACE_NAME);

	my $msg;
	my $lastSentTime=0;
	while (1)
	{
		$socket->recv($msg, 4096);
		#print "received:\n";
		my $device = $self->_get_device($msg);
		if ($device and time >= ($lastSentTime + $MESSAGE_FREQUENCY))
		{
			#print "found doxie!\n";
			$object->sendEvent('IMPORT_SCANS', $self->_document_args($device, $msg));
			my %empty;
			$object->sendEvent('PROCESS_SCANS', \%empty);
			$lastSentTime = time;
		}
	}
}

1;

