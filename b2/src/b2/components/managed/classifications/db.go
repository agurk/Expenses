package classifications

import (
	"database/sql"
	"errors"
)

type dbClassification struct {
	ID          uint64
	Description string
	Hidden      bool
	From        string
	To          sql.NullString
}

func parseDate(str *sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}

func cleanDate(in string) string {
	if len(in) < 10 {
		return in
	}
	date := in[:10]
	if date == "0001-01-01" {
		return ""
	}
	return date
}

func result2classification(result *dbClassification) *Classification {
	classification := new(Classification)
	classification.ID = result.ID
	classification.Description = result.Description
	classification.From = cleanDate(result.From)
	classification.To = cleanDate(parseDate(&result.To))
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

func createClassification(classification *Classification, db *sql.DB) error {
	classification.Lock()
	defer classification.Unlock()
	res, err := db.Exec(`insert into
							classificationdef (
								name,
								validFrom,
								validTo,
								isExpense)
							values ($1, $2, $3, $4)`,
		classification.Description,
		classification.From,
		classification.To,
		classification.Hidden)

	if err != nil {
		return err
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		classification.ID = uint64(rid)
	} else {
		return errors.New("Error creating new expense")
	}
	return nil
}

func updateClassification(classification *Classification, db *sql.DB) error {
	classification.RLock()
	defer classification.RUnlock()
	_, err := db.Exec(`
		update
			classificationdef
		set
			name = $1,
			validFrom = $2,
			validTo = $3,
			isExpense = $4
		where
			cid = $5`,
		classification.Description,
		classification.From,
		classification.To,
		classification.Hidden,
		classification.ID)
	return err
}
