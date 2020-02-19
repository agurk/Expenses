package analysis

import (
	"b2/moneyutils"
	"b2/webhandler"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// WebHandler is the http webhandler for managing requests for analysis graphs
type WebHandler struct {
	db       *sql.DB
	rates    *moneyutils.FxValues
	path     string
	longpath string
}

// Instance returns an instatiated & configured instance of the webhandler
// used to build analysis graphs
func Instance(path string, db *sql.DB) *WebHandler {
	handler := new(WebHandler)
	handler.db = db
	handler.rates = new(moneyutils.FxValues)
	handler.rates.Initalize(db)
	handler.path = path
	handler.longpath = path + "/"
	return handler
}

// Path returns the path at which the handler is expecting requests
func (handler *WebHandler) Path() string {
	return handler.path
}

//LongPath is the same as Path() but returns the path with a terminating /
func (handler *WebHandler) LongPath() string {
	return handler.longpath
}

// Handle is a webhandler function that can be passed to net/http
func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch req.URL.Path[len("/analysis/"):] {
		case "totals":
			params, err := processParams(req.URL.Query())
			if err != nil {
				webhandler.ReturnError(err, w)
				return
			}
			results, err := totals(params, handler.rates, handler.db)
			if err != nil {
				webhandler.ReturnError(err, w)
				return
			}
			json, _ := json.Marshal(results)
			fmt.Fprintln(w, string(json))
			w.Header().Set("Content-Type", "application/json")
		case "graph":
			params, err := processParams(req.URL.Query())
			if err != nil {
				webhandler.ReturnError(err, w)
				return
			}
			gParams := gInitialise(params)
			results, err := graph(gParams, handler.rates, handler.db)
			if err != nil {
				webhandler.ReturnError(err, w)
				return
			}
			fmt.Fprintln(w, results)
			w.Header().Set("Content-Type", "image/svg+xml")
		case "assets":
			results, err := assets(handler.rates, handler.db)
			if err != nil {
				webhandler.ReturnError(err, w)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json, _ := json.Marshal(results)
			fmt.Fprintln(w, string(json))
		default:
			http.Error(w, http.StatusText(404), 404)
			return
		}
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(405), 405)
	}
}
