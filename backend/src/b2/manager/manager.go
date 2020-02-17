package manager

// Manager represents the functionality of a manager
type Manager interface {
	Get(uint64) (Thing, error)
	Find(interface{}) ([]Thing, error)
	// its possible the thing passed to the new function will have the id set
	New(Thing) error
	Save(Thing) error
	Merge(Thing, Thing, string) error
	Delete(Thing) error
	Overwrite(Thing) (Thing, error)
	NewThing() Thing
	GetComponent() Component
}

// Thing describes the features of something that a component manages
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

// Component represents the functionality that is relevant to managing a Thing
// These generally do not have to be implemented and can throw an error instead
type Component interface {
	Load(uint64) (Thing, error)
	AfterLoad(Thing) error
	Find(interface{}) ([]uint64, error)
	// This function is used by the manager to find an existing version of the Thing
	// which it will merge with the posted version. if found. A zero is returned if there is
	// no match
	FindExisting(Thing) (uint64, error)
	Create(Thing) error
	Update(Thing) error
	Delete(Thing) error
	Combine(Thing, Thing, string) error
	NewThing() Thing
}
