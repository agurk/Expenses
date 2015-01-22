#!/usr/bin/perl

use strict;
use warnings;

package AutomaticClassifier;
use Moose;

use DataTypes::Expense;

#has 'categories' => (required => 1);
has 'numbers' => (is => 'ro', required => 1);

sub BUILD
{
	my ($self) = @_;
}

sub classify
{
	my ($self, $expense) = @_;
	my $result = $self->_tryStrategies($expense);
}

sub _tryStrategies
{
	my ($self, $expense) = @_;
	my @result;
	my $match = $self->_StrategyExactMatch($expense);
	if (defined $match)
	{
		$expense->setExpenseClassification($$match[0]);
		print "selected $$match[0] by exact match, with likelihood $$match[1] for expense ",$expense->getExpenseDescription(),"\n";
		return;
	}	
	$match = $self->_StrategyStatisticalMatch($expense);
	{
		$expense->setExpenseClassification($$match[0]);
		print "selected $$match[0] by statistical match, with likelihood $$match[1] for expense $expense->getExpenseDescription()\n";
		return;
	}	
}

sub _getValidClassifications
{
	my ($self, $expense) = @_;
	return $self->numbers()->getValidClassifications($expense);
}

sub _StrategyStatisticalMatch
{
	my ($self, $expense) = @_;
	my %classifications;
	my $total = 0;
	my $biggest = 'NO_VALUE';
	my $biggestValue = 0;
	$classifications{$_}++ for (@{$self->_getValidClassifications($expense)});
	foreach my $row ($self->numbers()->getClassificationStats($expense))
	{
		if (defined $classifications{$$row[0]})
		{
			if ($$row[1] > $biggestValue)
			{
				$biggest = $$row[0];
				$biggestValue = $$row[1];
			}

			$total += $$row[1];
		}
	}
	return if ($biggest eq 'NO_VALUE');
	my @return;
	$return[0] = $biggest;
	$return[1] = $biggestValue / $total;
	return \@return;
}

sub _StrategyExactMatch
{
	my ($self, $expense) = @_;
	my %classifications;
	my $total = 0;
	my $biggest = 'NO_VALUE';
	my $biggestValue = 0;
	$classifications{$_}++ for (@{$self->_getValidClassifications($expense)});
	my $results = $self->numbers()->getExactMatches($expense);
	foreach my $row (@$results)
	{
		if (defined $classifications{$$row[0]})
		{
			if ($$row[1] > $biggestValue)
			{
				$biggest = $$row[0];
				$biggestValue = $$row[1];
			}
			
			$total += $$row[1];
		}
	}
	return if ($biggest eq 'NO_VALUE');
	my @returnable;
	$returnable[0] = $biggest;
	print $biggestValue,' ',$total,"\n";
	$returnable[1] = $biggestValue / $total;
	return \@returnable;
}

1;

