package expenses

import (
    "database/sql"
    "sync"
    "errors"
    "b2/manager"
)

type ExManager struct {
    db *sql.DB
    expenses exMap
}

type exMap struct {
    sync.RWMutex
    m map[uint64]*Expense
}

func (em *ExManager) Initalize (db *sql.DB) error {
    em.db = db
    em.expenses.m = make(map[uint64]*Expense)
    return nil
}

func (em *ExManager) Get(eid uint64) (manager.Thing, error) {
    em.expenses.RLock()
    if expense, ok := em.expenses.m[eid]; ok {
        em.expenses.RUnlock()
        return expense, nil
    }
    em.expenses.RUnlock()
    expense, err := loadExpense(eid, em.db)
    if (err != nil ) {
        return nil, err
    }
    err = loadDocuments(expense, eid, em.db)
    if err == nil && expense != nil {
        em.expenses.Lock()
        defer em.expenses.Unlock()
        // check someone hasn't already inserted it while we were creating it
        if  newEx, ok := em.expenses.m[eid]; ok {
            return newEx, nil
        }
        em.expenses.m[eid] = expense
    }
    return expense, err
}

func (em *ExManager) GetMultiple(from, to string) ([]manager.Thing, error) {
    // create empty array so we return [] not null
    expenses := []manager.Thing{}
    eids, err := findExpenses(from, to, em.db)
    for _, eid := range eids {
        expense, err := em.Get(eid)
        if (err == nil ) {
            expenses = append (expenses, expense)
        }
    }

    return expenses, err
}

func (em *ExManager) Save(ex manager.Thing) error {
    expense, ok := ex.(*Expense)
    if !ok {
        return errors.New("Non expense passed to function")
    }
    oldEx, err := em.Get(expense.ID)
    if err != nil {
        if err.Error() == "404" {
            err := createExpense(expense, em.db)
            if err != nil && expense.ID > 0 {
                em.expenses.Lock();
                defer em.expenses.Unlock()
                em.expenses.m[expense.ID] = expense
            }
            return err
        }
        return errors.New("Error loading existing expense")
    } else if expense == oldEx {
        return updateExpenes(expense, em.db)
    } else {
        return errors.New("Expense pointer different to one in manager")
    }
}

func (em *ExManager) Overwrite(ex manager.Thing) (manager.Thing, error) {
    expense, ok := ex.(*Expense)
    if !ok {
        return nil, errors.New("Non expense passed to function")
    }
    oldExpense, err := em.Get(expense.ID)
    if err != nil {
        return nil, errors.New("Error loading existing expense")
    }
    oldEx, ok := oldExpense.(*Expense)
    if !ok {
        return nil, errors.New("Non expense passed to function")
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
    oldEx.Documents = expense.Documents
    expense.RUnlock()
    oldEx.Unlock()
    return oldEx, em.Save(oldEx)
}

