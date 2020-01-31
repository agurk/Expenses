package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WebHandler struct {
	manager  Manager
	Path     string
	LongPath string
}

// path is expected to be in the format /path (note no trailing /)
func Instance(path string, manager Manager) *WebHandler {
	handler := new(WebHandler)
	handler.manager = manager
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

func returnError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	switch err.Error() {
	case "404":
		http.Error(w, http.StatusText(404), 404)
	default:
		http.Error(w, err.Error(), 400)
	}
}

func (handler *WebHandler) getThing(idRaw string) (Thing, error) {
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	thing, err := handler.manager.Get(id)
	if err != nil {
		return nil, err
	}

	return thing, nil
}

func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) <= len(handler.LongPath) {
		handler.MultipleHandler(w, req)
	} else {
		handler.IndividualHandler(w, req)
	}
	//fmt.Println(req.URL)
	//fmt.Println(req.URL.Path)
	//fmt.Println(req.URL.Query)
}

func (handler *WebHandler) IndividualHandler(w http.ResponseWriter, req *http.Request) {
	idRaw := req.URL.Path[len(handler.LongPath):]
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch req.Method {
	case "GET":
		thing, err := handler.getThing(idRaw)
		if err != nil {
			returnError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		thing.RLock()
		json, err := json.Marshal(thing)
		fmt.Fprintln(w, string(json))
		thing.RUnlock()
		return

	// Save new
	case "POST":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing := handler.manager.NewThing()
		err := decoder.Decode(&thing)
		if err != nil {
			returnError(err, w)
			return
		}
		err = handler.manager.New(thing)
		if err != nil {
			returnError(err, w)
			return
		} else {
			thing.RLock()
			w.Header().Set("Location", fmt.Sprintf("%s%d", handler.LongPath, thing.GetID()))
			thing.RUnlock()
		}

	// replace existing
	case "PUT":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing := handler.manager.NewThing()
		err := decoder.Decode(&thing)
		if err != nil {
			returnError(err, w)
			return
		}
		_, err = handler.manager.Overwrite(thing)
		if err != nil {
			returnError(err, w)
			return
		}

	// update existing
	case "PATCH":
		thing, err := handler.getThing(idRaw)
		if err != nil {
			returnError(err, w)
			return
		}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing.Lock()
		err = decoder.Decode(&thing)
		thing.Unlock()
		if err != nil {
			returnError(err, w)
			return
		}
		err = handler.manager.Save(thing)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

	case "MERGE":
		thing, err := handler.getThing(idRaw)
		if err != nil {
			returnError(err, w)
			return
		}
		type merge struct {
			ID         uint64 `json:"id"`
			Parameters string `json:"parameters"`
		}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		mergeData := new(merge)
		err = decoder.Decode(&mergeData)
		if err != nil {
			returnError(err, w)
			return
		}
		mergeThing, err := handler.manager.Get(mergeData.ID)
		if err != nil {
			returnError(err, w)
			return
		}
		err = handler.manager.Merge(thing, mergeThing, mergeData.Parameters)
		if err != nil {
			returnError(err, w)
			return
		}

	case "DELETE":
		thing, err := handler.getThing(idRaw)
		if err != nil {
			returnError(err, w)
			return
		}
		err = handler.manager.Delete(thing)
		if err != nil {
			returnError(err, w)
			return
		}

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, MERGE, DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(405), 405)
	}
}

func (handler *WebHandler) MultipleHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		things, err := handler.manager.Find(req.URL.Query())
		if err != nil {
			returnError(err, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		for _, thing := range things {
			thing.RLock()
			defer thing.RUnlock()
		}
		json, _ := json.Marshal(things)
		fmt.Fprintln(w, string(json))
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(405), 405)
	}
}
