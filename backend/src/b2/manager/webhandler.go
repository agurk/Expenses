package manager

import (
	"b2/errors"
	"b2/webhandler"
	"encoding/json"
	"fmt"
	"net/http"
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

func (handler *WebHandler) getThing(req *http.Request) (Thing, error) {
	id, err := webhandler.GetID(req, handler.LongPath)
	if err != nil {
		return nil, errors.Wrap(err, "webhandler.getThing")
	}

	thing, err := handler.manager.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "webhandler.getThing")
	}

	return thing, nil
}

func (handler *WebHandler) Handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch req.Method {
	case "GET":
		thing, err := handler.getThing(req)
		if err != nil {
			// assuming if no ID given in the path then the user wanted to perform a request
			// against multiple things
			if errors.ErrorType(err) == errors.NoID {
				things, err := handler.manager.Find(req.URL.Query())
				if err != nil {
					webhandler.ReturnError(err, w)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				for _, thing := range things {
					thing.RLock()
					defer thing.RUnlock()
				}
				json, _ := json.Marshal(things)
				fmt.Fprintln(w, string(json))
			} else {
				webhandler.ReturnError(err, w)
			}
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
			webhandler.ReturnError(err, w)
			return
		}
		err = handler.manager.New(thing)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		thing.RLock()
		w.Header().Set("Location", fmt.Sprintf("%s%d", handler.LongPath, thing.GetID()))
		thing.RUnlock()

	// replace existing
	case "PUT":
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing := handler.manager.NewThing()
		err := decoder.Decode(&thing)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		_, err = handler.manager.Overwrite(thing)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}

	// update existing
	case "PATCH":
		thing, err := handler.getThing(req)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		decoder := json.NewDecoder(req.Body)
		decoder.DisallowUnknownFields()
		thing.Lock()
		err = decoder.Decode(&thing)
		thing.Unlock()
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		err = handler.manager.Save(thing)
		if err != nil {
			errors.Print(err)
			panic(err)
		}

	case "MERGE":
		thing, err := handler.getThing(req)
		if err != nil {
			webhandler.ReturnError(err, w)
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
			webhandler.ReturnError(err, w)
			return
		}
		mergeThing, err := handler.manager.Get(mergeData.ID)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		err = handler.manager.Merge(thing, mergeThing, mergeData.Parameters)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}

	case "DELETE":
		thing, err := handler.getThing(req)
		if err != nil {
			webhandler.ReturnError(err, w)
			return
		}
		err = handler.manager.Delete(thing)
		if err != nil {
			webhandler.ReturnError(err, w)
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
