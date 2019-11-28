package classifications 

import "database/sql"
import "b2/manager"

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

func getClassifications(db *sql.DB) ([]manager.Thing, error) {
    rows, err := db.Query("select cid, name, validfrom, validto, isexpense from classificationdef")
    if err != nil {
        return nil, err
    }
    var classifications []manager.Thing
    defer rows.Close()
    for rows.Next() {
        class := new(dbClassification)
        err = rows.Scan(&class.ID,
                        &class.Description,
                        &class.From,
                        &class.To,
                        &class.Hidden)
        if err != nil {
            return nil, err
        }
        classifications = append(classifications, result2classification(class))

    }
    return classifications, err
}

