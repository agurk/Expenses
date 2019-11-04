package main

import (
    "b2/expenses"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "errors"
)

type Env struct {
    manager *expenses.ExManager
}

func returnError (err error, w http.ResponseWriter) {
    switch err.Error() {
    case "404":
        http.Error(w, http.StatusText(404), 404)
    default:
        http.Error(w, err.Error(), 400)
    }
}

func (env *Env) getExpense (eidRaw string) (*expenses.Expense, error) {
    eid, err := strconv.ParseUint(eidRaw, 10, 64)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }

    expense, err := env.manager.GetExpense(eid)
    if err != nil {
        return nil, err
    }

    return expense, nil
}

func (env *Env) expensesHandler(w http.ResponseWriter, req *http.Request) {
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

        expenses, err := env.manager.GetExpenses(from, to)
        if err != nil {
            returnError(err, w)
            return
        }
        w.Header().Set("Content-Type", "application/json")
       // for _, expense := range expenses {
        //    expense.RLock()
        //
            json, _ := json.Marshal(expenses)
            fmt.Fprintln(w, string(json))
          //  expense.RUnlock()
       // }
    default:
        http.Error(w, http.StatusText(405), 405)
    }
}


func (env *Env) expenseHandler(w http.ResponseWriter, req *http.Request) {
    //fmt.Println(req.URL.Path[len("/expenses/"):])
    eidRaw := req.URL.Path[len("/expenses/"):]

    switch req.Method {
    case "GET":
        expense, err := env.getExpense(eidRaw)
        if err != nil {
            returnError(err, w)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        expense.RLock()
        json, err := json.Marshal(expense)
        fmt.Fprintln(w, string(json))
        expense.RUnlock()

    // Save new
    case "POST":
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        var e expenses.Expense
        err := decoder.Decode(&e)
        if err != nil {
            returnError(err, w)
            return
        }
        fmt.Println(e)
        err = env.manager.SaveExpense(&e)
        if err != nil {
            returnError(err, w)
            return
        } else {
            e.RLock()
            location := "/expenses/" + strconv.FormatUint(e.ID, 10)
            e.RUnlock()
            w.Header().Set("Location",location)
            //http.Success(w, http.StatusText(201), 201)
        }

    // replace existing
    case "PUT":
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        var e expenses.Expense
        err := decoder.Decode(&e)
        if err != nil {
            returnError(err, w)
            return
        }
        fmt.Println(e)
        _, err = env.manager.OverwriteExpense(&e)
        if err != nil {
            returnError(err, w)
            return
        }

    // update existing
    case "PATCH":
        expense, err := env.getExpense(eidRaw)
        if err != nil {
            returnError(err, w)
            return
        }
        decoder := json.NewDecoder(req.Body)
        decoder.DisallowUnknownFields()
        expense.Lock()
        err = decoder.Decode(&expense)
        expense.Unlock()
        if err != nil {
            returnError(err, w)
            return
        }
        err = env.manager.SaveExpense(expense)
        if err != nil {
            fmt.Println(err)
            panic(err)
        }

    case "OPTIONS":
        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH")
        w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5000")
        w.Header().Set("Access-Control-Allow-Headers", "content-type")
    default:
        http.Error(w, http.StatusText(405), 405)
    }
}

func main() {
    env := new (Env)
    env.manager = new (expenses.ExManager)
    err := env.manager.Initalize("/home/timothy/src/Expenses/expenses.db")
    if err != nil {
        log.Panic(err)
    }

    http.HandleFunc("/expenses/", env.expenseHandler)
    http.HandleFunc("/expenses", env.expensesHandler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
