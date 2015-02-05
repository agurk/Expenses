#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta

class ReadOnlyData:

    def Expense(self, eid):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select date, description, printf("%.2f", amount), cd.name, e.eid, confirmed from expenses e, classifications c, classificationdef cd where e.eid = {0} and e.eid = c.eid and c.cid = cd.cid;'.format(eid)
        cursor = conn.execute(query)
        for row in cursor:
            return row

