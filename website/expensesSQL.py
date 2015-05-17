#!/usr/bin/python

def _expenseSQL():
    sql = 'select date, description, printf("%.2f", amount), cd.name, e.eid, c.confirmed, tag, amountfx, ccyfx, fxrate, commission '
    return sql

def _baseSQL():
    sql = _expenseSQL() +  ' from expenses e left join tagged t on e.eid = t.eid, classifications c, classificationdef cd where  e.eid = c.eid and c.cid = cd.cid '
    return sql;

def getExpense(eid):
    sql = _baseSQL() + ' and e.eid = {0};'
    return sql.format(eid)

def getAllOneMonthsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') order by date desc;'
    print sql.format(date)
    return sql.format(date)

def getSomeOneMonthsExpenses(date):
    sql = _baseSQL() + ' and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') and (cd.isexpense or not c.confirmed) order by date desc;'
    return sql.format(date)

def getSimilarExpenses(search):
    sql = _expenseSQL() + " from expenses e left join classifications c on e.eid = c.eid left join classificationdef cd on c.cid = cd.cid left join tagged t on e.eid = t.eid left join documentexpensemapping d on e.eid = d.eid  where (e.description like '%{0}%' or cd.name like '%{0}%') order by e.date desc;"
    return sql.format(search)

def getRawLines(eid):
    sql = 'select r.rid, r.rawStr from rawdata r, ExpenseRawMapping erm where erm.rid = r.rid and erm.eid={0};'
    return sql.format(eid)

def getDocuments(eid):
    sql = 'select d.did, substr(text, 1, 101) from documents d, DocumentExpenseMapping dem where d.did = dem.did and dem.eid={0};'
    return sql.format(eid)

def getDocument(did):
    sql = 'select filename, text, deleted from documents where did={0}'
    return sql.format(did)

def getMatchingExpenses(did):
    sql = 'select e.eid, e.description from documentexpensemapping d, expenses e where did = {0} and e.eid = d.eid'
    return sql.format(did)

def getNextDocID(did):
    sql = 'select min (did) from documents where did > {0} and deleted = 0'
    return sql.format(did)

def getPreviousDocID(did):
    sql = 'select max (did) from documents where did < {0} and deleted = 0'
    return sql.format(did)
    
