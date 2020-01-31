package webhandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func ReturnError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	switch err.Error() {
	case "404":
		http.Error(w, http.StatusText(404), 404)
	default:
		http.Error(w, err.Error(), 400)
	}
}

// Gets the id from the path of an incoming request
// assuming the format is /path/id and id is a uint64
func GetID(req *http.Request, path string) (uint64, error) {
	if len(req.URL.Path) <= len(path) {
		return 0, errors.New("No id specified")
	}
	id, err := strconv.ParseUint(req.URL.Path[len(path):], 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
