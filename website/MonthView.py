#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta

class MonthView:

    def __init__(self, date):
        self.date = date
        #date=time.strftime("%Y-%m-%d"

    def get_cursor(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select count (*), classificationdef.name, sum(amount) from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid group by classifications.cid;'.format(self.date)
        cursor = conn.execute(query)
        return cursor

    def IndividualExpenses(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select date, description, amount, classificationdef.name from expenses, classifications, classificationdef where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid order by date desc;'.format(self.date)
        cursor = conn.execute(query)
        return cursor

    def add_months(self, sourcedate, months):
        month = sourcedate.tm_mon - 1 + months
        year = sourcedate.tm_year + month / 12
        month = month % 12 + 1
        day = 1
        return datetime.date(year,month,day)

    def PreviousMonth(self):
        previous = time.strptime(self.date, "%Y-%m-%d")
        return self.add_months(previous, -1)

    def NextMonth(self):
        nextM = time.strptime(self.date, "%Y-%m-%d")
        return self.add_months(nextM, 1)
