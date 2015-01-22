#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta

class Search:

#    def __init__(self):

    def SimilarExpenses(self, description):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = "select expenses.eid,  date, description, printf('%.2f',amount), cd.name from expenses,classifications c, classificationdef cd where (description like '%{0}%' or name like '%{0}%') and expenses.eid = c.eid and c.cid=cd.cid order by date desc".format(description)
        cursor = conn.execute(query)
        return cursor

