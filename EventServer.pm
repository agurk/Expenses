#!/usr/bin/perl

use strict;
use warnings;

package EventServer;

use EventSettings;
use EventGenerator;
use EventReceiver;

sub main
{
	my $pid = fork();
	EventGenerator::runGenerator(0) if ($pid == 0);
	sleep 2;
	EventReceiver::runReceiver() unless ($pid == 0);
}

main();

