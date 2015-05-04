#!/usr/bin/perl

use strict;
use warnings;

use Cwd qw(abs_path getcwd);
BEGIN
{
    push (@INC, getcwd());
    no if $] >= 5.018, warnings => "experimental";
}

use Net::DBus;
use Net::DBus::Reactor;

use EventSettings;

use Try::Tiny;
use Switch;

use DataTypes::Expense;

use DocumentData::Loaders::Loader;
use DocumentData::Loaders::Loader_Doxie;
use DocumentData::Processors::Processor;

use Database::DAL;
use Database::DocumentsDB;

my $docProcessor = Processor->new();
my $documentDB = DocumentDB->new();
my $documentsDB = DocumentsDB->new();


sub handleMessage
{
	my ($message, $args) = @_;
	switch ($message) {
		case 'PROCESS_DOCUMENT' { $docProcessor->processDocument($$args{'did'}) }
		case 'DELETE_DOCUMENT' { _delete_document($args) }
		case 'IMPORT_SCANS' { Loader_Doxie->new()->loadDocument() }
		case 'PROCESS_SCANS' { _process_scans() }
	}
}

sub _process_scans
{
    foreach (@{$documentsDB->getUnclassifiedDocuments})
    {
		print "-> Found Scan: $_\n";
        $docProcessor->processDocument($_);
    }
}

sub _delete_document
{
	my ($args) = @_;
    my $document = $documentDB->getDocument($$args{'did'});
    $document->removeAllExpenseIDs();
    $document->setDeleted(1);
    $documentDB->saveDocument($document);

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


