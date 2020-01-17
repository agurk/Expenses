package analysis

import (
	"b2/fxrates"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebHandler struct {
	db    *sql.DB
	rates *fxrates.FxValues
}

func (handler *WebHandler) Initalize(db *sql.DB) {
	handler.db = db
	handler.rates = new(fxrates.FxValues)
	handler.rates.Initalize(db)
}

func returnError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	switch err.Error() {
	case "404":
		http.Error(w, http.StatusText(404), 404)
	default:
		http.Error(w, err.Error(), 400)
	}
}

func (handler *WebHandler) Handler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch req.URL.Path[len("/analysis/"):] {
		case "totals":
			params, err := processParams(req.URL.Query())
			if err != nil {
				returnError(err, w)
				return
			}
			results, err := totals(params, handler.rates, handler.db)
			if err != nil {
				returnError(err, w)
				return
			}
			json, _ := json.Marshal(results)
			fmt.Fprintln(w, string(json))
			w.Header().Set("Content-Type", "application/json")
		case "graph":
			params, err := processParams(req.URL.Query())
			if err != nil {
				returnError(err, w)
				return
			}
			gParams := gInitialise(params)
			results, err := graph(gParams, handler.rates, handler.db)
			if err != nil {
				returnError(err, w)
				return
			}
			fmt.Fprintln(w, results)
			w.Header().Set("Content-Type", "image/svg+xml")
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
