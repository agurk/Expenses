#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta

class Config:

    def AllClassifications(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select cid, name, validfrom, validto, isexpense from classificationdef';
        return conn.execute(query)

