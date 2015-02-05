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
        query = "select date, description, printf('%.2f', amount), cd.name, e.eid, confirmed from expenses e left join classifications c on e.eid = c.eid left join classificationdef cd on c.cid = cd.cid where (e.description like '%{0}%' or cd.name like '%{0}%') order by e.date desc".format(description)
        cursor = conn.execute(query)
        return cursor

