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

func findDocuments(from, to string, db *sql.DB) ([]uint64, error) {
    rows, err := db.Query("select did from documents where date between $1 and $2", from, to)
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

func loadExpenses(d *Document, db *sql.DB) error {
    rows, err := db.Query(`
        select
            dem.eid,
            dem.confirmed,
            exes.description,
            exes.date
        from
            DocumentExpenseMapping dem,
            Documents docs,
            Expenses exes
        where
            docs.did = dem.did
            and dem.eid = exes.eid
            and dem.did = $1`,
            d.ID)
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        ex := new(Expense)
        err = rows.Scan(&ex.ID, &ex.Confirmed, &ex.Description, &ex.Date)
        if err != nil {
            return err
        }
        ex.Date = cleanDate(ex.Date)
        d.Expenses = append(d.Expenses, ex)

    }
    return err
}

func createDocument(e *Document, db *sql.DB) error {
    return nil 
}

func updateExpenes(e *Document, db *sql.DB) error {
    return nil
}

