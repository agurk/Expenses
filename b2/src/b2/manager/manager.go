package manager

type Manager interface {
	Get(uint64) (Thing, error)
	Find(interface{}) ([]Thing, error)
	New(Thing) error
	Save(Thing) error
	Merge(Thing, Thing) error
	Delete(Thing) error
	Overwrite(Thing) (Thing, error)
	NewThing() Thing
	Process(uint64)
	LoadDeps(uint64)
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
	Find(interface{}) ([]uint64, error)
	FindExisting(Thing) (uint64, error)
	Create(Thing) error
	Update(Thing) error
	Delete(Thing) error
	Combine(Thing, Thing) error
	NewThing() Thing
	Process(uint64)
}
