package docexmappings

import (
    "database/sql"
    "errors"
)

func loadMapping(dmid uint64, db *sql.DB) (*Mapping, error) {
    rows, err := db.Query(`
        select
            did,
            eid,
            confirmed
        from
            DocumentExpenseMapping
        where
            dmid = $1`,
            dmid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    mapping := new(Mapping)
    if rows.Next() {
        err = rows.Scan(&mapping.DID,
                        &mapping.EID,
                        &mapping.Confirmed)
        mapping.ID = dmid
    } else {
        return nil, errors.New("404")
    }
    if err != nil {
        return nil, err
    }
    return mapping, nil
}


func findMappings(idType string, id uint64, db *sql.DB) ([]uint64, error) {
    query := ""
    switch idType {
    case "expense":
        query = "select dmid from DocumentExpenseMapping where eid = $1"
    case "document":
        query = "select dmid from DocumentExpenseMapping where did = $1"
    default:
        return nil, errors.New("no valid idType")
    }
    rows, err := db.Query(query,id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var dmids []uint64
    for rows.Next() {
        var dmid uint64
        err = rows.Scan(&dmid)
        if err != nil {
            return nil, err
        }
        dmids = append(dmids, dmid)
    }
    return dmids, err
}

