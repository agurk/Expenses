#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config
import expensesSQL

class FXValues:

    def __init__(self):
        self.foo='true'
        #conn = sqlite3.connect(config.SQLITE_DB)
#=        conn.text_factory = str 
##        cursor = conn.execute(expensesSQL.getCCYFormats())
#        for row in cursor:
#            self.ccyFormats[row[0]] = row[1]

    def FXAmount(self, amount, baseCCY, ccy, date):
        baseCCY = baseCCY.upper()
        ccy = ccy.upper()
        if baseCCY == ccy:
            return amount
        if ccy=='EUR':
            return self._toEur(amount, baseCCY)
        if ccy=='GBP':
            return self._toGbp(amount, baseCCY)
        if ccy == 'DKK':
            return self._toDkk(amount, baseCCY)

    def _toEur(self, amount, baseCCY):
        if baseCCY == 'GBP':
            return amount / 0.72
        if baseCCY == 'DKK':
            return amount / 7.46

    def _toGbp(self, amount, baseCCY):
        if baseCCY == 'EUR':
            return amount / 1.39
        if baseCCY == 'DKK':
            return amount / 10.36

    def _toDkk(self, amount, baseCCY):
        if baseCCY == 'GBP':
            return amount / 0.097
        if baseCCY == 'EUR':
            return amount / 0.13

