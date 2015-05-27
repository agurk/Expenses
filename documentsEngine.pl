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
use DocumentData::Loaders::Loader_File;
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
		case 'IMPORT_SCANS' { _import_scans($args) }
		case 'PROCESS_SCANS' { _process_scans() }
		case 'PIN_ITEM'	{ _pin_item($args) }
		case 'CONFIRM_DOC_EXPENSE'	{ $documentsDB->confirmDocEx($$args{'dmid'}) }
		case 'REMOVE_DOC_EXPENSE'	{ $documentsDB->removeDocEx($$args{'dmid'}) }
		case 'IMPORT_FILES'			{ _import_files($args)	}
		case 'RECLASSIFY_DOC'		{ $docProcessor->reclassifyDocument($$args{'did'}) }
	}
}

sub _import_files
{
	my ($args) = @_;
	my $scanner = Loader_File->new();
	$scanner->loadDocument($$args{'path'});
}

sub _import_scans
{
	my ($args) = @_;
	my $scanner = Loader_Doxie->new();
	$scanner->setAddress($$args{'uri'}) if (defined $$args{'uri'});
	$scanner->setPassword($$args{'password'}) if (defined $$args{'password'});
	$scanner->loadDocument();
}

sub _pin_item
{
	my ($args) = @_;
	if (defined $$args{'did'} and defined $$args{'eid'})
	{
		print 'Joining document: ',$$args{'did'},' with expense: ',$$args{'eid'},"\n";
		my $document = $documentDB->getDocument($$args{'did'});
		$document->addExpenseID($$args{'eid'});
		$documentDB->saveDocument($document);
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

