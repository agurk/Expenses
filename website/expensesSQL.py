#!/usr/bin/python
# -*- coding: utf-8 -*-

def _expenseSQL():
    sql = 'select date, description, amount, e.ccy, cd.name, e.eid, c.confirmed, tag, amountfx, ccyfx, fxrate, commission, e.aid, e.temporary, e.reference, e.modified, e.detaileddescription, e.processdate'
    return sql

def _baseSQL():
    sql = _expenseSQL() +  ' from expenses e left join tagged t on e.eid = t.eid left join classifications c on e.eid = c.eid left join classificationdef cd on c.cid = cd.cid where 1=1 '
    return sql;

def getExpense(eid):
    sql = _baseSQL() + ' and e.eid = {0};'
    return sql.format(eid)

def getAllOneMonthsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') order by date desc, description;'
    return sql.format(date)

def getAllOneYearsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) < date(\'{0}\',\'start of year\',\'+1 year\') order by date desc, description;'
    return sql.format(date)

def getSomeOneMonthsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and (cd.isexpense or not c.confirmed) order by date desc, description;'
    return sql.format(date)

def getSomeOneYearsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) < date(\'{0}\',\'start of year\',\'+1 year\') and (cd.isexpense or not c.confirmed) order by date desc, description;'
    return sql.format(date)

def _classificationSearch(classification):
    if classification== '':
        return ''
    else:
        return " and cd.name = '{0}' ".format(classification)

def getSimilarExpenses(search, classification):
    sql = _expenseSQL() + " from expenses e left join classifications c on e.eid = c.eid left join classificationdef cd on c.cid = cd.cid left join tagged t on e.eid = t.eid where e.eid in (select distinct e.eid from expenses e left join classifications c on e.eid = c.eid left join classificationdef cd on c.cid = cd.cid left join tagged t on e.eid = t.eid left join documentexpensemapping d on e.eid = d.eid  where (e.description like '%{0}%' or cd.name like '%{0}%')) " + _classificationSearch(classification) +" order by e.date desc, description;"
    return sql.format(search)

def getRawLines(eid):
    sql = 'select r.rid, r.rawStr from rawdata r, ExpenseRawMapping erm where erm.rid = r.rid and erm.eid={0};'
    return sql.format(eid)

def getRelatedExpenses(eid):
    sql = 'select e.eid, e.description from expenses e, ExpenseRawMapping erm where e.eid = erm.eid and e.eid <> {0} and erm.rid in (select distinct rid from ExpenseRawMapping where eid = {0});'
    return sql.format(eid)

def getDocuments(eid):
    sql = 'select d.did, filename from documents d, DocumentExpenseMapping dem where d.did = dem.did and dem.eid={0};'
    return sql.format(eid)

def getDocument(did):
    sql = 'select filename, text, deleted from documents where did={0}'
    return sql.format(did)

def getMatchingExpenses(did):
    sql = 'select e.eid, e.description from documentexpensemapping d, expenses e where did = {0} and e.eid = d.eid'
    return sql.format(did)

def getCCYFormats():
    sql = 'select ccy, format from _CCYFormats'
    return sql 

def getFXMonth(month, year):
    date = '{0}-{1}-01'.format(year, month)
    sql = 'select date, ccy1, ccy2, rate from _FXRates where strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')'
    return sql.format(date)
