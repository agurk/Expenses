package manager

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type CachingManager struct {
	component ManagerComponent
	thingMap  map[uint64]Thing
	sync.RWMutex
}

func (m *CachingManager) Initalize(component ManagerComponent) {
	m.thingMap = make(map[uint64]Thing)
	m.component = component
}

// Get a single thing by id
func (m *CachingManager) Get(id uint64) (Thing, error) {
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
		if err := thing.Check(); err != nil {
			return nil, err
		}
		// To think about: using the id specifed as an arg, rather than the things ID
		m.thingMap[id] = thing
		err = m.component.AfterLoad(thing)
	}
	return thing, err
}

func (m *CachingManager) Find(params interface{}) ([]Thing, error) {
	// create empty array so we return [] not null
	things := []Thing{}
	ids, err := m.component.Find(params)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		thing, err := m.Get(id)
		if err == nil {
			things = append(things, thing)
		} else {
			fmt.Println(id, err.Error())
		}
	}
	return things, err
}

func (m *CachingManager) New(thing Thing) error {
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

func (m *CachingManager) Save(thing Thing) error {
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

func (m *CachingManager) Merge(thing, thingToMerge Thing, params string) error {
	err := m.component.Combine(thing, thingToMerge, params)
	if err != nil {
		return errors.New("Error merging things")
	}
	// deal with errors below
	err = m.Save(thing)
	if err != nil {
		return err
	}
	err = m.Delete(thingToMerge)
	if err != nil {
		return err
	}
	delete(m.thingMap, thingToMerge.GetID())
	return nil
}

func (m *CachingManager) Delete(thing Thing) error {
	err := m.component.Delete(thing)
	delete(m.thingMap, thing.GetID())
	return err
}

// overwrite the existing version of the thing with the new version provided to it
func (m *CachingManager) Overwrite(thing Thing) (Thing, error) {
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

func (m *CachingManager) NewThing() Thing {
	return m.component.NewThing()
}

func (m *CachingManager) Process(id uint64) {
	m.component.Process(id)
}

func (m *CachingManager) LoadDeps(id uint64) {
	thing, err := m.Get(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = m.component.AfterLoad(thing)
	if err != nil {
		fmt.Println(err)
	}
}
