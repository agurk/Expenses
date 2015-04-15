#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta

class Expense:

    def Expense(self, eid):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select date, description, printf("%.2f", amount), cd.name, e.eid, confirmed, tag, amountfx, ccyfx, fxrate, commission from expenses e left join tagged t on e.eid = t.eid, classifications c, classificationdef cd where e.eid = {0} and e.eid = c.eid and c.cid = cd.cid;'.format(eid)
        cursor = conn.execute(query)
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
