package documents

import (
    "net/http"
    "errors"
    "encoding/json"
    "strconv"
    "fmt"
)

type WebHandler struct {
    manager *DocManager
}

func (handler *WebHandler) Initalize (manager *DocManager) error {
    handler.manager = manager
    return nil
}

func returnError (err error, w http.ResponseWriter) {
    switch err.Error() {
    case "404":
        http.Error(w, http.StatusText(404), 404)
    default:
        http.Error(w, err.Error(), 400)
    }
}

func (handler *WebHandler) getDocument(didRaw string) (*Document, error) {
    did, err := strconv.ParseUint(didRaw, 11, 64)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    document, err := handler.manager.GetDocument(did)
    if err != nil {
        return nil, err
    }

    return document, nil
}

func (handler *WebHandler) DocumentHandler(w http.ResponseWriter, req *http.Request) {
    didRaw := req.URL.Path[len("/documents/"):]
    w.Header().Set("Access-Control-Allow-Origin", "*")

    switch req.Method {
    case "GET":
        document, err := handler.getDocument(didRaw)
        if err != nil {
            returnError(err, w)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        document.RLock()
        json, err := json.Marshal(document)
        fmt.Fprintln(w, string(json))
        document.RUnlock()

    // Save new
    case "POST":
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        var d Document 
        err := decoder.Decode(&d)
        if err != nil {
            returnError(err, w)
            return
        }
        fmt.Println(d)
        err = handler.manager.SaveDocument(&d)
        if err != nil {
            returnError(err, w)
            return
        } else {
            d.RLock()
            location := "/documents/" + strconv.FormatUint(d.ID, 10)
            d.RUnlock()
            w.Header().Set("Location",location)
        }

    // replace existing
    case "PUT":
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        var d Document 
        err := decoder.Decode(&d)
        if err != nil {
            returnError(err, w)
            return
        }
        fmt.Println(d)
        _, err = handler.manager.OverwriteDocument(&d)
        if err != nil {
            returnError(err, w)
            return
        }

    // update existing
    case "PATCH":
        document, err := handler.getDocument(didRaw)
        if err != nil {
            returnError(err, w)
            return
        }
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        document.Lock()
        err = decoder.Decode(&document)
        document.Unlock()
        if err != nil {
            returnError(err, w)
            return
        }
        err = handler.manager.SaveDocument(document)
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

func (handler *WebHandler) DocumentsHandler(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case "GET":
        var from, to string
        for key, elem := range req.URL.Query() {
            fmt.Println(key)
            fmt.Println(elem)
            // Query() returns empty string as value when no value set for key
            if (len(elem) != 1 || elem[0] == "" ) {
                returnError(errors.New("Invalid query parameter " + key), w)
                return
            }
            switch key {
            case "date":
                // todo: validate date
                from = elem[0]
                to = elem[0]
            case "from":
                from = elem[0]
            case "to":
                to = elem[0]
            default:
                returnError(errors.New("Invalid query parameter " + key), w)
                return
            }
        }

        if ( to == "" || from == "" ) {
            returnError(errors.New("Missing date in date range"), w)
            return
        }

        documents, err := handler.manager.GetDocuments(from, to)
        if err != nil {
            returnError(err, w)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        json, _ := json.Marshal(documents)
        fmt.Fprintln(w, string(json))
    case "OPTIONS":
        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "content-type")
    default:
        http.Error(w, http.StatusText(405), 405)
    }
}

