package expenses

import (
    "b2/docexmappings"
    "b2/rawexmappings"
    "net/url"
    "database/sql"
    "errors"
    "b2/manager"
    "fmt"
    "strconv"
)

type ExManager struct {
    db *sql.DB
    docMappings *manager.Manager
    rawMappings *manager.Manager
}

func Instance(db *sql.DB, docMappings *manager.Manager, rawMappings *manager.Manager) *manager.Manager {
    em := new (ExManager)
    em.db = db
    em.docMappings = docMappings
    em.rawMappings = rawMappings
    general := new (manager.Manager)
    general.Initalize(em)
    return general
}

func (em *ExManager) Load(eid uint64) (manager.Thing, error) {
    return loadExpense(eid, em.db)
}

func (em *ExManager) AfterLoad(ex manager.Thing) (error) {
    expense, ok := ex.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    v := url.Values{}
	v.Set("expense", strconv.FormatUint(expense.ID,10))
    mapps, err := em.docMappings.GetMultiple(v) 
    for _, thing := range mapps {
        mapping, ok := thing.(*(docexmappings.Mapping))
        if !ok {
            return errors.New("Non mapping returned from function")
        }
        expense.Documents = append (expense.Documents, mapping)
    }
    if ( err != nil ) {
        return err
    }
    mapps, err = em.rawMappings.GetMultiple(v) 
    for _, thing := range mapps {
        mapping, ok := thing.(*(rawexmappings.Mapping))
        if !ok {
            return errors.New("Non mapping returned from function")
        }
        expense.Rawdata = append (expense.Rawdata, mapping)
    }
    return err
}

func (em *ExManager) Find(params url.Values) ([]uint64, error) {
    var from, to string
    for key, elem := range params {
        fmt.Println(key)
        fmt.Println(elem)
        // Query() returns empty string as value when no value set for key
        if (len(elem) != 1 || elem[0] == "" ) {
            return nil, errors.New("Invalid query parameter " + key)
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
            return nil, errors.New("Invalid query parameter " + key)
        }
    }

    if ( to == "" || from == "" ) {
        return nil, errors.New("Missing date in date range")
    }

    return findExpenses(from, to, em.db)
}

func (em *ExManager) Create(ex manager.Thing) error {
    expense, ok := ex.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    return createExpense(expense, em.db)
}

func (em *ExManager) Update(ex manager.Thing) error {
    expense, ok := ex.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    return updateExpense(expense, em.db)
}

func (em *ExManager) Merge(from manager.Thing, to manager.Thing) error {
    expense, ok := from.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    oldEx, ok := to.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    expense.RLock()
    oldEx.Lock()
    oldEx.TransactionReference = expense.TransactionReference
    oldEx.Description = expense.Description
    oldEx.DetailedDescription = expense.DetailedDescription
    oldEx.AccountID = expense.AccountID
    oldEx.Date = expense.Date
    oldEx.ProcessDate = expense.ProcessDate
    oldEx.Amount = expense.Amount
    oldEx.Currency = expense.Currency
    oldEx.Commission = expense.Commission
    oldEx.FX = expense.FX
    oldEx.Metadata = expense.Metadata
    //oldEx.Documents = expense.Documents
    expense.RUnlock()
    oldEx.Unlock()
    return nil
}

func (em *ExManager) NewThing() manager.Thing {
    return new(Expense)
}

