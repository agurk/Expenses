package exrecords

import (
	"b2/backend"
	"b2/components/managed/expenses"
	"b2/webhandler"
	"encoding/json"
	"fmt"
	"net/http"
)

type postData struct {
	Eid     uint64   `json:"eid"`
	Group   uint64   `json:"group"`
	Members []uint64 `json:"members"`
}

type WebHandler struct {
	backend  *backend.Backend
	Path     string
	LongPath string
}

func Instance(path string, backend *backend.Backend) *WebHandler {
	handler := new(WebHandler)
	handler.backend = backend
	handler.Path = path
	handler.LongPath = path + "/"
	return handler
}

func (handler *WebHandler) GetPath() string {
	return handler.Path
}

func (handler *WebHandler) GetLongPath() string {
	return handler.LongPath
}

func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.Header().Set("Access-Control-Allow-Origin", "*")
		groups, err := getSplitwiseGroups(handler.backend.Splitwise.BearerToken)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		json, _ := json.Marshal(groups)
		fmt.Fprintln(w, string(json))
		w.Header().Set("Content-Type", "application/json")
		return
	case "POST":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		data := new(postData)
		err := decoder.Decode(&data)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		exp, _ := handler.backend.Expenses.Get(data.Eid)
		err = addSplitwiseExpense(data, exp.(*expenses.Expense), handler.backend.Splitwise.BearerToken, handler.backend.Splitwise.User)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		err = handler.backend.Expenses.Save(exp)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(404), 404)
		return
	}
}
