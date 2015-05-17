#!/usr/bin/python

import sqlite3
import config

class MetaData:

    def Classifications(self, eid):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = "select cid,name from classificationdef,expenses e where e.eid={0} and date(validfrom) <= date(e.date) and (validto = '' or date(validto) >= date(e.date)) order by name".format(eid)
        cursor = conn.execute(query)
        return cursor

    def AllClassifications(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select cid, name, validfrom, validto, isexpense from classificationdef';
        return conn.execute(query)
