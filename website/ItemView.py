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
        rawString = ''
        for row in cursor:
            rawString = rawString + "\n" + row[0]
        return rawString 
