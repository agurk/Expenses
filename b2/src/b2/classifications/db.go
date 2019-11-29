package classifications 

import (
    "database/sql"
    "errors"
)

type dbClassification struct {
    ID uint64
    Description string
    Hidden bool
    From string
    To sql.NullString
}

func parseSQLstr(str *sql.NullString) string {
        if !str.Valid {
            return ""
        }
        return str.String
}

func result2classification(result *dbClassification) *Classification {
    classification := new(Classification)
    classification.ID = result.ID
    classification.Description = result.Description
    classification.From = result.From
    classification.To = parseSQLstr(&result.To)
    classification.Hidden = result.Hidden
    return classification
}

func loadClassification(cid uint64, db *sql.DB) (*Classification, error) {
    rows, err := db.Query(`
        select
            cid,
            name,
            validfrom,
            validto,
            isexpense
        from
            classificationdef
        where
            cid = $1`,
            cid)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    dbClass := new(dbClassification)
    if rows.Next() {
        err = rows.Scan(&dbClass.ID,
                        &dbClass.Description,
                        &dbClass.From,
                        &dbClass.To,
                        &dbClass.Hidden)
    } else {
        return nil, errors.New("404")
    }
    if err != nil {
        return nil, err
    }
    return result2classification(dbClass), nil
}


func findClassifications(db *sql.DB) ([]uint64, error) {
    rows, err := db.Query("select cid from classificationdef")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var cids []uint64
    for rows.Next() {
        var cid uint64
        err = rows.Scan(&cid)
        if err != nil {
            return nil, err
        }
        cids = append(cids, cid)
    }
    return cids, err
}

