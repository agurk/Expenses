package docexmappings

import (
	"b2/manager"
	"errors"
	"sync"
)

type Mapping struct {
	ID        uint64 `json:"id"`
	EID       uint64 `json:"expenseId"`
	DID       uint64 `json:"documentId"`
	Confirmed bool   `json:"confirmed"`
	sync.RWMutex
}

func (mapping *Mapping) Type() string {
	return "docexpensemapping"
}

func (mapping *Mapping) GetID() uint64 {
	return mapping.ID
}

func (mapping *Mapping) Merge(newThing manager.Thing) error {
	return errors.New("Not implemented")
}

func (mapping *Mapping) Overwrite(newThing manager.Thing) error {
	return errors.New("Not implemented")
}

func (mapping *Mapping) Check() error {
	return nil
}
