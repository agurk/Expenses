#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config

class MonthView:

    def __init__(self, date):
        self.date = date
        #date=time.strftime("%Y-%m-%d"

    def TotalAmount(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select printf("%.2f", sum(amount) * -1) from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and classificationdef.isexpense;'.format(self.date)
        cursor = conn.execute(query)
        for row in cursor:
            totalAmount = row[0]
        return totalAmount

    def OverallExpenses(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select count (*), classificationdef.name, printf("%.2f", sum(amount) * -1) as amt from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and classificationdef.isexpense group by classifications.cid order by amt+0 asc;'.format(self.date)
        cursor = conn.execute(query)
        return cursor

    def IndividualExpensesAll(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select date, description, printf("%.2f", amount), classificationdef.name, expenses.eid, classifications.confirmed, tag, d.did from expenses left join tagged on expenses.eid = tagged.eid left join documentexpensemapping d on expenses.eid = d.eid, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid order by date desc;'.format(self.date)
        cursor = conn.execute(query)
        return cursor

    def IndividualExpenses(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select date, description, printf("%.2f", amount), classificationdef.name, expenses.eid, classifications.confirmed, tag, d.did from expenses left join tagged on expenses.eid = tagged.eid left join documentexpensemapping d on expenses.eid = d.eid, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid and (classificationdef.isexpense or not classifications.confirmed) order by date desc;'.format(self.date)
        cursor = conn.execute(query)
        return cursor

    def add_months(self, sourcedate, months):
        month = sourcedate.tm_mon - 1 + months
        year = sourcedate.tm_year + month / 12
        month = month % 12 + 1
        day = 1
        return datetime.date(year,month,day)

    def get_date(self, sourcedate):
        month = sourcedate.tm_mon
        year = sourcedate.tm_year
        day = 1
        return datetime.date(year,month,day)

    def PreviousMonth(self):
        previous = time.strptime(self.date, "%Y-%m-%d")
        return self.add_months(previous, -1)

    def ThisMonth(self):
        thisM = time.strptime(self.date, "%Y-%m-%d")
        return self.get_date(thisM)

    def NextMonth(self):
        nextM = time.strptime(self.date, "%Y-%m-%d")
        return self.add_months(nextM, 1)

    def MonthName(self):
        month = time.strptime(self.date, "%Y-%m-%d").tm_mon
        year = time.strptime(self.date, "%Y-%m-%d").tm_year
        return {
             1 : 'January',
             2 : 'February',
             3 : 'March',
             4 : 'April',
             5 : 'May',
             6 : 'June',
             7 : 'July',
             8 : 'August',
             9 : 'September',
            10 : 'October',
            11 : 'November',
            12 : 'December',
         }[month] + ' ' + str(year)
