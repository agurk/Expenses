#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config
from Expense import Expense

class OverallExpenses:

  # def __init__(self, date):
    #    self.date = date
        #date=time.strftime("%Y-%m-%d"

    def OverallExpenses(self, date, baseCCY='GBP'):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select classificationdef.name, amount, ccy from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and classificationdef.isexpense;'.format(date)
        cursor = conn.execute(query)
        allExes = {};
        for row in cursor:
            key = row[0]
            amount = self._getAmount(row[1], row[2], baseCCY, date)
            if key in allExes.keys():
                allExes[key] += amount
            else:
                allExes[key] = amount
        return allExes

    def _getAmount(self, amount, CCY, baseCCY, date):
        if CCY == baseCCY:
            return amount
        fxRate=self._getFXRate(CCY, baseCCY, date);
        return amount * fxRate

    def _getFXRate(self, CCY, baseCCY, date):
        return 0.095

    def TotalAmount(self, exes):
        totalAmount = 0
        for key in exes.keys():
            totalAmount+=exes[key]
        return totalAmount

