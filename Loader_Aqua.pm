#!/usr/bin/perl

package Loader_Aqua;
use Moose;

use strict;
use warnings;

use WWW::Mechanize;
use HTTP::Cookies;

use constant INTERNAL_FIELD_SEPARATOR => '|';

extends 'Loader';

has 'USER_NAME' => ( is => 'rw', isa=>'Str' );
has 'SURNAME' => ( is => 'rw', isa=>'Str' );
has 'SECRET_WORD' => ( is => 'rw', isa=>'Str' );
has 'SECRET_NUMBERS' => ( is => 'rw', isa=>'Str' );

#sub loadInput
#{
#    my $self = shift;
#    my @input_data;
#    open(my $file,"<",$self->file_name()) or warn "Cannot open: ",$self->file_name(),"\n";
#    foreach (<$file>)
#    {
#	next unless ($_ =~ m/abbr/);
#	chomp;
#	$_ =~ s/^[\w\t ]+//g;
#	$_ =~ s/^<td>//g;
#	$_ =~ s/<\/td><td>/INTERNAL_FIELD_SEPARATOR/g;
#	$_ =~ s/<\/td><td class="right">/INTERNAL_FIELD_SEPARATOR/g;
#	$_ =~ s/<abbr title="[^\"]*">/INTERNAL_FIELD_SEPARATOR/g;
#	$_ =~ s/<\/abbr><\/td>//g;
#	push(  @input_data, $_ );
#    }
#    close($file);
#    $self->set_input_data(\@input_data);
#}


#sub _getTransactionsFromLine
#{
#    my ($self, $line) = @_;
#    my @lines = split("\n",$line);
#    my @returnLines;
#    my $count = 0;
#    foreach(@lines)
#    {
#	print $_;
#	chomp;
#	next unless ($_ =~ m/.*attr.*/);
#	$_ =~ s/^\w+//g;
#	foreach ($_ =~ m/<br>(.*?)<\/br>/)
#	{
#	    print $_,"\n";
#	}
#	push (@returnLines, $_);
#    }
#    return \@returnLines;
#}

sub _makeRecord
{
    my ($self, $line) = @_;
    my @lineParts=split(/INTERNAL_FIELD_SEPARATOR/, $$line);
    $lineParts[3] =~ s/^[^0-9]*//;
    if ($lineParts[4] =~ m/DR/)
    {
		$lineParts[4] =~ s/DR//;
    } else {
		$lineParts[4] =~ s/CR//;
		$lineParts[4] *= -1;
    }
	return Expense->new (	OriginalLine => $$line,
							ExpenseDate => $lineParts[0],
							ExpenseDescription => $lineParts[2],
							ExpenseAmount => $lineParts[3],
							AccountName => $self->account_name,
						)
}

sub _ignoreYear
{
	my ($self, $record) = @_;
	return 0 unless (defined $self->settings->DATA_YEAR);
	$record->getExpenseDate =~ m/([0-9]{2}$)/;
	my $found_year = $1;
	$self->settings->DATA_YEAR =~ m/([0-9]{2}$)/;
	return 0 if ($found_year eq $1);
	return 1;
}

sub _generateSecretNumbers
{
    my $self = shift;
    # start with 0 so we can use 1 for 1 array referencing
    # i.e. first number (we care about) is in array posn 1
    my @numbers = (0);
    push (@numbers, split("", $self->SECRET_NUMBERS));
    return \@numbers;
}

sub _getPasscodes
{
    my ($self, $agent) = @_;
    #my $agent = shift;
    my @returnCodes;
    my $values = $self->_generateSecretNumbers();
    $agent->content() =~ m/([0-9]).. number of your Passcode.*([0-9]).. number of your Passcode/s;
    #print $1," next ",$2,"\n";
    $returnCodes[0] = $$values[$1];
    $returnCodes[1] = $$values[$2];
    return \@returnCodes;
}

sub _useInputLine
{
    my ($self, $line) = @_;
    return 1 if ($_ =~ m/abbr/);
    return 0;
}

sub _procesInputLine
{
    my ($self, $line) = @_;
    chomp $line;
    $line =~ s/^[\w\t ]+//g;
    $line =~ s/^<td>//g;
    $line =~ s/<\/td><td>/INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<\/td><td class="right">/INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<abbr title="[^\"]*">/INTERNAL_FIELD_SEPARATOR/g;
    $line =~ s/<\/abbr><\/td>//g;
    return $line;
}

sub _setOutputData
{
    my ($self, $lines) = @_;
    my @output;
    foreach (@$lines)
    {
	next unless ($self->_useInputLine($_));
	push(  @output, $self->_procesInputLine($_) );
    }
    $self->set_input_data(\@output);
}

sub _pullOnlineData
{
    my $self = shift;
    my $agent = WWW::Mechanize->new( cookie_jar => {} );
    $agent->get('https://service.aquacard.co.uk/aqua/web_channel/cards/security/logon/logon.aspx');
    $agent->form_id('mainform');
    $agent->set_fields('datasource_3a4651d1-b379-4f77-a6b1-5a1a4855a9fd' => $self->USER_NAME);
    $agent->set_fields('datasource_2e7ba395-e972-4a25-8d19-9364a6f06132' => $self->SURNAME);
    $agent->set_fields( '__EVENTTARGET' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
    $agent->set_fields( '__EVENTARGUMENT' =>'Target_5d196b33-e30f-442f-a074-1fe97c747474' );
    $agent->submit();

    $agent->form_id('mainform');
    $agent->set_fields('datasource_20056b2c-9455-4e42-aeba-9351afe0dbe1' => $self->SECRET_WORD );

    my $secretNumbers = $self->_getPasscodes($agent);
    $agent->set_fields('selectedvalue_dc28fef3-036f-4e48-98a0-586c9a4fbb3c' => $$secretNumbers[0]);
    $agent->set_fields('selectedvalue_f758fdb6-4b2c-4272-b785-cb3989b67901' => $$secretNumbers[1]);
    $agent->set_fields( '__EVENTTARGET' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
    $agent->set_fields( '__EVENTARGUMENT' =>'Target_53ab78d3-78ed-46f1-a777-1fd7957e1165' );
    $agent->submit();

    my @lines = split ("\n",$agent->content());
    $self->_setOutputData(\@lines);
    return 1;

}

1;

