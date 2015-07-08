#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config
import expensesSQL
from FXValues import FXValues

class Expense:

    ccyFormats={}
    fxValues = FXValues()

    def __init__(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        cursor = conn.execute(expensesSQL.getCCYFormats())
        for row in cursor:
            self.ccyFormats[row[0]] = row[1]

    def Expense(self, eid, ccy=''):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        cursor = conn.execute(expensesSQL.getExpense(eid))
        for row in cursor:
            return self._makeExpense(row, ccy, conn)

    def Expenses(self, date, allExes, ccy=''):
        if allExes == 'true':
           return self._Expenses(date, 'ALL', ccy)
        else:
            return self._Expenses(date, '', ccy)

    def _Expenses(self, date, condition, ccy):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        if condition == 'ALL':
            sql = expensesSQL.getAllOneMonthsExpenses(date)
        else:
            sql = expensesSQL.getSomeOneMonthsExpenses(date)
        cursor = conn.execute(sql)
        expenses=[]
        for row in cursor:
            expenses.append(self._makeExpense(row, ccy, conn))
        return expenses  

    def Search (self, search, ccy=''):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        cursor = conn.execute(expensesSQL.getSimilarExpenses(search))
        expenses=[]
        for row in cursor:
            expenses.append(self._makeExpense(row, ccy, conn))
        return expenses  

    def _makeExpense(self, row, ccy, conn):
        expense = {}
        if ccy == '' or ccy == 'base':
            expense['amount'] = row[2]
            expense['ccy'] = row[3]
            expense['fxcommission'] = row[11]
        else:
            expense['amount'] = self.fxValues.FXAmount(row[2], row[3], ccy, row[0])
            expense['ccy'] = ccy
            expense['fxcommission'] = row[11]
        expense['date'] = row[0]
        expense['description'] = row[1].decode('utf8', 'ignore')
        expense['pretty_amount'] = self._makePrettyAmount(expense['amount'], expense['ccy'])
        expense['classification'] = row[4]
        expense['eid'] = row[5]
        expense['confirmed'] = row[6]
        expense['tag'] = row[7]
        expense['fxamount'] = row[8]
        expense['fxccy'] = row[9]
        expense['fxrate'] = row[10]
        self._addRawIDs(expense, conn)
        self._addDocuments(expense, conn)
        return expense

    def _makePrettyAmount(self, amount, ccy):
        roundedAmount = '%.2f' % amount
        amount = self.ccyFormats[ccy].format(roundedAmount)
        return amount.decode('utf-8')

    def _addRawIDs(self, expense, db):
        cursor = db.execute(expensesSQL.getRawLines(expense['eid']))
        expense['rawlines'] = cursor

    def _addDocuments(self, expense, db):
        cursor = db.execute(expensesSQL.getDocuments(expense['eid']))
        expense['documents'] = cursor

