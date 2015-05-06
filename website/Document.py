#!/usr/bin/python

import sqlite3
import config
import expensesSQL

class Document:

    def Document(self, did):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str
        cursor = conn.execute(expensesSQL.getDocument(did))
        document = {}
        for row in cursor:
            document['did'] = did
            document['filename'] = '/static/data/documents/' + row[0]
            document['text'] = row[1]
            document['deleted'] = row[2]
            self._addMatchingExpenses(document, conn)
            self._addNextDocID(document, conn)
            self._addPreviousDocID(document, conn)
            return document

    def _addMatchingExpenses(self, document, db):
        cursor = db.execute(expensesSQL.getMatchingExpenses(document['did']))
        document['expenses'] = cursor

    def _addNextDocID(self, document, db):
        cursor = db.execute(expensesSQL.getNextDocID(document['did']))
        for row in cursor:
            document['nextID'] = row[0]

    def _addPreviousDocID(self, document, db):
        cursor = db.execute(expensesSQL.getPreviousDocID(document['did']))
        for row in cursor:
            document['previousID'] = row[0]

