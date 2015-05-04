#!/usr/bin/perl

use strict;
use warnings;

use Net::DBus;
use Net::DBus::Reactor;

use EventSettings;

sub handleMessage
{
	my ($message, $args) = @_;
	print "received $message:\n";
	if (defined $args)
	{
		foreach (keys %$args) {print "$_ ->	",$$args{$_},"\n"}
	}
}

sub main
{

	my $bus=Net::DBus->session();
	my $service=$bus->get_service($DBUS_SERVICE_NAME);
	my $object=$service->get_object($SERVICE_OBJECT_NAME, $DBUS_INTERFACE_NAME);
	
	
	$object->connect_to_signal($EVENT_TYPE, \&handleMessage);
	
	my $reactor=Net::DBus::Reactor->main();
	$reactor->run();
}

main();


