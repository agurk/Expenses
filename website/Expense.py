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
        expense = {}
        for row in cursor:
            expense['date'] = row[0]
            expense['description'] = row[1]
            expense['amount'] = row[2]
            expense['name'] = row[3]
            expense['eid'] = row[4]
            expense['confirmed'] = row[5]
            expense['tag'] = row[6]
            expense['fxamount'] = row[7]
            expense['fxccy'] = row[8]
            expense['fxrate'] = row[9]
            expense['fxcommission'] = row [10]
            return expense 
