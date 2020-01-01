package manager

import (
	"net/url"
)

type Manager interface {
	Get(uint64) (Thing, error)
	GetMultiple(url.Values) ([]Thing, error)
	New(Thing) error
	Save(Thing) error
	Merge(Thing, Thing) error
	Delete(Thing) error
	Overwrite(Thing) (Thing, error)
	NewThing() Thing
}

type Thing interface {
	Type() string
	RLock()
	RUnlock()
	Lock()
	Unlock()
	GetID() uint64
	Merge(Thing) error
	Overwrite(Thing) error
	Check() error
}

type ManagerComponent interface {
	Load(uint64) (Thing, error)
	AfterLoad(Thing) error
	FindFromUrl(url.Values) ([]uint64, error)
	FindExisting(Thing) (uint64, error)
	Create(Thing) error
	Update(Thing) error
	Delete(Thing) error
	Combine(Thing, Thing) error
	NewThing() Thing
}
