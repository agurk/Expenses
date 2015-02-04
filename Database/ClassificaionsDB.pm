#!/usr/bin/perl

package ClassificationsDB;
use Moose;
extends 'DAL';

use DataTypes::Classification;

sub getClassification
{
    my ($self, $classificationID) = @_; 
    my $dbh = $self->_openDB();
    my $sth = $dbh->prepare('select name, validFrom, validTo, isExpense from classificationdef where cid = ?');
    $sth->execute(@bindValues);

	my $row = $sth->fetchrow_arrayref

	my $classification = $Classification->new();
	$classification->setClassificationID($classificationID);
	$classification->setDescription($$row[0]);
	$classification->setValidFrom($$row[1]);
	$classification->setValidTo($$row[2]);
	$classification->isExpense($$row[3]);

    $dbh->disconnect();
	return $classification;
}

sub saveClassification
{
}
