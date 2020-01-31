package backend

import (
	"b2/manager"
	"b2/webhandler"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"net/http"
)

type Backend struct {
	Documents            manager.Manager
	Expenses             manager.Manager
	Mappings             manager.Manager
	Classifications      manager.Manager
	DB                   *sql.DB
	DocumentsProcessChan chan uint64
	ExpensesProcessChan  chan uint64
	DocumentsDepsChan    chan uint64
	ExpensesDepsChan     chan uint64
	Splitwise            Splitwise
}

type Splitwise struct {
	User        uint64
	BearerToken string
}

func Instance(dataSourceName string) *Backend {
	backend := new(Backend)
	err := backend.loadDB(dataSourceName)
	if err != nil {
		panic(err)
	}
	backend.DocumentsProcessChan = make(chan uint64, 100)
	backend.DocumentsDepsChan = make(chan uint64, 100)
	backend.ExpensesProcessChan = make(chan uint64, 100)
	backend.ExpensesDepsChan = make(chan uint64, 100)
	go backend.docsProcessListen()
	go backend.docsDepsListen()
	go backend.expensesProcessListen()
	go backend.expensesDepsListen()
	return backend
}

func (backend *Backend) expensesDepsListen() {
	for {
		id := <-backend.ExpensesDepsChan
		fmt.Println("expense:", id)
		backend.Expenses.LoadDeps(id)
	}
}
func (backend *Backend) docsDepsListen() {
	for {
		id := <-backend.DocumentsDepsChan
		fmt.Println("document", id)
		backend.Documents.LoadDeps(id)
	}
}
func (backend *Backend) expensesProcessListen() {
	for {
		id := <-backend.ExpensesProcessChan
		fmt.Println("expense:", id)
		backend.Expenses.Process(id)
	}
}
func (backend *Backend) docsProcessListen() {
	for {
		id := <-backend.DocumentsProcessChan
		fmt.Println("document", id)
		backend.Documents.Process(id)
	}
}

func (backend *Backend) Process(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		type dataStruct struct {
			ID   uint64 `json:"id"`
			Type string `json:"type"`
		}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		data := new(dataStruct)
		err := decoder.Decode(&data)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		switch data.Type {
		case "document":
			backend.DocumentsProcessChan <- data.ID
		case "expense":
			backend.ExpensesProcessChan <- data.ID
		default:
			http.Error(w, http.StatusText(400), 400)
		}
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(405), 405)
	}
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
