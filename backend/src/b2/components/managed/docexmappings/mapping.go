package docexmappings

import (
	"b2/errors"
	"b2/manager"
	"sync"
)

type Mapping struct {
	ID        uint64 `json:"id"`
	EID       uint64 `json:"expenseId"`
	DID       uint64 `json:"documentId"`
	Confirmed bool   `json:"confirmed"`
	deleted   bool   `json:-`
	sync.RWMutex
}

func (mapping *Mapping) Type() string {
	return "docexpensemapping"
}

func (mapping *Mapping) GetID() uint64 {
	return mapping.ID
}

func (mapping *Mapping) Merge(newThing manager.Thing) error {
	return errors.New("Not implemented", errors.NotImplemented, "mapping.Merge")
}

func (mapping *Mapping) Overwrite(newThing manager.Thing) error {
	return errors.New("Not implemented", errors.NotImplemented, "mapping.Overwrite")
}

func (mapping *Mapping) Check() error {
	mapping.RLock()
	defer mapping.RUnlock()
	if mapping.deleted {
		return errors.New("Mapping deleted", nil, "mapping.Check")
	}
	return nil
}
