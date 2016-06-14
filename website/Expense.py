#!/usr/bin/python

import sqlite3
import time
import datetime
import re
from datetime import date, timedelta
import config
import expensesSQL
from FXValues import FXValues

class Expense:

    ccyFormats={}
    fxValues = FXValues()

    def __init__(self):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        cursor = conn.execute(expensesSQL.getCCYFormats())
        for row in cursor:
            self.ccyFormats[row[0]] = row[1]
        conn.close()

    def Expense(self, eid, ccy=''):
        expense = []
        if eid:
            conn = sqlite3.connect(config.SQLITE_DB, uri=True)
            cursor = conn.execute(expensesSQL.getExpense(eid))
            for row in cursor:
                expense = self._makeExpense(row, ccy, conn)
            conn.close()
        return expense

    def NewExpense(self, did='', ccy=''):
        empty = [''] * 13
        #date
        empty[0]=time.strftime("%Y-%m-%d")
        #amount
        empty[2] = '0'
        expense = self._makeExpense(empty, ccy, '')
        if did:
            conn = sqlite3.connect(config.SQLITE_DB, uri=True)
            self._addSingleDocument(expense, did, conn)
        return expense

    def Expenses(self, date, allExes, ccy=''):
        if allExes == 'true':
           return self._Expenses(date, 'ALL', ccy)
        else:
            return self._Expenses(date, '', ccy)

    def _Expenses(self, date, condition, ccy):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        if condition == 'ALL':
            sql = expensesSQL.getAllOneMonthsExpenses(date)
        else:
            sql = expensesSQL.getSomeOneMonthsExpenses(date)
        expenses=[]
        cursor = conn.execute(sql)
        cursor = conn.execute(sql)
        for row in cursor:
            expenses.append(self._makeExpense(row, ccy, conn))
        conn.close()
        return expenses  

    def Search (self, search, ccy=''):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        cursor = conn.execute(expensesSQL.getSimilarExpenses(search))
        expenses=[]
        for row in cursor:
            expenses.append(self._makeExpense(row, ccy, conn))
        conn.close()
        return expenses  

    def _makeExpense(self, row, ccy, conn):
        expense = {}
        if ccy == '' or ccy == 'base':
            expense['amount'] = row[2]
            expense['ccy'] = row[3]
            expense['fxcommission'] = row[11]
        elif ccy == 'original':
            self._originalCCY(expense, row, ccy)
        else:
            expense['amount'] = self.fxValues.FXAmount(row[2], row[3], ccy, row[0])
            expense['ccy'] = ccy
            expense['fxcommission'] = row[11]
        expense['date'] = row[0]
        expense['description'] = row[1]
        self._fixAmount(expense)
        expense['pretty_amount'] = self._makePrettyAmount(expense['amount'], expense['ccy'])
        expense['classification'] = row[4]
        expense['eid'] = row[5]
        expense['confirmed'] = row[6]
        expense['tag'] = row[7]
        expense['fxamount'] = row[8]
        expense['fxccy'] = row[9]
        expense['fxrate'] = row[10]
        expense['aid'] = row[12]
        if row[13] == '' or row[13] == 0 or row[13] == None:
            expense['temporary'] = False
        else:
            expense['temporary'] = True 
        expense['reference'] = row[14]
        expense['modified'] = row[15]
        expense['detaileddescription'] = row[16]
        self._addRawIDs(expense, conn)
        self._addDocuments(expense, conn)
        return expense

    def _fixAmount(self, expense):
        amnt = expense['amount']
        if isinstance( amnt, float):
            return
        amnt = re.sub(r'[,.]([0-9]{3}[.,])',r'\1', amnt)
        expense['amount'] = float(amnt.replace(',','.'))

    def _originalCCY(self, expense, row, ccy):
        if row[9] is None or row[9] == '':
            expense['ccy'] = row[3]
            expense['amount'] = row[2]
        else:
            expense['ccy'] = row[9]
            expense['amount'] = row[8]
            expense['fxcommission'] = row[11]

    def _makePrettyAmount(self, amount, ccy):
        amount = float(amount)
        roundedAmount = '%.2f' % amount
        if ccy in self.ccyFormats.keys():
            amount = self.ccyFormats[ccy].format(roundedAmount)
        else:
            amount = str(ccy) + ' ' + roundedAmount
        return amount

    def _addRawIDs(self, expense, db):
        # TODO clean up not-UTF8 in the db
        results = []
        if db:
            try:
                for row in db.execute(expensesSQL.getRawLines(expense['eid'])):
                    results.append(row)
            except:
                db.text_factory = lambda x: str(x, 'latin1')
                for row in db.execute(expensesSQL.getRawLines(expense['eid'])):
                    results.append(row)
                db.text_factory = str
        expense['rawlines'] = results 

    def _addDocuments(self, expense, db):
        if db:
            cursor = db.execute(expensesSQL.getDocuments(expense['eid']))
            documents=[]
            for row in cursor:
                documents.append({'did': row[0], 'filename': row[1]})
            expense['documents'] = documents

    def _addSingleDocument(self, expense, did, db):
        cursor = db.execute(expensesSQL.getDocument(did))
        documents=[]
        for row in cursor:
            documents.append({'did': did, 'filename': row[0]})
        expense['documents'] = documents
