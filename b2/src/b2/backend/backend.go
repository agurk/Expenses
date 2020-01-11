package backend

import (
	"b2/manager"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

type Backend struct {
	Documents       manager.Manager
	Expenses        manager.Manager
	Mappings        manager.Manager
	Classifications manager.Manager
	DB              *sql.DB
}

func Instance(dataSourceName string) *Backend {
	backend := new(Backend)
	err := backend.loadDB(dataSourceName)
	if err != nil {
		panic(err)
	}
	return backend
}

func (backend *Backend) loadDB(dataSourceName string) error {
	sqlite3conn := []*sqlite3.SQLiteConn{}
	sql.Register("expenses_db",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
					//fmt.Println("here", op)
					switch op {
					case sqlite3.SQLITE_INSERT:
						//	fmt.Println("Notified of insert on db", db, "table", table, "rowid", rowid)
					}
				})
				return nil
			},
		})
	db, err := sql.Open("expenses_db", dataSourceName)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	backend.DB = db
	return nil
}
