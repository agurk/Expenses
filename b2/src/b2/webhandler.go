package main

import (
	"b2/manager"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WebHandler struct {
	manager *manager.Manager
	path    string
}

func (handler *WebHandler) Initalize(path string, manager *manager.Manager) {
	handler.manager = manager
	handler.path = path
}

func returnError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	fmt.Println(err.Error())
	switch err.Error() {
	case "404":
		http.Error(w, http.StatusText(404), 404)
	default:
		http.Error(w, err.Error(), 400)
	}
}

func (handler *WebHandler) getThing(idRaw string) (manager.Thing, error) {
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
		d := handler.manager.NewThing()
		err := decoder.Decode(&d)
		// TODO: ignore if ID specified
		if err != nil {
			returnError(err, w)
			return
		}
		fmt.Println(d)
		err = handler.manager.New(d)
		if err != nil {
			returnError(err, w)
			return
		} else {
			d.RLock()
			location := handler.path + strconv.FormatUint(d.GetID(), 10)
			d.RUnlock()
			w.Header().Set("Location", location)
		}

	// replace existing
	case "PUT":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		d := handler.manager.NewThing()
		err := decoder.Decode(&d)
		if err != nil {
			returnError(err, w)
			return
		}
		fmt.Println(d)
		_, err = handler.manager.Overwrite(d)
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

	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH")
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
