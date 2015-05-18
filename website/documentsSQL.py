#!/usr/bin/python

def _documentSQL():
    sql = 'select did, date, filename, text, textModDate, deleted '
    return sql

def _baseSQL():
    sql = _documentSQL() +  ' from documents '
    return sql

def getDocument(did):
    sql= _baseSQL() + ' where did={0};'
    return sql.format(did)

def getExpenses(did):
    sql = 'select e.eid, date, description, confirmed, dmid from expenses e, documentexpensemapping dem where e.eid = dem.eid and dem.did = {0};'
    return sql.format(did)

def getAllDocuments():
    return _baseSQL() + ' where not deleted order by did desc;'

def getNextDocID(did):
    sql = 'select min (did) from documents where did > {0} and deleted = 0'
    return sql.format(did)

def getPreviousDocID(did):
    sql = 'select max (did) from documents where did < {0} and deleted = 0'
    return sql.format(did)
