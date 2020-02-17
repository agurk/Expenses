package suggestions

import (
	"b2/backend"
	"b2/webhandler"
	"encoding/json"
	"fmt"
	"net/http"
)

// WebHandler handles http connections to provide suggestions of changes
// to an existing expense
type WebHandler struct {
	backend  *backend.Backend
	path     string
	longpath string
}

// Instance returns an instantiated webhandler for suggestions
func Instance(path string, backend *backend.Backend) *WebHandler {
	handler := new(WebHandler)
	handler.path = path
	handler.longpath = path + "/"
	handler.backend = backend
	return handler
}

// Path returns the path that the webhandler expects to be serving on
func (handler *WebHandler) Path() string {
	return handler.path
}

// LongPath returns the path the webhandler expects to be serving on appended with a trailing /
func (handler *WebHandler) LongPath() string {
	return handler.longpath
}

// Handle is a standard net/http request handler
func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch req.Method {
	case "GET":
		id, err := webhandler.GetID(req, handler.longpath)
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
