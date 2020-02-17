package backend

import "b2/manager"

type component interface {
	Process(uint64)
	Load(uint64) (manager.Thing, error)
	AfterLoad(manager.Thing) error
}

type docmgr interface {
	ReclassifyAll() error
}
