#!/usr/bin/python3

import config
import sqlite3
import re

from OverallExpenses import OverallExpenses
from FXValues import FXValues
import time
from calendar import monthrange
from datetime import datetime
import calendar

class Analysis:

    def __init__(self, ccy='GBP'):
        self.oe = OverallExpenses()
        self.fxValues = FXValues()
        self.ccy = ccy
        self.startYear = 2011
        self.endYear = 2019

    def DaysInMonth(self, date):
        graphDate = datetime.strptime(self.date, '%Y-%m-%d')
        days = monthrange(graphDate.year, graphDate.month)
        return days[1]
    
    def YearlySpend(self):
        salary={}
        reinem={}
        expenses={}
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        for year in range (self.startYear, self.endYear):
            salary[year] = 0
            reinem[year] = 0
            expenses[year] = 0
            date = '{0}-01-01'.format(year)
            query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and (c.cid = 17) and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) <= date(\'{0}\',\'start of year\', \'+12 months\')'.format(date)
            cursor = conn.execute(query)
            for row in cursor:
                salary[year] += self.fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
            query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and (c.cid=18) and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) <= date(\'{0}\',\'start of year\', \'+12 months\')'.format(date)
            cursor = conn.execute(query)
            for row in cursor:
                reinem[year] += self.fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
            query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and c.cid = 12 and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) <= date(\'{0}\',\'start of year\', \'+12 months\')'.format(date)
            cursor = conn.execute(query)
            for row in cursor:
                expenses[year] += self.fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
        #for key in salary.keys():
        #    print(key, ';', salary[key],';',reinem[key],';',expenses[key])
        conn.close()
        output = {}
        output['salary'] = salary
        output['reimbursements'] = reinem
        output['expenses'] = expenses
        return output 
    
    
    def CumulativeSpend(dateIn):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) <= date(\'{0}\',\'start of year\', \'+12 months\')'.format(dateIn)
        cursor = conn.execute(query)
        # TODO leap years
        amounts = [0] * (366 + 1)
        for row in cursor:
            dayNo = datetime.strptime(row[2], "%Y-%m-%d").timetuple().tm_yday
            amounts[dayNo] += fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
        for i in range(1, 366):
            amounts[i] = amounts[i] + amounts[i-1]
        conn.close()
        return amounts
    
    
    def yearTotals(self):
        totals = {}
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense'.format(self.startYear, self.endYear)
        cursor = conn.execute(query)
        for row in cursor:
            year = datetime.strptime(row[2], '%Y-%m-%d').year
            if year in totals.keys():
                totals[year] += self.fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
            else:
                totals[year] = self.fxValues.FXAmount(row[0],row[1],self.ccy,row[2])
        return totals

        
    
    def yearTotals2(self):
        totals = {}
        for year in range (self.startYear, self.endYear):
            for month in range(1, 13):
                if month < 10:
                    month = '0' + str(month)
                date = '{0}-{1}-01'.format(year, month)
                overall = self.oe.OverallExpenses(date)
                # kludgy hack  as we have all the actual category data here
                for key in overall.keys():
                    if key in totals.keys():
                        totals[year] += overall[key]
                    else:
                        totals[year] = overall[key]
    #        print(year)
    #        for key in totals.keys():
    #            print(key,';',totals[key])
    #        print()
        return totals
        
    
    def monthTotals(self):
        for year in range (self.startYear, self.endYear):
            for month in range(1, 13):
                if month < 10:
                    month = '0' + str(month)
                date = '{0}-{1}-01'.format(year, month)
                print(date,';', end='')
                overall = oe.OverallExpenses(date)
                print(oe.TotalAmount(overall))
