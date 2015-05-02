#!/usr/bin/python

def getExpense(eid):
    sql = 'select date, description, printf("%.2f", amount), cd.name, e.eid, confirmed, tag, amountfx, ccyfx, fxrate, commission from expenses e left join tagged t on e.eid = t.eid, classifications c, classificationdef cd where e.eid = {0} and e.eid = c.eid and c.cid = cd.cid;'
    return sql.format(eid)

def getRawLines(eid):
    sql = 'select r.rid, r.rawStr from rawdata r, ExpenseRawMapping erm where erm.rid = r.rid and erm.eid={0};'
    return sql.format(eid)

def getDocuments(eid):
    sql = 'select d.did, substr(text, 1, 101) from documents d, DocumentExpenseMapping dem where d.did = dem.did and dem.eid={0};'
    return sql.format(eid)
