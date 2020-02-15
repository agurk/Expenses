package classifications

import (
	"b2/errors"
	"database/sql"
	"fmt"
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
		return nil, errors.New("Classification not found", errors.ThingNotFound, "classifications.loadClassification", true)
	}
	if err != nil {
		return nil, err
	}
	classification.From = cleanDate(classification.From)
	classification.To = cleanDate(classification.To)
	return classification, nil
}

func findClassifications(query *Query, db *sql.DB) ([]uint64, error) {
	var args []interface{}
	dbQuery := "select cid from classificationdef "
	where := false
	if query.From != "" || query.Date != "" {
		dbQuery += " where strftime(ValidFrom) <= strftime($1) "
		where = true
		if query.Date == "" {
			args = append(args, query.From)
		} else {
			args = append(args, query.Date)
		}
	}
	if query.To != "" || query.Date != "" {
		if where {
			dbQuery += " and "
		} else {
			dbQuery += " where "
		}
		if query.Date == "" {
			args = append(args, query.To)
		} else {
			args = append(args, query.Date)
		}
		dbQuery += fmt.Sprintf(` ( ValidTo = "" or  strftime(ValidTo) >= strftime($%d))`, len(args))
	}
	rows, err := db.Query(dbQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "classifications.findClassifications (dbQuery)")
	}
	defer rows.Close()
	var cids []uint64
	for rows.Next() {
		var cid uint64
		err = rows.Scan(&cid)
		if err != nil {
			return nil, errors.Wrap(err, "classifications.findClassifications (rows.Scan)")
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
		return errors.Wrap(err, "classifications.createClassification")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		classification.ID = uint64(rid)
	} else if err == nil {
		return errors.New("Error creating new classification", errors.InternalError, "classifications.createClassification", false)
	}
	return errors.Wrap(err, "classifications.createClassification")
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
	return errors.Wrap(err, "classifications.updateClassification")
}

func deleteClassification(c *Classification, db *sql.DB) error {
	rows, err := db.Query("select count(*) from classifications where cid = $1", c.ID)
	if err != nil {
		return errors.Wrap(err, "classifications.deleteClassification (count)")
	}
	defer rows.Close()
	for rows.Next() {
		var count uint64
		err = rows.Scan(&count)
		if err != nil {
			return errors.Wrap(err, "classifications.deleteClassification(count rows.Scan)")
		}
		if count > 0 {
			return errors.New(fmt.Sprintf("Cannot delete classification as it's being used by %d expenses", count),
				nil, "classifications.deleteClassification", true)
		}
	}
	_, err = db.Exec(`
		delete from
			classificationdef
		where
			cid = $1`,
		c.ID)
	return errors.Wrap(err, "classifications.deleteClassification (delete)")
}
