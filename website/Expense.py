#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config
import expensesSQL

class Expense:

    def Expense(self, eid):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        cursor = conn.execute(expensesSQL.getExpense(eid))
        for row in cursor:
            return self.makeExpense(row, conn)

    def Expenses(self, date, condition=''):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        if condition == 'ALL':
            sql = expensesSQL.getAllOneMonthsExpenses(date)
        else:
            sql = expensesSQL.getSomeOneMonthsExpenses(date)
        cursor = conn.execute(sql)
        expenses=[]
        for row in cursor:
            expenses.append(self.makeExpense(row, conn))
        return expenses  

    def makeExpense(self, row, conn):
        expense = {}
        expense['date'] = row[0]
        expense['description'] = row[1].decode('utf8', 'ignore')
        expense['amount'] = row[2]
        expense['name'] = row[3]
        expense['eid'] = row[4]
        expense['confirmed'] = row[5]
        expense['tag'] = row[6]
        expense['fxamount'] = row[7]
        expense['fxccy'] = row[8]
        expense['fxrate'] = row[9]
        expense['fxcommission'] = row [10]
        self._addRawIDs(expense, conn)
        self._addDocuments(expense, conn)
        return expense
    

    def _addRawIDs(self, expense, db):
        cursor = db.execute(expensesSQL.getRawLines(expense['eid']))
        expense['rawlines'] = cursor

    def _addDocuments(self, expense, db):
        cursor = db.execute(expensesSQL.getDocuments(expense['eid']))
        expense['documents'] = cursor

