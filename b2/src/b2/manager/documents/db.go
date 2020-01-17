package documents

import "database/sql"
import "fmt"
import "errors"

func cleanDate(date string) string {
	// todo improve date handling
	if date == "" {
		return date
	}
	return date[0:10]
}

func findDocuments(query *Query, db *sql.DB) ([]uint64, error) {
	dbQuery := `
		select
			distinct(d.did)
		from 
			documents d
		left join
			DocumentExpenseMapping dem on d.did = dem.did
		where
			not d.deleted`
	if query.Starred == true {
		dbQuery += ` and d.Starred`
	} else {
		dbQuery += ` and not d.Starred`
	}

	if query.Archived == true {
		dbQuery += ` and d.archived`
	} else {
		dbQuery += ` and not d.archived`
	}

	if query.Unmatched == true {
		dbQuery += ` and 
				(not dem.confirmed 
				or dem.confirmed is null)`
	}
	dbQuery += ` order by d.did desc`
	rows, err := db.Query(dbQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dids []uint64
	for rows.Next() {
		var did uint64
		err = rows.Scan(&did)
		if err != nil {
			return nil, err
		}
		dids = append(dids, did)
	}
	return dids, err
}

func loadDocument(did uint64, db *sql.DB) (*Document, error) {
	rows, err := db.Query(`
        select
            d.date,
            d.filename,
            d.text,
            d.deleted,
			d.filesize,
			d.starred,
			d.archived
        from
            documents d
        where
            d.did = $1`,
		did)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	document := new(Document)
	if rows.Next() {
		err = rows.Scan(&document.Date,
			&document.Filename,
			&document.Text,
			&document.Deleted,
			&document.Filesize,
			&document.Starred,
			&document.Archived)
		document.ID = did
	} else {
		return nil, errors.New("404")
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return document, nil
}

func createDocument(d *Document, db *sql.DB) error {
	err := d.Check()
	if err != nil {
		return err
	}
	d.Lock()
	defer d.Unlock()
	rows, err := db.Query(`
		select
			did
		from
			documents
		where
			filename = $1
			and filesize = $2`,
		d.Filename, d.Filesize)
	if err != nil {
		return nil
	}
	if rows.Next() {
		return errors.New("existing document, not saving")
	}
	res, err := db.Exec(`
		insert into
			documents (
				filename,
				date,
				text,
				filesize,
				deleted,
				starred,
				archived
				)
			values ($1, $2, $3, $4, $5, $6, $7)`,
		d.Filename,
		d.Date,
		d.Text,
		d.Filesize,
		d.Deleted,
		d.Starred,
		d.Archived)
	if err != nil {
		return err
	}
	did, err := res.LastInsertId()
	if err == nil && did > 0 {
		d.ID = uint64(did)
	} else {
		return errors.New("Error saving new document")
	}
	return nil
}

func updateDocument(d *Document, db *sql.DB) error {
	d.RLock()
	defer d.RUnlock()
	_, err := db.Exec(`
		update
			documents
		set
			filename = $1,
			date = $2,
			text = $3,
			filesize = $4,
			deleted = $5,
			starred = $6,
			archived = $7
		where
			did = $8`,
		d.Filename,
		d.Date,
		d.Text,
		d.Filesize,
		d.Deleted,
		d.Starred,
		d.Archived,
		d.ID)
	return err
}

func deleteDocument(d *Document, db *sql.DB) error {
	// Assuming we're getting a locked document
	_, err := db.Exec(`delete from documents where did = $1`, d.ID)
	return err
}
