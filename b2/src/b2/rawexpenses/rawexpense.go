package rawexpenses

import (
	"b2/manager"
	"b2/rawexmappings"
	"errors"
	"sync"
)

type RawExpense struct {
	ID        uint64                   `json:"id"`
	Date      string                   `json:"date"`
	Data      string                   `json:"data"`
	AccountID int                      `json:"accountId"`
	Expenses  []*rawexmappings.Mapping `json:"expenses"`
	sync.RWMutex
}

func (ex *RawExpense) Type() string {
	return "rawexpense"
}

func (ex *RawExpense) GetID() uint64 {
	return ex.ID
}

func (ex *RawExpense) Merge(newThing manager.Thing) error {
	return ex.Overwrite(newThing)
}

func (ex *RawExpense) Overwrite(newThing manager.Thing) error {
	rawexpense, ok := newThing.(*RawExpense)
	if !ok {
		return errors.New("Non rawexpense passed to overwrite function")
	}
	rawexpense.RLock()
	ex.Lock()
	ex.Date = rawexpense.Date
	ex.Data = rawexpense.Data
	ex.AccountID = rawexpense.AccountID
	rawexpense.RUnlock()
	ex.Unlock()
	return nil
}

func (ex *RawExpense) Check() error {
	return nil
}
