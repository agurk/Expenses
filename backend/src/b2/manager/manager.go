package manager

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
	Process(uint64)
	LoadDeps(uint64)
	GetComponent() ManagerComponent
}

// What is returned from a managed component
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

// The implementation of the details for each component share this interface
// with the manager
type ManagerComponent interface {
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
	Process(uint64)
}
