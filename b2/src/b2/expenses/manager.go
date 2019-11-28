package expenses

import (
    "database/sql"
    "sync"
    "errors"
)

type ExManager struct {
    db *sql.DB
    expenses exMap
}

type exMap struct {
    sync.RWMutex
    m map[uint64]*Expense
}

func (manager *ExManager) Initalize (db *sql.DB) error {
    manager.db = db
    manager.expenses.m = make(map[uint64]*Expense)
    return nil
}

func (manager *ExManager) GetClassifications() ([]*Classification, error) {
    return getClassifications(manager.db)
}

func (manager *ExManager) GetExpense(eid uint64) (*Expense, error) {
    manager.expenses.RLock()
    if expense, ok := manager.expenses.m[eid]; ok {
        manager.expenses.RUnlock()
        return expense, nil
    }
    manager.expenses.RUnlock()
    expense, err := loadExpense(eid, manager.db)
    if (err != nil ) {
        return nil, err
    }
    err = loadDocuments(expense, eid, manager.db)
    if err == nil && expense != nil {
        manager.expenses.Lock()
        defer manager.expenses.Unlock()
        // check someone hasn't already inserted it while we were creating it
        if  newEx, ok := manager.expenses.m[eid]; ok {
            return newEx, nil
        }
        manager.expenses.m[eid] = expense
    }
    return expense, err
}

func (manager *ExManager) GetExpenses(from, to string) ([]*Expense, error) {
    // create empty array so we return [] not null
    expenses := []*Expense{}
    eids, err := findExpenses(from, to, manager.db)
    for _, eid := range eids {
        expense, err := manager.GetExpense(eid)
        if (err == nil ) {
            expenses = append (expenses, expense)
        }
    }

    //expens, err := manager.GetExpense(1234)
    //foo := []*Expense{expens}
    return expenses, err
}

func (manager *ExManager) SaveExpense(expense *Expense) error {
    oldEx, err := manager.GetExpense(expense.ID)
    if err != nil {
        if err.Error() == "404" {
            err := createExpense(expense, manager.db)
            if err != nil && expense.ID > 0 {
                manager.expenses.Lock();
                defer manager.expenses.Unlock()
                manager.expenses.m[expense.ID] = expense
            }
            return err
        }
        return errors.New("Error loading existing expense")
    } else if expense == oldEx {
        return updateExpenes(expense, manager.db)
    } else {
        return errors.New("Expense pointer different to one in manager")
    }
}

func (manager *ExManager) OverwriteExpense(expense *Expense) (*Expense, error) {
    oldEx, err := manager.GetExpense(expense.ID)
    if err != nil {
        return nil, errors.New("Error loading existing expense")
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
    return oldEx, manager.SaveExpense(oldEx)
}

