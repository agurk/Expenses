package docexmappings

import (
	"b2/errors"
	"database/sql"
	"fmt"
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
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "mapping.loadMapping")
	}
	mapping := new(Mapping)
	if rows.Next() {
		err = rows.Scan(&mapping.DID,
			&mapping.EID,
			&mapping.Confirmed)
		mapping.ID = dmid
	} else {
		return nil, errors.New("Mapping not found", errors.ThingNotFound, "mapping.loadMapping", true)
	}
	if err != nil {
		return nil, errors.Wrap(err, "mapping.loadMapping")
	}
	return mapping, nil
}

func findMappings(query *Query, db *sql.DB) ([]uint64, error) {
	var sqlQuery string
	var id uint64
	if query.ExpenseID > 0 {
		sqlQuery = "select dmid from DocumentExpenseMapping where eid = $1"
		id = query.ExpenseID
	} else if query.DocumentID > 0 {
		sqlQuery = "select dmid from DocumentExpenseMapping where did = $1"
		id = query.DocumentID
	} else {
		return nil, errors.New("no valid idType", nil, "mapping.findMappings", false)
	}
	rows, err := db.Query(sqlQuery, id)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "mapping.findMappings")
	}
	var dmids []uint64
	for rows.Next() {
		var dmid uint64
		err = rows.Scan(&dmid)
		if err != nil {
			return nil, errors.Wrap(err, "mapping.findMappings")
		}
		dmids = append(dmids, dmid)
	}
	return dmids, errors.Wrap(err, "mapping.findMappings")
}

func updateMapping(mapping *Mapping, db *sql.DB) error {
	err := mapping.Check()
	if err != nil {
		return errors.Wrap(err, "mappings.updateMapping")
	}
	mapping.RLock()
	defer mapping.RUnlock()
	// Todo: Check values are legit before writing
	_, err = db.Exec(` update DocumentExpenseMapping
							set eid = $1,
							did = $2,
							confirmed = $3
						where
							dmid = $4`,
		mapping.EID, mapping.DID, mapping.Confirmed, mapping.ID)
	return errors.Wrap(err, "mappings.updateMapping")
}

func createMapping(mapping *Mapping, db *sql.DB) error {
	mapping.Lock()
	defer mapping.Unlock()
	rows, err := db.Query(`select * from DocumentExpenseMapping where eid = $1 and did = $2`, mapping.EID, mapping.DID)
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "mappings.createMapping")
	}
	for rows.Next() {
		return errors.New(fmt.Sprintf("Error creating mapping as existing mapping for expense %d and document %d", mapping.EID, mapping.DID), nil, "mappings.createMapping", true)
	}
	res, err := db.Exec(`
		insert into
			DocumentExpenseMapping (
				did,
				eid,
				confirmed)
			 values
				($1, $2, $3)`,
		mapping.DID, mapping.EID, mapping.Confirmed)
	if err != nil {
		return errors.Wrap(err, "mappings.createMapping")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		mapping.ID = uint64(rid)
	} else {
		return errors.New("Error creating new mapping", errors.InternalError, "mappings.createMapping", false)
	}
	return nil
}

func deleteMapping(mapping *Mapping, db *sql.DB) error {
	mapping.Lock()
	defer mapping.Unlock()
	_, err := db.Exec("delete from DocumentExpenseMapping where dmid = $1", mapping.ID)
	return errors.Wrap(err, "mappings.deleteMapping")
}
