#!/usr/bin/perl

use strict;
use warnings;

package AutomaticClassifier;
use Moose;

use DataTypes::Expense;

#has 'categories' => (required => 1);
has 'numbers' => ();

sub BUILD
{
	my ($self) = @_;
}

sub classify
{
	my ($self, $expense) = @_;
	my $exactMatch = $self->_StrategyExactMatch($expense);
	if ($exactMatch) return keys %$exactMatch;
}

sub _getValidClassifications
{
	my ($self, $expense) = @_;
	return $self->numbers->getValidClassifications($expense);
}

sub _StrategyExactMatch
{
	my ($self, $expense) = @_;
	my %classifications;
	my $total = 0;
	my $biggest = 'NO_VALUE';
	my $biggestValue = 0;
	$classifications{$_}++ for (@{$self->getValidClassifications($expense)});
	foreach my $row ($self->numbers->getExactMatches($expense))
	{
		if (exists $classifications{$$row[0]})
		{
			if ($$row[1] > $biggestValue)
			{
				$biggest = $$row[0];
				$biggestValue = $$row[1];
			}
			
			$total += $$row[1];
		}
	}
	return if ($biggest='NO_VALUE');
	my %return;
	$return{$$biggest} = $biggestValue / $total;
	return \%return;
}
