#!/usr/bin/perl

use strict;
use warnings;

use Selenium::Chrome;
use HTTP::Request::Common;
use CACertOrg::CA;
$ENV{PERL_LWP_SSL_VERIFY_HOSTNAME} = 0;
use JSON;
use LWP::UserAgent;

use Switch;

sub sendExpense 
{
    my $line = shift;
    my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0, SSL_verify_mode => 0x00, SSL_ca_file      => CACertOrg::CA::SSL_ca_file() });
    my $json = JSON->new->allow_nonref;
    my $header = ['Content-Type' => 'application/json; charset=UTF-8'];
    #my $url = 'https://localhost:8000/expenses/';
    my $url = 'https://debian.home:8000/expenses/';
    my $encoded_data = $json->encode($line);
    my $request = HTTP::Request->new('POST', $url, $header, $encoded_data);
    my $response = $ua->request($request);
    print ("Saving: $encoded_data\n");
    print ("Response: ", $response->code,"\n");
    #print ($response->message,"\n");
    return $response->code;
}

sub sendAssetSeries 
{
    my $line = shift;
    my $ua = LWP::UserAgent->new(ssl_opts => { verify_hostname => 0, SSL_verify_mode => 0x00, SSL_ca_file      => CACertOrg::CA::SSL_ca_file() });
    my $json = JSON->new->allow_nonref;
    my $header = ['Content-Type' => 'application/json; charset=UTF-8'];
    #  my $url = 'https://localhost:8000/assets/series';
    my $url = 'https://debian.home:8000/assets/series';
    my $encoded_data = $json->encode($line);
    my $request = HTTP::Request->new('POST', $url, $header, $encoded_data);
    my $response = $ua->request($request);
    print ("Saving: $encoded_data\n");
    print ("Response: ", $response->code,"\n");
    #print ($response->message,"\n");
    return $response->code;
}

sub newExpense
{
    my ($accountId, $currency) = @_;
    my %expense;
    $expense{'id'} = 0;
    $expense{'transactionReference'} = '';
    $expense{'description'} = '';
    $expense{'detailedDescription'} = '';
    $expense{'accountId'} = $accountId;
    $expense{'date'} = '';
    $expense{'processDate'} = '';
    $expense{'currency'} = $currency;
    $expense{'commission'} = '0';
    $expense{'metadata'} = {};
    $expense{'metadata'}{'temporary'} = $JSON::false;
    $expense{'fx'} = {};
    $expense{'fx'}{'amount'} = 0;
    $expense{'fx'}{'currency'} = '';
    $expense{'fx'}{'rate'} = 0;
    return \%expense;
}

sub newAssetSeries
{
    my ($assetId) = @_;
    my %as;
    $as{'id'} = 0;
    $as{'amount'} = '';
    $as{'assetId'} = $assetId;
    $as{'date'} = '';
    return \%as;
}

sub _processInPageExpense
{
    my ($driver, $account, $ccy) = @_;
    my $line = newExpense($account, $ccy);

    my $message = 0;

    for (my $i = 1; ;$i++)
    {
        my $table = '//*[@id="ctl00_ExternalContent_IntroArea_WPManager_DbgGWP1_grd1_Table1"]/tbody';
        my $type  = $table . "/tr[$i]/td[1]/span";
        my $value = $table . "/tr[$i]/td[2]/div/table/tbody/tr/td/span";

        last unless ( $driver->find_element_by_xpath($value) );
        my $actualType = $driver->find_element($type)->get_text();
        my $actualValue = $driver->find_element($value)->get_text();

        if ($actualType ne '' )
        {
            $message = 0;
            $message = 1 if ($actualType eq 'Message:');

            _addValueToLine($line, $actualType, $actualValue);
        }
        elsif ($message)
        {
            $actualType = 'Message:';
            $actualValue = $driver->find_element($value)->get_text();
            _addValueToLine($line, $actualType, $actualValue);
        }
    }

    return $line;

}

sub _processRubrikExpense
{
    my ($driver, $account, $ccy) = @_;
    my $line = newExpense($account, $ccy);

    my $message = 0;

    for (my $i = 1; ;$i++)
    {
        my $table = '//*[@id="rubrikInner"]/tbody';
        my $type  = $table . "/tr[$i]/td[1]/span";
        my $value = $table . "/tr[$i]/td[2]/div/table/tbody/tr/td/span";

        last unless ( $driver->find_element_by_xpath($value) );
        my $actualType = $driver->find_element($type)->get_text();
        my $actualValue = $driver->find_element($value)->get_text();

        if ($actualType ne '' )
        {
            $message = 0;
            $message = 1 if ($actualType eq 'Message:');

            _addValueToLine($line, $actualType, $actualValue);
        }
        elsif ($message)
        {
            $actualType = 'Message:';
            $actualValue = $driver->find_element($value)->get_text();
            _addValueToLine($line, $actualType, $actualValue);
        }
    }

    return $line;
}

sub _addReference
{
    my ($line, $value, $force) = @_;
    $force = 0 unless (defined $force);
    return if (($line->{'transactionReference'} ne '') && (not ($force)));
    $line->{'transactionReference'} = $value;
}

sub _addMessage
{
    my ($line, $value) = @_;
    if ($value =~ m/PaymentID: (.*)/)
    {
        _addReference($line, $1, 1);
    } else {
        $line->{'detailedDescription'} .= $value . "\n";
    }
}

sub _addValueToLine
{
    my ($line, $key, $value) = @_;
    switch ($key)
    {
        case 'Reference number:' { _addReference($line, $value) }
        case 'Reference:' { _addReference($line, $value) }
        case 'Text:' { $line->{'description'} = _cleanText($value) }
        case 'Amount:' { $value =~ s/,//g; $value =~ s/"//g; $line->{'amount'} = $value  }
        case 'Amount in DKK:' { $value =~ s/,//g; $value =~ s/"//g; $line->{'amount'} = $value  }
        case 'Date:' { $line->{'date'} = _formatDate($value) }
        case 'Value date:' { $line->{'processDate'} = _formatDate($value) }
        case 'Currency traded:' {$line->{'fx'}{'currency'} = $value }
        case 'Exchange rate mark-up:' {$value =~ s/,//g; $value =~ s/"//g; $line->{'commission'} = $value  }
        case 'Exchange rate:' { $value =~ s/ //g; $line->{'fx'}{'rate'} = $value + 0 }
        case 'Amount in foreign currency:' { $value =~ s/ //g; $value =~ s/"//g; $line->{'fx'}{'amount'} = $value + 0}
        case 'Status:' { $line->{'detailedDescription'} .= 'Status: ' . $value . "\n"}
        case 'Message:' { _addMessage($line, $value) }
        case 'Creditor message:' { _addMessage($line, $value) }

        # Fields for a transfer
        case 'Text on account statement:' { $line->{'description'} = _cleanText($value) }
        case 'Bank:' { $line->{'detailedDescription'} .= 'Bank: ' . $value . "\n"}
        case 'Type of payment:' { $line->{'detailedDescription'} .= 'Type of payment: ' . $value . "\n"}
        case 'From account:' { $line->{'detailedDescription'} .= 'From account: ' . $value . "\n"}
        case 'Payment reference:' { $line->{'detailedDescription'} .= 'Payment reference: ' . $value . "\n"}
        case 'To account:' { $line->{'detailedDescription'} .= 'To account: ' . $value . "\n"}

        # Direct Debit - assuming date already set
        case 'Agreement no.:' { $line->{'transactionReference'} = $value . '-' . $line->{'date'}  }
    }
}

sub _processExpenseLine
{
    my ($driver) = @_;
    my $line = newExpense(6, 'DKK');

    my $frame = '//*[@id="indhold"]';
    $driver->switch_to_frame($driver->find_element($frame));

    my $message = 0;

    for (my $i = 1; ;$i++)
    {
        my $table = '//*[@id="parent.top.indhold.R1overflow"]/table/tbody';
        my $row   = $table . "/tr[$i]";
        my $type  = $row . '/td[1]';
        my $value = $row . '/td[2]';

        last unless ( $driver->find_element_by_xpath($row) );
        next unless ( $driver->find_element_by_xpath($value) );
        my $actualType = $driver->find_element($type)->get_text();
        my $actualValue = $driver->find_element($value)->get_text();

        if ($actualType ne '' )
        {
            $message = 0;
            $message = 1 if ($actualType eq 'Message:');
            $message = 1 if ($actualType eq 'Creditor message:');

            _addValueToLine($line, $actualType, $actualValue);
        }
        elsif ($message)
        {
            $actualType = 'Message:';
            $actualValue = $driver->find_element($value)->get_text();
            _addValueToLine($line, $actualType, $actualValue);
        }
    }

    $line->{'detailedDescription'} =~ s/\n*$//;

    return $line;
}

sub _formatDate
{
    my ($date) = @_;
    $date =~ m/([0-9]{2}).([0-9]{2}).([0-9]{4})/;
    return "$3-$2-$1";
}

sub _cleanText
{
    my ($text) = @_;
    $text =~ s/&amp;/&/g;
    $text =~ s/ \)\)\)\)//g;
    return $text;
}

sub pullOnlineData
{
    my ($account, $ccy, $assetId) = @_;
    my $driver = Selenium::Chrome->new( custom_args => '--proxy-auto-detect');
    #$driver->debug_on;

    #$driver->get('https://www.danskebank.dk/en-dk/Personal/Pages/personal.aspx?secsystem=J2');
    $driver->get('https://danskebank.dk/en/personal/help?n-login=pbnetbank');
	sleep 90;

    my @loadedExpenses;
    my $justLoaded = 1;
    my $lastDate = "";

    # starting at 2, as first row is header
    for (my $i = 2; ;$i++)
    {
        sleep 5 if ($justLoaded);
        $justLoaded = 0;
        # define xpaths
        my $expenseTable = '//*[@id="db-tl-table"]';
        my $expenseRow = $expenseTable .    "/tbody/tr[$i]";
        my $reconciledBox = $expenseRow .   "/td[12]/div/input";
        my $categorisation = $expenseRow .  "/td[2]/div/div/a";
        # column number varies if their classifications are shown
        my $expenseDetails = $expenseRow .  "/td[4]/div/a";
        my $balance = $expenseRow .  "/td[10]";
        my $date = $expenseRow .  "/td[1]";

        my @dataElements;
        $dataElements[0] = '//*[@id="ctl00_ExternalContent_IntroArea_WPManager_DbgGWP1_grd1_Table1"]/tbody';
        $dataElements[1] = '/html/body/div/form/table[2]';
        $dataElements[2] = '//*[@id="rubrikInner"]';

        # Wait until page fully loaded 
#        _waitForElement($agent, $expenseTable);

        # if should process
        last unless ($driver->find_element_by_xpath($expenseRow));
        next unless ($driver->find_element_by_xpath($expenseDetails));

        # deal with values
        my $dateVal = $driver->find_element_by_xpath($date)->get_text();
        unless ($dateVal eq $lastDate) {
            $lastDate = $dateVal;
            my $balVal = $driver->find_element_by_xpath($balance)->get_text();;
            $balVal =~ s/,//g;
            $balVal =~ s/ //g;
            unless ($balVal eq "" ) {
                print ("Sending balance of $balVal on $dateVal\n");
                my $as = newAssetSeries($assetId);
                $as->{'date'} = _formatDate($dateVal);
                $as->{'amount'} = $balVal;
                sendAssetSeries($as);
            }
        }

        my $processedEx = 0;
        if ( $driver->find_element_by_xpath($reconciledBox) )
        {
            next if ( $driver->find_element_by_xpath($reconciledBox)->is_selected() );
            $processedEx = 1;
        }

        # follow link
        my $element = $driver->find_element($expenseDetails);
        $driver->mouse_move_to_location(element => $element);
        $element->click();
        $justLoaded = 1;

        # process
        sleep 5;
        my $line;

		if ($driver->find_element_by_xpath($dataElements[0]))
        {
            $line = _processInPageExpense($driver, $account, $ccy);
            unless ( $processedEx )
            {
                #$line->setAmount( $line->getAmount() * -1 );
                $line->{'metadata'}{'temporary'} = $JSON::true;
            }
        }
        #elsif ($driver->find_element_by_xpath($dataElements[2]))
        #{
        #    $line = _processRubrikExpense($driver, $account, $ccy);
        #}
        else
        {
            $line = _processExpenseLine($driver);
        }
        $driver->go_back();

        sleep 3;
            if ((sendExpense($line) == '200' ) && ($processedEx))
            {
                $driver->find_element($reconciledBox)->click();
                sleep 3;
            }
    
    }
}

sub _setAllValues
{
    my ($agent, $xpath, $value) = @_;
    foreach ($agent->xpath($xpath))
    {
        $_->{value} = $value;
    }
}

sub _waitForElementSelenium
{
	my ($driver, $element) = @_;
	#print "Waiting for: >>$element<<\n";
	# 300s max wait time before failing
	for (my $i = 0; $i < 300; $i++)
	{
		my $return = 0;
		try	  { $return = ($driver->find_element($element)); }
		catch { print "Got error: $_. Ignoring...\n"; };
		return 1 if ($return);
		sleep 1;
	}
	return 0;
}

pullOnlineData(6, 'DKK', 11);

