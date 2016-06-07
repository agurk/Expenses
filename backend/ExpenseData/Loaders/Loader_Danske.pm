#!/usr/bin/perl

package Loader_Danske;
use Moose;
extends 'Loader';

use strict;
use warnings;

use WWW::Mechanize::Firefox;
use HTML::Form;

use DataTypes::GenericRawLine;

use Switch;

has 'UserName' => ( is => 'rw', isa=>'Str', writer=>'setUserName', reader=>'getUserName');
has 'Password' => ( is => 'rw', isa=>'Str', writer=>'setPassword', reader=>'getPassword');
has 'AccountName' => ( is => 'rw', isa=>'Str', writer=>'setAccountName', reader=>'getAccountName');

# build string formats:
# file; directory, filename
# notfile; username; password, accountname
sub BUILD
{
	my ($self) = @_;
	my @buildParts = split (';' ,$self->build_string);
	if ($buildParts[0])
	{
		$self->setUserName($buildParts[1]);
		$self->setPassword($buildParts[2]);
		$self->setAccountName($buildParts[3]);
	}
	else
	{
		$self->setDirectory($buildParts[1]);
		$self->setFileName($buildParts[2]);
	}
}

sub _processExpenseLine
{
	my ($self, $data) = @_;
	#print "looking at $filename\n";
	#chdir $self->getDirectory;
	#open (my $file, '<', $filename) or die "Cannot open file $filename\n";
	my $line = GenericRawLine->new();
	my @matches = ($data =~ m/nytabellinie2\("(.*)", ?"(.*)", ?".*"\)/g);
	while (scalar @matches > 1)
	{
		my $key = shift @matches;
		my $value = shift @matches;
		switch ($key)
		{
			case 'Reference number:' { $line->setRefID($value) }
			case 'Text:' { $line->setDescription($value) }
			case 'Amount:' { $value =~ s/,//g; $line->setAmount($value) }
			case 'Date:' { $line->setTransactionDate($self->_formatDate($value)) }
			case 'Value date:' { $line->setProcessedDate($self->_formatDate($value)) }
			case 'Currency traded:' {$line->setFXCCY($value) }
			case 'Exchange rate:' { $value =~ s/ //g; $line->setFXRate($value) }
			case 'Amount in foreign currency:' { $value =~ s/ //g; $line->setFXAmount($value) }
		}
	}
	print 'Saving ',$line->toString,"\n";
	#if (defined $line->getRefID() and ! $line->getRefID() eq '')
	#{
#		$self->numbers_store()->addRawExpense($line->toString,$self->account_id());
#	}
#	unlink $filename or warn "could not delete $filename\n";
#	close ($file);
	return $line;
}

sub _processUnconfirmedExpense
{
	my ($self, $agent) = @_;
	my $line = GenericRawLine->new();
	$line->setDescription($agent->xpath('/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table/tbody/tr[2]/td/div[2]/table/tbody/tr[4]/td[2]/div/table/tbody/tr/td/span', one=>1)->{innerHTML});
    my $amount = $agent->xpath('/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table/tbody/tr[2]/td/div[2]/table/tbody/tr[5]/td[2]/div/table/tbody/tr/td/span', one=>1)->{innerHTML};
    $amount =~ s/,//g;
	$line->setAmount($amount * -1);
	$line->setTransactionDate($self->_formatDate($agent->xpath('/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table/tbody/tr[2]/td/div[2]/table/tbody/tr[6]/td[2]/div/table/tbody/tr/td/span', one=>1)->{innerHTML}));
	$line->setTemporary(1);
	print 'Saving ',$line->toString,"\n";
	return $line;
}

#sub readDirectory
#{
#	my ($self) = @_;
#	my $dir;
#	opendir ($dir, $self->getDirectory) or die "Cannot open ",$self->getDirectory,"\n";
#	my @files = readdir $dir;
#	foreach (@files)
#	{
#		$self->proccessFile($_) if ($_ eq $self->file_name);
#	}
#	closedir $dir;
#}

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

    my $agent = WWW::Mechanize::Firefox->new();
    $agent->get('https://www.danskebank.dk/en-dk/Personal/Pages/personal.aspx?secsystem=J2');

    ## setup user id
	die "Login button never appeared" 
		unless ($self->_waitForElement($agent, '/html/body/div[3]/div[2]/div/form/div/button[1]'));

    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/form/fieldset/input', $self->getUserName);
    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/form/fieldset/div/div/input', $self->getPassword);

	$agent->xpath('/html/body/div[3]/div[2]/div/form/div/button[1]', one=>1)->click;
	#TODO: check load happened correctly

    ## deal with nemid code..
	die "Nemid form never loaded" 
		unless ($self->_waitForElement($agent, '/html/body/div[3]/div[2]/div/div/form/fieldset/label'));
    my $numb = $agent->xpath('/html/body/div[3]/div[2]/div/div/form/fieldset/label', one=>1)->{textContent};
	print "\nType nemid for security number: $numb\n";
	chomp(my $secret = <>);

    $self->_setAllValues($agent, '/html/body/div[3]/div[2]/div/div/form/fieldset/input', $secret);
    $agent->xpath('/html/body/div[3]/div[2]/div/div/form/div[8]/button[1]', one=>1)->click;

	die "Account link never loaded" 
		unless ($self->_waitForElement($agent, "//a[(text()=\"" . $self->getAccountName . "\")]"));
    $agent->follow_link(text => $self->getAccountName);

	sleep 30;

    # starting at 2, as first row is header
    for (my $i = 2; ;$i++)
    {
		# Wait until page fully loaded
		$self->_waitForElement($agent, '/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table');

		#xpath is for expense row in document
		last unless ($agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]", any=>1));

        # if should process
		#td12 contains the reconcile box
		#td2 is the 'more details' link for the expense
		my $processedEx = 0;
        if ( $agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[12]/div/input", any=>1) ) {
			next if ( $agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[12]/div/input", one=>1)->{checked} );
			next if ($agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[2]/div/div/a", one=>1)->{innerHTML} eq 'Categorise');
			$processedEx = 1;
		} else {
			next unless ($agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[5]/div[1]/a", any=>1));
		}

        # follow link
        $agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[5]/div[1]/a", one=>1)->click;


        # process
        
		my $dataElement = '/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div/div/div/div[2]/div/table';
		$dataElement = '/html/body/div/form/table[2]' if ($processedEx);

		if ($self->_waitForElement($agent, $dataElement))
		{
			my $line;
			if ($processedEx) {
				$line = $self->_processExpenseLine($agent->xpath('/html/body/div/form/table[2]', one=>1)->{innerHTML});
			} else {
				$line = $self->_processUnconfirmedExpense($agent);
			}
			push (@result, $line->toString()) unless ($line->isEmpty()); 
			$agent->back();
			if ($processedEx) {
				$self->_waitForElement($agent, "/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[12]/div/input");
				$agent->xpath("/html/body/form/div[4]/div[3]/div/div/div[1]/div[3]/div[4]/div[1]/table/tbody/tr[$i]/td[12]/div/input", one=>1)->click;
				sleep 5;
			}
		} else {
			$agent->back();
		}
    }

	return \@result;
}

sub _setAllValues
{
    my ($self, $agent, $xpath, $value) = @_;
    foreach ($agent->xpath($xpath))
    {
        $_->{value} = $value;
    }
}

1;

