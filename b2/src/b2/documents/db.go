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
	//rows, err := db.Query("select did from documents where deleted = 0")
	rows, err := db.Query(`	select
								d.did
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
            d.deleted
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
			&document.Deleted)
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

func createDocument(e *Document, db *sql.DB) error {
	return nil
}

func updateExpenes(e *Document, db *sql.DB) error {
	return nil
}
