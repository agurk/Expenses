#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config
from Expense import Expense
from FXValues import FXValues

class OverallExpenses:

    def __init__(self):
        self.fxValues = FXValues()
    #    self.date = date
        #date=time.strftime("%Y-%m-%d"

    def monthQuery(self, date):
        return 'select classificationdef.name, amount, ccy from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and classificationdef.isexpense;'.format(date)

    def yearQuery(self, date):
        return 'select classificationdef.name, amount, ccy from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) < date(\'{0}\',\'start of year\',\'+1 year\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and classificationdef.isexpense;'.format(date)

    def OverallExpenses(self, date, period, baseCCY):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        if period == 'year':
            query = self.yearQuery(date)
        else:
            query = self.monthQuery(date)
        cursor = conn.execute(query)
        allExes = {};
        for row in cursor:
            key = row[0]
            amount = self.fxValues.FXAmount(row[1], row[2], baseCCY, date)
            if key in allExes.keys():
                allExes[key] += amount
            else:
                allExes[key] = amount
        conn.close()
        return allExes

    # Deprecated
    def TotalAmount(self, exes):
        totalAmount = 0
        for key in exes.keys():
            totalAmount+=exes[key]
        return totalAmount

