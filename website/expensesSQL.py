#!/usr/bin/python

def getExpense(eid):
    sql = 'select date, description, printf("%.2f", amount), cd.name, e.eid, confirmed, tag, amountfx, ccyfx, fxrate, commission from expenses e left join tagged t on e.eid = t.eid, classifications c, classificationdef cd where e.eid = {0} and e.eid = c.eid and c.cid = cd.cid;'
    return sql.format(eid)
