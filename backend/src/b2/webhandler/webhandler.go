package webhandler

import (
	"b2/errors"
	"net/http"
	"strconv"
)

// ReturnError returns an suitable error to the client whilst logging all known details
// onto the console
func ReturnError(err error, w http.ResponseWriter) {
	if err == nil {
		panic("nil error sent to Return Error function")
		return
	}
	errors.Print(err)
	w.Header().Set("Content-Type", "Content-Type: text/html; charset=UTF-8")
	var errText string
	var errCode int
	switch errors.ErrorType(err) {
	case errors.ThingNotFound:
		errText = http.StatusText(404)
		errCode = 404
	case errors.NotImplemented:
		errText = http.StatusText(501)
		errCode = 501
	case errors.InternalError:
		errText = http.StatusText(500)
		errCode = 500
	case errors.Forbidden:
		errText = http.StatusText(403)
		errCode = 403
	default:
		errText = http.StatusText(400)
		errCode = 400
	}
	if err2, ok := err.(*errors.Error); ok && err2.Public {
		errText += ". " + err.Error()
	}
	http.Error(w, errText, errCode)
}

// GetID gets the id from the path of an incoming request
// assuming the format is /path/id and id is a uint64
func GetID(req *http.Request, path string) (uint64, error) {
	if len(req.URL.Path) <= len(path) {
		return 0, errors.New("No ID provided for resource", errors.NoID, "webhandler.GetID", true)
	}
	id, err := strconv.ParseUint(req.URL.Path[len(path):], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "webhandler.GetID")
	}
	return id, nil
}
