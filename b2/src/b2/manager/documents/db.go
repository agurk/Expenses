package documents

import "database/sql"
import "fmt"
import "errors"

func cleanDate(date string) string {
	// horrible hack
	if date == "" {
		return date
	}
	return date[0:len("1234-12-12")]
}

func findDocuments(db *sql.DB) ([]uint64, error) {
	rows, err := db.Query(`
		select
			distinct(d.did)
		from 
			documents d
		left join
			DocumentExpenseMapping dem on d.did = dem.did
		where
			not d.deleted
			and 
				(not dem.confirmed 
				or dem.confirmed is null)
		order by
			d.did desc`)
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
			d.filesize
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
			&document.Filesize)
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
				deleted
				)
			values ($1, $2, $3, $4, $5)`,
		d.Filename,
		d.Date,
		d.Text,
		d.Filesize,
		d.Deleted)
	if err != nil {
		return nil
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
			deleted = $5
		where
			did = $6`,
		d.Filename,
		d.Date,
		d.Text,
		d.Filesize,
		d.Deleted,
		d.ID)
	return err
}

func deleteDocument(d *Document, db *sql.DB) error {
	// Assuming we're getting a locked document
	_, err := db.Exec(`delete from documents where did = $1`, d.ID)
	return err
}
