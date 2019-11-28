package main

import (
    "b2/expenses"
    "b2/documents"
    "b2/classifications"
    "fmt"
    "log"
    "net/http"
    "github.com/mattn/go-sqlite3"
    "database/sql"
)

func loadDB (dataSourceName string) (*sql.DB, error) {
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

    expensesM := new (expenses.ExManager)
    if err = expensesM.Initalize(db); err != nil {
        log.Panic(err)
    }

    exWebManager := new (WebHandler)
    if err = exWebManager.Initalize("/expenses/", expensesM); err != nil {
        log.Panic(err)
    }

    docsM := new (documents.DocManager)
    if err = docsM.Initalize(db); err != nil {
        log.Panic(err)
    }

    docWebManager := new (WebHandler)
    if err = docWebManager.Initalize("/documents/", docsM); err != nil {
        log.Panic(err)
    }

    clM := new (classifications.ClassificationManager)
    if err = clM.Initalize(db); err != nil {
        log.Panic(err)
    }

    clWebManager := new (WebHandler)
    if err = clWebManager.Initalize("/expense_classifications/", clM); err != nil {
        log.Panic(err)
    }

    http.HandleFunc("/expense_classifications", clWebManager.MultipleHandler)
    http.HandleFunc("/expenses/", exWebManager.IndividualHandler)
    http.HandleFunc("/expenses", exWebManager.MultipleHandler)
    http.HandleFunc("/documents/", docWebManager.IndividualHandler)

    //log.Fatal(http.ListenAndServe("localhost:8000", nil))
    log.Fatal(http.ListenAndServeTLS("localhost:8000", "certs/server.crt", "certs/server.key", nil))
}

