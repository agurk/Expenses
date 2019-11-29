package main

import (
    "b2/expenses"
    "b2/classifications"
    "b2/documents"
    "b2/manager"
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


    docsC := new (documents.DocManager)
    docsM := new (manager.Manager)
    docWebManager := new (WebHandler)
    expensesC := new (expenses.ExManager)
    expensesM := new (manager.Manager)
    exWebManager := new (WebHandler)
    clC := new (classifications.ClassificationManager)
    clM := new (manager.Manager)
    clWebManager := new (WebHandler)

    docsC.Initalize(db)
    docsM.Initalize(docsC)
    if err = docWebManager.Initalize("/documents/", docsM); err != nil {
        log.Panic(err)
    }

    expensesC.Initalize(db)
    expensesM.Initalize(expensesC)
    if err = exWebManager.Initalize("/expenses/", expensesM); err != nil {
        log.Panic(err)
    }

    clC.Initalize(db)
    clM.Initalize(clC)
    if err = clWebManager.Initalize("/expense_classifications/", clM); err != nil {
        log.Panic(err)
    }

    http.HandleFunc("/expense_classifications", clWebManager.MultipleHandler)
    http.HandleFunc("/expenses/", exWebManager.IndividualHandler)
    http.HandleFunc("/expenses", exWebManager.MultipleHandler)
    http.HandleFunc("/documents/", docWebManager.IndividualHandler)
    http.HandleFunc("/documents", docWebManager.MultipleHandler)

    //log.Fatal(http.ListenAndServe("localhost:8000", nil))
    log.Fatal(http.ListenAndServeTLS("localhost:8000", "certs/server.crt", "certs/server.key", nil))
}

