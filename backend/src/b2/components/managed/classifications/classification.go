package classifications

import (
	"b2/errors"
	"b2/manager"
	"sync"
)

// Classification is a representation of a classification used in the application
type Classification struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	From        string `json:"from"`
	To          string `json:"to"`
	sync.RWMutex
}

// Type returns a string representing what type the object is as it implements manager.Thing
func (classification *Classification) Type() string {
	return "classification"
}

// GetID returns the classifications ID. 0 is a new/unsaved classification
func (classification *Classification) GetID() uint64 {
	return classification.ID
}

// Merge is a sysnonym for Overwrite
func (classification *Classification) Merge(newThing manager.Thing) error {
	return errors.Wrap(classification.Overwrite(newThing), "classifications.Merge")
}

// Overwrite replaces the existing classifications Description, Hidde, From and To fields
// with those from the classification passed to it
func (classification *Classification) Overwrite(newThing manager.Thing) error {
	class, ok := newThing.(*Classification)
	if !ok {
		panic("Non classification passed to overwrite function")
	}
	class.RLock()
	classification.Lock()
	defer classification.Unlock()
	defer class.RUnlock()
	classification.Description = class.Description
	classification.Hidden = class.Hidden
	classification.From = class.From
	classification.To = class.To
	return nil
}

// Check returns an error if the From date is not set
func (classification *Classification) Check() error {
	// todo improve date check
	if len(classification.From) < 10 {
		return errors.New("Invalid from date. Must be specified", nil, "classification.Check", true)
	}
	return nil
}
