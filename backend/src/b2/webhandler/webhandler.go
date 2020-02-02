package webhandler

import (
	"b2/errors"
	"fmt"
	"net/http"
	"strconv"
)

func ReturnError(err error, w http.ResponseWriter) {
	if err == nil {
		panic("nil error sent to Return Error function")
		return
	}
	fmt.Println("Error servicing request:", err)
	if e, ok := err.(*errors.Error); ok {
		fmt.Println("Op Stack: ", e.OpStack())
	}
	w.Header().Set("Content-Type", "Content-Type: text/html; charset=UTF-8")
	switch errors.ErrorType(err) {
	case errors.ThingNotFound:
		http.Error(w, http.StatusText(404), 404)
	case errors.NotImplemented:
		http.Error(w, http.StatusText(501), 501)
	case errors.InternalError:
		http.Error(w, http.StatusText(500), 500)
	default:
		http.Error(w, err.Error(), 400)
	}
}

// Gets the id from the path of an incoming request
// assuming the format is /path/id and id is a uint64
func GetID(req *http.Request, path string) (uint64, error) {
	if len(req.URL.Path) <= len(path) {
		return 0, errors.New("No ID", errors.NoID, "webhandler.GetID")
	}
	id, err := strconv.ParseUint(req.URL.Path[len(path):], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "webhandler.GetID")
	}
	return id, nil
}
