#!/usr/bin/python

import sqlite3
import time
import datetime
from datetime import date, timedelta
import config

class Config:

    def AllClassifications(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select cid, name, validfrom, validto, isexpense from classificationdef';
        return conn.execute(query)

