#!/usr/bin/perl

package Loader_Danske;
use Moose;
extends 'Loader';

use strict;
use warnings;

use WWW::Mechanize;
use HTTP::Cookies;

use DataTypes::GenericRawLine;

use Switch;

has 'directory'=>(is => 'rw', isa=>'Str', writer=>'setDirectory', reader=>'getDirectory',);

has 'input_data' => ( is => 'rw', isa => 'ArrayRef', writer=>'set_input_data', reader=>'get_input_data', default=> sub { my @empty; return \@empty});
has 'USER_NAME' => ( is => 'rw', isa=>'Str', writer=>'setUserName');
has 'SURNAME' => ( is => 'rw', isa=>'Str', writer=>'setSurname' );
has 'SECRET_WORD' => ( is => 'rw', isa=>'Str', writer=>'setSecretWord' );
has 'SECRET_NUMBERS' => ( is => 'rw', isa=>'Str', writer=>'setSecretNo' );
has 'cutoff_date' => (is => 'rw', isa=>'Str', writer=>'setCutoffDate', reader=>'getCutoffDate');
has 'process_statement' => (is => 'rw', isa=>'Bool', writer=>'setProcessStatement', reader=>'getProcessStatement');

use constant DATE_INDEX => 0;
use constant DESCRIPTION_INDEX => 2;
use constant AMOUNT_INDEX => 3;
use constant CREDIT_DEBIT_INDEX => 4;

# build string formats:
# file; filename
# notfile; username; surname; secretword; secretnumber; processStatement; (cutoff date)
sub BUILD
{
	my ($self) = @_;
	my @buildParts = split (';' ,$self->build_string);
	$self->setDirectory($buildParts[0]);
	$self->setFileName($buildParts[1]);
}

sub proccessFile
{
	my ($self, $filename) = @_;
	print "looking at $filename\n";
	chdir $self->getDirectory;
	open (my $file, '<', $filename) or die "Cannot open file $filename\n";
	my $line = GenericRawLine->new();
	foreach (<$file>)
	{
		if ($_ =~ m/nytabellinie2\("(.*)", ?"(.*)", ?"(.*)"\)/)
		{
			my ($key, $value) = ($1, $2);
			switch ($key)
			{
				case 'Reference number:' { $line->setRefID($value) }
				case 'Text:' { $line->setDescription($value) }
				case 'Amount:' { $line->setAmount($value) }
				case 'Date:' { $line->setTransactionDate($self->_formatDate($value)) }
				case 'Value date:' { $line->setProcessedDate($self->_formatDate($value)) }
				case 'Currency traded:' {$line->setFXCCY($value) }
				case 'Exchange rate:' { $value =~ s/ //g; $line->setFXRate($value) }
				case 'Amount in foreign currency:' { $value =~ s/ //g; $line->setFXAmount($value) }
			}
		}
	}
	print 'Saving ',$line->toString,"\n";
	if (defined $line->getRefID() and ! $line->getRefID() eq '')
	{
		$self->numbers_store()->addRawExpense($line->toString,$self->account_id());
	}
	unlink $filename or warn "could not delete $filename\n";
	close ($file);
}

sub readDirectory 
{
	my ($self) = @_;
	my $dir;
	opendir ($dir, $self->getDirectory) or die "Cannot open ",$self->getDirectory,"\n";
	my @files = readdir $dir;
	foreach (@files)
	{
		$self->proccessFile($_) if ($_ eq $self->file_name);
	}
	closedir $dir;
}

sub _formatDate
{
	my ($self, $date) = @_;
	$date =~ m/([0-9]{2}).([0-9]{2}).([0-9]{4})/;
	return "$3-$2-$1";
}

###############################################################################
## All functions below for navigation of website ##############################
###############################################################################

sub _pullOnlineData
{
	my $self = shift;
	my @result;
	my $agent = WWW::Mechanize->new( cookie_jar => {} );
	$agent->get('https://www.danskebank.dk/en-dk/Personal/Pages/personal.aspx?secsystem=J2');
	print $agent->content();
#	$agent->form_id('mainform');
#	$agent->set_fields('datasource_3a4651d1-b379-4f77-a6b1-5a1a4855a9fd' => $self->USER_NAME);
#	$agent->set_fields('datasource_2e7ba395-e972-4a25-8d19-9364a6f06132' => $self->SURNAME);
#	$agent->set_fields( '__EVENTTARGET' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
#	$agent->set_fields( '__EVENTARGUMENT' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
#	$agent->submit();
#
#	$agent->form_id('mainform');
#	$agent->set_fields('datasource_20056b2c-9455-4e42-aeba-9351afe0dbe1' => $self->SECRET_WORD );
#
#	my $secretNumbers = $self->_getPasscodes($agent);
#	$agent->set_fields('selectedvalue_dc28fef3-036f-4e48-98a0-586c9a4fbb3c' => $$secretNumbers[0]);
#	$agent->set_fields('selectedvalue_f758fdb6-4b2c-4272-b785-cb3989b67901' => $$secretNumbers[1]);
#	$agent->set_fields( '__EVENTTARGET' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
#	$agent->set_fields( '__EVENTARGUMENT' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
#	$agent->submit();
#
#	my $pageNumber = 0;
#	$pageNumber = -2 if ($self->_getPageNumber($agent) == -1);
#
#	while ($self->_getPageNumber($agent) > $pageNumber)
#	{
#		my @lines = split ("\n",$agent->content());
#		$self->_setOutputData(\@lines);
#		$pageNumber = $self->_getPageNumber($agent);
#		$agent->click_button( name => $self->_getNextPageLinkName($agent) ) unless ($pageNumber == -1 or ! defined $self->_getNextPageLinkName($agent));
#	}
#
#	$self->_doPostback($agent, 'View statements');
#
#	if ($self->getProcessStatement())
#	{
#		for (my $i=0; $i < 5; $i++)
#		{
#			$agent->post('https://service.aquacard.co.uk/aqua/web_channel/cards/servicing/youraccount/statement.aspx',
#				[	'__EVENTTARGET' => '',
#					'__EVENTARGUMENT' => '',
#					'__VIEWSTATE' => '',
#					'Target_f62e6d84-ef4f-4981-9b63-2c14b74ea065' => $self->_getViewState($agent),
#					'fixle' => '',
#					'selectedvalue_cfb4c82a-5338-4170-bafd-55c914904bcc' => $i,
#				]);
#
#			$self->_doPostback($agent, 'Transactions');
#			$pageNumber = 0;
#			$pageNumber = -2 if ($self->_getPageNumber($agent) == -1);
#
#			while ($self->_getPageNumber($agent) > $pageNumber)
#			{
#				my @lines = split ("\n",$agent->content());
#				$self->_setOutputData(\@lines);
#				if ($self->_getPageNumber($agent) > $pageNumber)
#				{
#					$pageNumber = $self->_getPageNumber($agent);
#					$agent->click_button( name => $self->_getNextPageLinkName($agent) ) if (defined $self->_getNextPageLinkName($agent));
#				}
#			}
#			$agent->get('/sitecore/content/Aqua/Web_Channel/Cards/Servicing/YourAccount/Statement.aspx');
#		}
#	}
#
#	return $self->_returnStrings();

}


1;

