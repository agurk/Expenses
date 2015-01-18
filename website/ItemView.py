#!/usr/bin/python

import sqlite3

class ItemView:

    def __init__(self, expenseID):
        self.expenseID = expenseID

    def RawStr(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select rawstr from rawdata, expenserawmapping where rawdata.rid=expenserawmapping.rid and expenserawmapping.eid = {0};'.format(self.expenseID)
        cursor = conn.execute(query)
        rawRows = []
        for row in cursor:
            rawRows.append(row[0])
        return rawRows

    def Classifications(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = "select cid,name from classificationdef,expenses e where e.eid={0} and date(validfrom) <= date(e.date) and (validto = '' or date(validto) >= date(e.date)) order by name".format(self.expenseID)
        cursor = conn.execute(query)
        return cursor

    def Classification(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = "select cid from classifications where eid={0}".format(self.expenseID)
        cursor = conn.execute(query)
        cid = ''
        for row in cursor:
            cid = row[0]
        return cid

    def Amount(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = "select amount from expenses where eid={0}".format(self.expenseID)
        cursor = conn.execute(query)
        cid = ''
        for row in cursor:
            cid = row[0]
        return cid


