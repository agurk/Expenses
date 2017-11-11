#!/usr/bin/perl

package Loader_AMEX;
use Moose;
extends 'Loader';

use WWW::Mechanize::Firefox;
use LWP::ConnCache;

use Selenium::Chrome;

has 'AMEX_PASSWORD' => ( is => 'rw', isa=>'Str', writer => 'setAmexPass');
has 'AMEX_USERNAME' => ( is => 'rw', isa=>'Str', writer => 'setAmexUser');
has 'AMEX_CARD_NUMBER' => ( is => 'rw', isa=>'Str', writer => 'setAmexCardNo');
# Index is 0-rated
has 'AMEX_INDEX' => ( is => 'rw', isa=>'Str', writer => 'setAmexIndex');

use constant csvFile => '/home/timothy/Downloads/ofx.csv';

# build string formats:
# file; filename
# notfile; cardno; user; password; index
sub BUILD
{
    my ($self) = @_;
    my @buildParts = split (';' ,$self->build_string);
    # if it is a file
    if ($buildParts[0])
    {
        $self->setFileName($buildParts[1]);
    }
    else
    {
        $self->setAmexCardNo($buildParts[1]);
        $self->setAmexUser($buildParts[2]);
        $self->setAmexPass($buildParts[3]);
        $self->setAmexIndex($buildParts[4]);
    }
}

# The AMEX form, once that page has been reached is quite simple, and three input fields need to be set:
# From the DownloadForm:
# Format => download format, we're using 'CSV'
# selectradio => with the value of the card number
# selectradio => with the value set to the statement periods we want to download
sub _pullOnlineData
{
    my $self = shift;
    unlink csvFile;

	my $driver = Selenium::Chrome->new( custom_args => '--headless --no-proxy-server');

	$driver->get('https://www.americanexpress.com/');
	#$driver->find_element_by_id('//*[@id="sprite-ContinueButton_EN"]')->send_keys("\n");
	
	$driver->find_element('//*[@id="Username"]|//*[@id="UserID"]')->send_keys($self->AMEX_USERNAME);
	$driver->find_element('//*[@id="Password"]')->send_keys($self->AMEX_PASSWORD);
	$driver->find_element('//*[@id="loginLink"]|//*[@id="loginButton"]')->click();

    sleep 10;
	
	#$driver->find_element_by_id('//*[@id="sprite-ContinueButton_EN"]')->send_keys("\n");

	$driver->get('https://global.americanexpress.com/myca/intl/download/emea/download.do?request_type=&Face=en_GB&BPIndex=0&sorted_index=0&inav=gb_myca_pc_statement_export_statement_data');

	sleep 5;
	
	$driver->find_element('//*[@id="sprite-ContinueButton_EN"]')->send_keys("\n");

	$driver->find_element('//*[@id="CSV"]')->click();
	
	$driver->find_element('//*[@id="selectCard1'.$self->AMEX_INDEX.'"]')->click();
	
	$driver->find_element('//*[@id="radioid'.$self->AMEX_INDEX.'0"]')->click();
	$driver->find_element('//*[@id="radioid'.$self->AMEX_INDEX.'1"]')->click();
	$driver->find_element('//*[@id="radioid'.$self->AMEX_INDEX.'2"]')->click();
	$driver->find_element('//*[@id="radioid'.$self->AMEX_INDEX.'3"]')->click();
	
	$driver->find_element('//*[@id="myBlueButton1"]')->click();

    sleep 15;

	$driver->quit();

    # TODO: catch download properly
    $self->setFileName(csvFile);
    return $self->_loadCSVRows();

}

1;

