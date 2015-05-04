#!/usr/bin/perl

use strict;
use warnings;

package EventServer;

use EventSettings;

use Net::DBus;
use Net::DBus::Service;
use Net::DBus::Reactor;

use base qw(Net::DBus::Object);
use Net::DBus::Exporter qw(ExpensesTime.interface);

sub new
{
	my ($class, $service) = @_;
	my $self=$class->SUPER::new($service, $SERVICE_OBJECT_NAME);
	bless $self, $class;
	return $self;
}

dbus_signal($EVENT_TYPE, ["string", ["dict", "string", "string"]]);
dbus_signal($EVENT_TYPE, ["string"]);

sub sendEvent
{
	my ($self, $messageType, $args) = @_;
	print "sending: $messageType\n";
	if (defined $args)
	{
		$self->emit_signal($EVENT_TYPE, $messageType);
	}
	else
	{
		$self->emit_signal($EVENT_TYPE, $messageType, $args);
	}
}

dbus_method('sendEvent', ["string", ["dict", "string", "string"]], []);
dbus_method('sendEvent', ['string'],[]);

package main;
use EventSettings;

main
{
	my $bus=Net::DBus->session();
	my $service=$bus->export_service($DBUS_SERVICE_NAME);
	my $object=EventServer->new($service);
	Net::DBus::Reactor->main->run();
}

main();

