package main

import (
    "b2/expenses"
    "b2/documents"
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

    exWebManager := new (expenses.WebHandler)
    if err = exWebManager.Initalize(expensesM); err != nil {
        log.Panic(err)
    }

    docsM := new (documents.DocManager)
    if err = docsM.Initalize(db); err != nil {
        log.Panic(err)
    }

    docWebManager := new (documents.WebHandler)
    if err = docWebManager.Initalize(docsM); err != nil {
        log.Panic(err)
    }

    http.HandleFunc("/expenses/", exWebManager.ExpenseHandler)
    http.HandleFunc("/expenses", exWebManager.ExpensesHandler)
    http.HandleFunc("/expense_classifications", exWebManager.ClassificationsHandler)

    http.HandleFunc("/documents/", docWebManager.DocumentHandler)
    //log.Fatal(http.ListenAndServe("localhost:8000", nil))
    log.Fatal(http.ListenAndServeTLS("localhost:8000", "server.crt", "server.key", nil))
}

