#!/usr/bin/perl

use strict;
use warnings;

package EventReceiver;

use POSIX qw/strftime/;

use Net::DBus;
use Net::DBus::Reactor;

use EventSettings;

sub handleMessage
{
	my ($message, $args) = @_;
	my $now = strftime "%Y-%m-%d %H:%M:%S", localtime;
	print $now,", received: $message";
	if (keys %$args)
	{
		print ":\n";
		foreach (keys %$args) {print ' ' x 32,"$_ ->	",$$args{$_},"\n"}
	}
	else
	{
		print "\n";
	}
}

sub runReceiver
{
	my $bus=Net::DBus->session();
	my $service=$bus->get_service($DBUS_SERVICE_NAME);
	my $object=$service->get_object($SERVICE_OBJECT_NAME, $DBUS_INTERFACE_NAME);
	$object->connect_to_signal($EVENT_TYPE, \&handleMessage);
	print "Listening for Events\n";
	my $reactor=Net::DBus::Reactor->main();
	$reactor->run();
}

1;

