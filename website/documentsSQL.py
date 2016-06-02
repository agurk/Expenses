#!/usr/bin/python

def _documentSQL():
    sql = 'select distinct(d.did), d.date, d.filename, d.text, d.textModDate, d.deleted '
    return sql

def _baseSQL():
    sql = _documentSQL() +  ' from documents d '
    return sql

def _mappedSQL():
    sql = _baseSQL() + 'left join DocumentExpenseMapping dem on d.did = dem.did '
    return sql

def getDocument(did):
    sql= _baseSQL() + ' where did={0};'
    return sql.format(did)

def getExpenses(did):
    sql = 'select e.eid, date, description, confirmed, dmid from expenses e, documentexpensemapping dem where e.eid = dem.eid and dem.did = {0};'
    return sql.format(did)

def getAllDocuments():
    return _baseSQL() + ' where not deleted order by did desc;'

def getUnmappedDocuments():
    return _mappedSQL() + ' where not d.deleted and (not dem.confirmed or dem.confirmed is null) order by d.did desc'

def getNextDocID(did):
    sql = 'select min (did) from documents where did > {0} and deleted = 0'
    return sql.format(did)

def getPreviousDocID(did):
    sql = 'select max (did) from documents where did < {0} and deleted = 0'
    return sql.format(did)

def getLinkedDocs(did):
    sql = 'select d.did, filename from documents d, DocumentExpenseMapping dem where d.did = dem.did and dem.eid=(select eid from DocumentExpenseMapping where did={0}) and d.did <> {0};'
    return sql.format(did)
