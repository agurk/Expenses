#!/usr/bin/env perl 
#===============================================================================
#
#         FILE: ExpenseDB.pm
#
#  DESCRIPTION: Data Access Layer for Expense Object
#
#      OPTIONS: ---
# REQUIREMENTS: ---
#         BUGS: ---
#        NOTES: ---
#       AUTHOR: Timothy Moll
# ORGANIZATION: 
#      VERSION: 0.2
#      CREATED: 23/12/14 11:19:12
#     REVISION: ---
#===============================================================================

package ExpenseDB;
use Moose;
extends 'DAL';

use constant RAW_TABLE=>'RawData';
use constant EXPENSES_TABLE=>'Expenses';
use constant CLASSIFIED_DATA_TABLE=>'Classifications';
use constant CLASSIFICATION_DEFINITION_TABLE=>'ClassificationDef';
use constant ACCOUNT_DEFINITION_TABLE=>'AccountDef';
use constant LOADER_DEFINITION_TABLE=>'LoaderDef';
use constant PROCESSOR_DEFINITION_TABLE=>'ProcessorDef';
use constant ACCOUNT_LOADERS_TABLE=>'AccountLoaders';
use constant EXPENSE_RAW_MAPPING_TABLE => 'ExpenseRawMapping';

use strict;
use warnings;
use utf8;

use DBI;
use DataTypes::Expense;
use Time::Piece;

sub getExpense
{
	my ($self, $expenseID) = @_;
	my $dbh = $self->_openDB();
	my $query = 'select e.aid, e.description, e.amount, e.ccy, e.amountFX, e.ccyFX, e.fxRate, e.commission, e.date, e.modified, c.cid, c.confirmed from expenses e, classifications c where e.eid = ? and e.eid = c.eid';
	my $sth = $dbh->prepare($query);
	$sth->execute($expenseID);
	
	my $row = $sth->fetchrow_arrayref();
	for (my $i = 0; $i < 12; $i++)
	{
		$$row[$i] = '' unless (defined $$row[$i]);
	}

	my $expense = Expense->new(	ExpenseID=>$expenseID,
								AccountID=>$$row[0],
								Description=>$$row[1],
								Amount=>$$row[2],
								Currency=>$$row[3],
								FXAmount=>$$row[4],
								FXCCY=>$$row[5],
								FXRate=>$$row[6],
								Commission=>$$row[7],
								Date=>$$row[8],
								Modified=>$$row[9],
								Classification=>$$row[10],
								Confirmed=>$$row[11],
						   	  );
	
	$query = 'select rid from expenserawmapping where eid = ?';
	$sth = $dbh->prepare($query);
	$sth->execute($expenseID);

	foreach my $row ( $sth->fetchrow_arrayref())
	{
		$expense->addRawID($$row[0]);
	}

	$query = 'select tag from tagged where eid = ?';
	$sth = $dbh->prepare($query);
	$sth->execute($expenseID);

	foreach my $row ( $sth->fetchrow_arrayref())
	{
		if ($row)
		{
			$expense->setTagged($$row[0]) if ($$row[0]);
		}
	}

	return $expense;
}

sub _createNewExpense
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();

	my $insertString='insert into '.EXPENSES_TABLE.' (aid, description, amount, ccy, amountFX, ccyFX, fxRate, commission, date) values (?, ?, ?, ?, ?, ?, ?, ?, ?)';
	my $sth = $dbh->prepare($insertString);
	$sth->execute($self->_makeTextQuery($expense->getAccountID()),
			$self->_makeTextQuery($expense->getDescription()),
			$expense->getAmount(),
			$expense->getCCY(),
			$expense->getFXAmount(),
			$expense->getFXCCY(),
			$expense->getFXRate(),
			$expense->getCommission(),
			$expense->getDate());
	$sth->finish();

	# TODO: make this a bit safer
	$sth=$dbh->prepare('select max(eid) from expenses');
	$sth->execute();
	$expense->setExpenseID($sth->fetchrow_arrayref()->[0]);
	$sth->finish();

	$self->_setExpensesRawClassification($expense);
	$self->setTagged($expense->getExpenseID, $expense->isTagged());
}

sub setTagged
{
	my ($self, $eid, $tag) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('delete from tagged where eid = ?');
	$sth->execute($eid);
	$sth->finish();

	if ($tag)
	{
		$sth = $dbh->prepare('insert into tagged (eid, tag) values (?, ?)');
		$sth->execute($eid, $tag);
	}
	$sth->finish();
}

sub _setExpensesRawClassification
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('select distinct rid from expenserawmapping where eid = ?');
	$sth->execute($expense->getExpenseID);

	my %RIDS;
	while (my $row = $sth->fetchrow_arrayref())
	{
		$RIDS{$$row[0]} = 0 if (defined $row);
	}

	foreach (@{$expense->getRawIDs()})
	{
		if (exists $RIDS{$_})
		{
			delete $RIDS{$_};
			next;
		}
		my $insertString='insert into '. EXPENSE_RAW_MAPPING_TABLE .' (eid, rid) values (?, ?)';
		my $sth=$dbh->prepare($insertString);
		$sth->execute($expense->getExpenseID(), $self->_makeTextQuery($_));
		$sth->finish();
	}

	foreach (keys %RIDS)
	{
		my $query = 'delete from expenserawmapping where eid = ? and rid = ?';
		my $sth=$dbh->prepare($query);
		$sth->execute($expense->getExpenseID(), $self->_makeTextQuery($_));
		$sth->finish();
	}


	$sth=$dbh->prepare('delete from classifications where eid = ?');
	$sth->execute($expense->getExpenseID);
	
	my $confirmed = 0;
	$confirmed = 1 if ($expense->isConfirmed());

	my $insertString2='insert into '.CLASSIFIED_DATA_TABLE.' (eid, cid, confirmed) values (?, ?, ?)';
	$sth = $dbh->prepare($insertString2);
	$sth->execute($self->_makeTextQuery($expense->getExpenseID()), $self->_makeTextQuery($expense->getClassification()), $confirmed);

}

sub _updateExpense
{
	my ($self, $expense) = @_;
	my $dbh = $self->_openDB();
	my $query = 'update expenses set aid = ?, description = ?, amount = ?, ccy = ?, amountFX = ?, ccyFX = ?, fxRate = ?, commission = ?, date = ? where eid = ?';
	my $sth = $dbh->prepare($query);
	$sth->execute($self->_makeTextQuery($expense->getAccountID()),
			$self->_makeTextQuery($expense->getDescription()),
			$expense->getAmount(),
			$expense->getCCY(),
			$expense->getFXAmount(),
			$expense->getFXCCY(),
			$expense->getFXRate(),
			$expense->getCommission(),
			$expense->getDate(),
			$expense->getExpenseID);
	$sth->finish();

	$self->_setExpensesRawClassification($expense);
	$self->setTagged($expense->getExpenseID, $expense->isTagged());
}

sub saveExpense
{
	my ($self, $expense) = @_;
	if ($expense->getExpenseID > -1)
	{
		$self->_updateExpense($expense);
	}
	else
	{
		$self->_createNewExpense($expense);
	}
}

sub mergeExpenses
{
	my ($self, $primaryExpense, $secondaryExpense) = @_;
	my $dbh = $self->_openDB();
	$dbh->{AutoCommit} = 0;

	eval
	{
		my $sth=$dbh->prepare('select rid from expenserawmapping where eid = ?');
		$sth->execute($secondaryExpense);
		foreach my $row ( $sth->fetchrow_arrayref())
		{
			my $sth2 = $dbh->prepare('insert into expenserawmapping (eid, rid) values(?,?)');
			$sth2->execute($primaryExpense, $row->[0]);
		}
		$sth = $dbh->prepare('delete from expenses where eid = ?');
		$sth->execute($secondaryExpense);
		$sth = $dbh->prepare('delete from expenserawmapping where eid = ?');
		$sth->execute($secondaryExpense);
		$sth = $dbh->prepare('delete from classifications where eid = ?');
		$sth->execute($secondaryExpense);
		$sth = $dbh->prepare('delete from tagged where eid = ?');
		$sth->execute($secondaryExpense);

		$dbh->commit();

	};

    if($@)
	{
		warn "Error inserting the link and tag: $@\n";
		$dbh->rollback();
	}

}

sub confirmClassification
{
	my ($self, $expenseID) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('update classifications set confirmed = 1 where eid = ?');
	$sth->execute($expenseID);
	$sth->finish();
}

# Removes existing classifications so can be used also to update an existing one
sub saveClassification
{
	my ($self, $expenseID, $classificationID, $confirmed) = @_;
	my $dbh = $self->_openDB();
	$dbh->{AutoCommit} = 0;

	eval
	{
		my $sth = $dbh->prepare('delete from classifications where eid = ?');
		$sth->execute($expenseID);
		$sth->finish();
		$sth = $dbh->prepare('insert into classifications (eid, cid, confirmed) values (?, ?, ?)');
		$sth->execute($expenseID, $classificationID, $confirmed);
		$sth->finish();
		$dbh->commit();
		$dbh->disconnect();
	};
    
	if($@)
	{
		warn "Error saving classification $classificationID for expense $expenseID\n";
		$dbh->rollback();
	}
}

sub saveAmount
{
	my ($self, $expenseID, $amount) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare('update expenses set amount = ?, modified = ? where eid = ?');
	$sth->execute($amount, $self->_getCurrentDateTime() ,$expenseID);
	$sth->finish();
	$dbh->disconnect();
}



sub findExpense
{
	my ($self, $aid, $date, $description, $amount, $ccy) = @_;
	my $dbh = $self->_openDB();
	my $sth = $dbh->prepare("select eid from expenses where aid = ? and date = ? and description = ? and amount = ? and ccy = ?");
    $sth->execute($aid, $date, $description, $amount, $ccy);

	my $row = $sth->fetchrow_array;
	if ($row)
	{
		return $self->getExpense($row);
	}
	else
	{
		return;
	}

}

sub duplicateExpense
{
	my ($self, $eid) = @_;
	my $originalExpense = $self->getExpense($eid);
	my $newExpense = Expense->new ( AccountID=>$originalExpense->getAccountID,
									Description=>$originalExpense->getDescription,
									Amount=>$originalExpense->getAmount,
									Currency=>$originalExpense->getCCY,
									FXAmount=>$originalExpense->getFXAmount,
									FXCCY=>$originalExpense->getFXCCY,
									FXRate=>$originalExpense->getFXRate,
									Commission=>$originalExpense->getCommission,
									Date=>$originalExpense->getDate,
									Classification=>$originalExpense->getClassification,
						   	  );
	foreach (@{$originalExpense->getRawIDs})
	{
		$newExpense->addRawID($_);
	}

	$self->saveExpense($newExpense);
}

1;

