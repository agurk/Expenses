package classifications

import (
	"database/sql"
	"errors"
)

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
	classification := new(Classification)
	if rows.Next() {
		err = rows.Scan(&classification.ID,
			&classification.Description,
			&classification.From,
			&classification.To,
			&classification.Hidden)
	} else {
		return nil, errors.New("404")
	}
	if err != nil {
		return nil, err
	}
	classification.From = cleanDate(classification.From)
	classification.To = cleanDate(classification.To)
	return classification, nil
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
		return errors.New("Error creating new classification")
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
