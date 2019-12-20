package rawexmappings

import (
    "database/sql"
    "errors"
)

func loadMapping(mid uint64, db *sql.DB) (*Mapping, error) {
    rows, err := db.Query(`
        select
            rid,
            eid
        from
            ExpenseRawMapping
        where
            mid = $1`,
            mid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    mapping := new(Mapping)
    if rows.Next() {
        err = rows.Scan(&mapping.RID,
                        &mapping.EID)
        mapping.ID = mid
    } else {
        return nil, errors.New("404")
    }
    if err != nil {
        return nil, err
    }
    return mapping, nil
}


func findMappings(query *Query, db *sql.DB) ([]uint64, error) {
    var sqlQuery string
    var id uint64
    if query.ExpenseId > 0 {
        sqlQuery = "select mid from ExpenseRawMapping where eid = $1"
        id = query.ExpenseId
    } else if query.RawId > 0 {
        sqlQuery = "select mid from ExpenseRawMapping where rid = $1"
        id = query.RawId
    } else {
        return nil, errors.New("no valid idType")
    }
    rows, err := db.Query(sqlQuery,id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var rids []uint64
    for rows.Next() {
        var rid uint64
        err = rows.Scan(&rid)
        if err != nil {
            return nil, err
        }
        rids = append(rids, rid)
    }
    return rids, err
}

