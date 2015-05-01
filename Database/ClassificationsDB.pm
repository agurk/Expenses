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
    $sth->execute($classificationID);

	my $row = $sth->fetchrow_arrayref();

	my $classification = Classification->new();
	$classification->setClassificationID($classificationID);
	$classification->setDescription($$row[0]);
	$classification->setValidFrom($$row[1]);
	$classification->setValidTo($$row[2]);
	$classification->setExpense($$row[3]);

	return $classification;
}

sub saveClassification
{
    my ($self, $classification) = @_; 
    my $dbh = $self->_openDB();
	if ($classification->getClassificationID() > 0)
	{
		my $sth = $dbh->prepare('update classificationdef set name = ?,validFrom=?, validTo=?, isExpense=? from where cid = ?');
		$sth->execute($classification->getDescription(),
					  $classification->getValidFrom(),
					  $classification->getValidTo(),
					  $classification->isExpense(),
					  $classification->getClassificationID());
	}
	else
	{
		my $sth = $dbh->prepare('insert into classificationdef (name, validFrom, validTo, isExpense) values (?, ?, ?, ?)');
		$sth->execute($classification->getDescription(),
					  $classification->getValidFrom(),
					  $classification->getValidTo(),
					  $classification->isExpense(),);
	}
}

sub deleteClassification
{
    my ($self, $classificationID) = @_; 
    my $dbh = $self->_openDB();
    my $sth = $dbh->prepare('delete from classificationdef where cid = ?');
    $sth->execute($classificationID);
}
