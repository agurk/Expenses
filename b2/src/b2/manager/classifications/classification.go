package classifications

import (
	"b2/manager"
	"errors"
	"sync"
)

type Classification struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	From        string `json:"from"`
	To          string `json:"to"`
	sync.RWMutex
}

func (classification *Classification) Type() string {
	return "classification"
}

func (classification *Classification) GetID() uint64 {
	return classification.ID
}

func (classification *Classification) Merge(newThing manager.Thing) error {
	return errors.New("Not implemented")
}

func (classification *Classification) Overwrite(newThing manager.Thing) error {
	return errors.New("Not implemented")
}

func (classification *Classification) Check() error {
	return nil
}
