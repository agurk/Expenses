package suggestions

import (
	"b2/backend"
	"b2/webhandler"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebHandler struct {
	backend  *backend.Backend
	Path     string
	LongPath string
}

func Instance(path string, backend *backend.Backend) *WebHandler {
	handler := new(WebHandler)
	handler.Path = path
	handler.LongPath = path + "/"
	handler.backend = backend
	return handler
}

func (handler *WebHandler) GetPath() string {
	return handler.Path
}

func (handler *WebHandler) GetLongPath() string {
	return handler.LongPath
}

func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch req.Method {
	case "GET":
		id, err := webhandler.GetID(req, handler.LongPath)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		suggestions, err := getSuggestions(id, handler.backend)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		json, _ := json.Marshal(suggestions)
		fmt.Fprintln(w, string(json))
		w.Header().Set("Content-Type", "application/json")
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(404), 404)
		return
	}
}
