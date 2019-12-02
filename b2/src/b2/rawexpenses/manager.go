package rawexpenses

import (
    "net/url"
    "database/sql"
    "errors"
    "b2/manager"
    "b2/rawexmappings"
    "strconv"
)

type RawManager struct {
    db *sql.DB
    expenseMappings *manager.Manager
    processor chan uint64
}

func Instance(db *sql.DB, expenseMappings *manager.Manager, c chan uint64) *manager.Manager {
    rm := new (RawManager)
    rm.db = db
    rm.processor = c
    rm.expenseMappings = expenseMappings
    general := new (manager.Manager)
    general.Initalize(rm)
    return general
}

func (rm *RawManager) Load(eid uint64) (manager.Thing, error) {
    return loadRawExpense(eid, rm.db)
}

func (rm *RawManager) AfterLoad(ex manager.Thing) (error) {
    rawexpense, ok := ex.(*RawExpense)
    if !ok {
        return errors.New("Non rawexpense passed to function")
    }
    v := url.Values{}
    v.Set("raw", strconv.FormatUint(rawexpense.ID,10))
    mapps, err := rm.expenseMappings.GetMultiple(v)
    for _, thing := range mapps {
        mapping, ok := thing.(*(rawexmappings.Mapping))
        if !ok {
            return errors.New("Non mapping returned from function")
        }
        rawexpense.Expenses = append (rawexpense.Expenses, mapping)
    }
    return err
}

func (rm *RawManager) Find(params url.Values) ([]uint64, error) {
    return findRawExpenses(rm.db)
}

func (rm *RawManager) Create(ex manager.Thing) error {
    rawexpense, ok := ex.(*RawExpense)
    if !ok {
        return errors.New("Non rawexpense passed to function")
    }
    err := createRawExpense(rawexpense, rm.db)
    if (err == nil ) {
        rm.processor <- rawexpense.ID
    }
    return err
}

func (rm *RawManager) Update(ex manager.Thing) error {
    rawexpense, ok := ex.(*RawExpense)
    if !ok {
        return errors.New("Non rawexpense passed to function")
    }
    return updateRawExpense(rawexpense, rm.db)
}

func (rm *RawManager) Merge(from manager.Thing, to manager.Thing) error {
    rawexpense, ok := from.(*RawExpense)
    if !ok {
        return errors.New("Non rawexpense passed to function")
    }
    oldEx, ok := to.(*RawExpense)
    if !ok {
        return errors.New("Non rawexpense passed to function")
    }
    rawexpense.RLock()
    oldEx.Lock()
    oldEx.Date = rawexpense.Date
    oldEx.Data = rawexpense.Data
    oldEx.AccountID = rawexpense.AccountID
    rawexpense.RUnlock()
    oldEx.Unlock()
    return nil
}

func (rm *RawManager) NewThing() manager.Thing {
    return new(RawExpense)
}

