package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WebHandler struct {
	manager Manager
	path    string
}

func (handler *WebHandler) Initalize(path string, manager Manager) {
	handler.manager = manager
	handler.path = path
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

func (handler *WebHandler) IndividualHandler(w http.ResponseWriter, req *http.Request) {
	idRaw := req.URL.Path[len(handler.path):]
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

	// Save new
	case "POST":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing := handler.manager.NewThing()
		err := decoder.Decode(&thing)
		// TODO: ignore if ID specified
		if err != nil {
			returnError(err, w)
			return
		}
		fmt.Println(thing)
		err = handler.manager.New(thing)
		if err != nil {
			returnError(err, w)
			return
		} else {
			// todo: add location
			//thing.RLock()
			//location := handler.path + strconv.FormatUint(thing.GetID(), 10)
			//thing.RUnlock()
			//w.Header().Set("Location", location)
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
		fmt.Println(thing)
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
			ID uint64 `json:"id"`
		}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		mergeData := new(merge)
		fmt.Println(mergeData)
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
		err = handler.manager.Merge(thing, mergeThing)
		if err != nil {
			returnError(err, w)
			return
		}

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, MERGE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
	default:
		http.Error(w, http.StatusText(405), 405)
	}
}

func (handler *WebHandler) MultipleHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		things, err := handler.manager.GetMultiple(req.URL.Query())
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
