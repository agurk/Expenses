package rawexpenses 

import (
    "sync"
    "b2/rawexmappings"
)

type RawExpense struct {
    ID uint64 `json:"id"`
    Date string `json:"date"`
    Data string `json:"data"`
    AccountID int `json:"accountId"`
    Expenses []*rawexmappings.Mapping `json:"expenses"`
    sync.RWMutex
}

func (ex *RawExpense) Type() string {
    return "rawexpense"
}

func (ex *RawExpense) GetID() uint64 {
    return ex.ID
}

