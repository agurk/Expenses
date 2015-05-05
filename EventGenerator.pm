#!/usr/bin/perl

use strict;
use warnings;

package EventGenerator;

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
		dbus_signal($EVENT_TYPE, ["string", ["dict", "string", "string"]]);
		dbus_method('sendEvent', ["string", ["dict", "string", "string"]], []);
#		dbus_signal($EVENT_TYPE, ["string"]);
#		dbus_method('sendEvent', ['string'],[]);
	return $self;
}


sub sendEvent
{
	my ($self, $messageType, $args) = @_;
	$self->emit_signal($EVENT_TYPE, $messageType, $args);
}


sub runGenerator
{
	my $bus=Net::DBus->session();
	my $service=$bus->export_service($DBUS_SERVICE_NAME);
	my $object=EventGenerator->new($service);
	print "Waiting for Events\n";
	Net::DBus::Reactor->main->run();
}

1;

