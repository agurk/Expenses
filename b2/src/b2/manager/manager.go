package manager

import (
	"errors"
	"net/url"
	"strconv"
	"sync"
)

type ManagerInterface interface {
	Load(uint64) (Thing, error)
	AfterLoad(Thing) error
	FindFromUrl(url.Values) ([]uint64, error)
	FindExisting(Thing) (uint64, error)
	Create(Thing) error
	Update(Thing) error
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

type Manager struct {
	component ManagerInterface
	thingMap  map[uint64]Thing
	sync.RWMutex
}

func (m *Manager) Initalize(component ManagerInterface) {
	m.thingMap = make(map[uint64]Thing)
	m.component = component
}

// Get a single thing by id
func (m *Manager) Get(id uint64) (Thing, error) {
	m.RLock()
	if thing, ok := m.thingMap[id]; ok {
		m.RUnlock()
		return thing, nil
	}
	m.RUnlock()
	thing, err := m.component.Load(id)
	if err == nil && thing != nil {
		m.Lock()
		defer m.Unlock()
		// check someone hasn't already inserted it while we were creating it
		if newThing, ok := m.thingMap[id]; ok {
			return newThing, nil
		}
		// To think about: using the id specifed as an arg, rather than the things ID
		m.thingMap[id] = thing
		err = m.component.AfterLoad(thing)
	}
	if err := thing.Check(); err != nil {
		return nil, err
	}
	return thing, err
}

func (m *Manager) GetMultiple(params url.Values) ([]Thing, error) {
	// create empty array so we return [] not null
	things := []Thing{}
	ids, err := m.component.FindFromUrl(params)
	for _, id := range ids {
		thing, err := m.Get(id)
		if err == nil {
			if err := thing.Check(); err != nil {
				return nil, err
			}
			things = append(things, thing)
		}
	}
	return things, err
}

func (m *Manager) New(thing Thing) error {
	if err := thing.Check(); err != nil {
		return err
	}
	existingID, err := m.component.FindExisting(thing)
	if err != nil {
		return err
	} else if existingID > 0 {
		existing, err := m.Get(existingID)
		if err != nil {
			return err
		}
		existing.Merge(thing)
		m.Save(existing)
	} else {
		err := m.component.Create(thing)
		if err != nil && thing.GetID() > 0 {
			m.Lock()
			defer m.Unlock()
			m.thingMap[thing.GetID()] = thing
		}
		return err
	}
	return nil
}

func (m *Manager) Save(thing Thing) error {
	if err := thing.Check(); err != nil {
		return err
	}
	oldThing, err := m.Get(thing.GetID())
	if err != nil {
		return errors.New("Error loading existing " + thing.Type() + " from id " + strconv.FormatUint(thing.GetID(), 10))
	} else if thing == oldThing {
		return m.component.Update(thing)
	} else {
		return errors.New("Conflicting ID '" + strconv.FormatUint(thing.GetID(), 10) + "' tring to save " + thing.Type())
	}
}

func (m *Manager) Overwrite(thing Thing) (Thing, error) {
	if err := thing.Check(); err != nil {
		return nil, err
	}
	// check is right type?
	oldThing, err := m.Get(thing.GetID())
	if err != nil {
		return nil, errors.New("Error loading existing " + thing.Type())
	}
	oldThing.Overwrite(thing)
	return oldThing, m.Save(oldThing)
}

func (m *Manager) NewThing() Thing {
	return m.component.NewThing()
}
