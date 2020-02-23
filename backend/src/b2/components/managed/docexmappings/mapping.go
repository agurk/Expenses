package docexmappings

import (
	"b2/errors"
	"b2/manager"
	"sync"
)

// Mapping represents a one-to-one link between a mapping and an expense
type Mapping struct {
	ID        uint64 `json:"id"`
	EID       uint64 `json:"expenseId"`
	DID       uint64 `json:"documentId"`
	Confirmed bool   `json:"confirmed"`
	deleted   bool   `json:-`
	sync.RWMutex
}

// Cast a manager.Thing into a *Mapping or panic
func Cast(thing manager.Thing) *Mapping {
	mapping, ok := thing.(*Mapping)
	if !ok {
		panic("Non mapping passed to function")
	}
	return mapping
}

// Type returns a string representation of the mapping type
func (mapping *Mapping) Type() string {
	return "docexpensemapping"
}

// GetID returns the ID of a mapping
func (mapping *Mapping) GetID() uint64 {
	return mapping.ID
}

// Merge is not implemented for mappings
func (mapping *Mapping) Merge(newThing manager.Thing) error {
	return errors.New("Not implemented", errors.NotImplemented, "mapping.Merge", true)
}

// Overwrite is not implemented for mappings
func (mapping *Mapping) Overwrite(newThing manager.Thing) error {
	return errors.New("Not implemented", errors.NotImplemented, "mapping.Overwrite", true)
}

// Check returns an error if the mapping has been deleted (i.e. you're holding a pointer to
// an old object
func (mapping *Mapping) Check() error {
	mapping.RLock()
	defer mapping.RUnlock()
	if mapping.deleted {
		return errors.New("Mapping deleted", nil, "mapping.Check", true)
	}
	return nil
}
