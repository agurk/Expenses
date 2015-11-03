#!/usr/bin/perl

use strict;
use warnings;

package EventSettings;
use parent 'Exporter';

our $EVENT_TYPE = 'ExpensesTimeEvent';
our $SERVICE_OBJECT_NAME = '/ExpensesTime/EventServiceObject';
our $DBUS_SERVICE_NAME = 'ExpensesTime.events';
our $DBUS_INTERFACE_NAME = 'ExpensesTime.interface';

our @EXPORT = qw($EVENT_TYPE $SERVICE_OBJECT_NAME $DBUS_SERVICE_NAME $DBUS_INTERFACE_NAME);

1;

