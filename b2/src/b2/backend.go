package main

import (
	"b2/analysis"
	"b2/manager"
	"b2/manager/classifications"
	"b2/manager/docexmappings"
	"b2/manager/documents"
	"b2/manager/expenses"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func loadDB(dataSourceName string) (*sql.DB, error) {
	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register("expenses_db",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
					fmt.Println("here", op)
					switch op {
					case sqlite3.SQLITE_INSERT:
						fmt.Println("Notified of insert on db", db, "table", table, "rowid", rowid)
					}
				})
				return nil
			},
		})
	db, err := sql.Open("expenses_db", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	var db *sql.DB
	var err error

	if db, err = loadDB("/home/timothy/src/Expenses/expenses.db"); err != nil {
		log.Panic(err)
	}

	docExMapping := docexmappings.Instance(db)

	docWebManager := new(manager.WebHandler)
	exWebManager := new(manager.WebHandler)
	clWebManager := new(manager.WebHandler)
	analWebManager := new(analysis.WebHandler)

	analWebManager.Initalize(db)
	docWebManager.Initalize("/documents/", documents.Instance(db, docExMapping))
	exWebManager.Initalize("/expenses/", expenses.Instance(db, docExMapping))
	clWebManager.Initalize("/expense_classifications/", classifications.Instance(db))

	http.HandleFunc("/expense_classifications", clWebManager.MultipleHandler)
	http.HandleFunc("/expenses/", exWebManager.IndividualHandler)
	http.HandleFunc("/expenses", exWebManager.MultipleHandler)
	http.HandleFunc("/documents/", docWebManager.IndividualHandler)
	http.HandleFunc("/documents", docWebManager.MultipleHandler)
	http.HandleFunc("/analysis/", analWebManager.Handler)

	//log.Fatal(http.ListenAndServe("localhost:8000", nil))
	log.Fatal(http.ListenAndServeTLS("localhost:8000", "certs/server.crt", "certs/server.key", nil))
}
