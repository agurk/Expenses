package manager

import (
	"b2/errors"
	"fmt"
	"sync"
)

// CachingManager is desigend to hold a copy of a Thing once it has been created and perform
// standard CRUD operations on it
type CachingManager struct {
	component Component
	thingMap  map[uint64]Thing
	sync.RWMutex
}

// Initalize sets up the manager
func (m *CachingManager) Initalize(component Component) {
	m.thingMap = make(map[uint64]Thing)
	m.component = component
}

// Component returns the component for the manager
func (m *CachingManager) Component() Component {
	return m.component
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
			return nil, errors.Wrap(err, "cachingManager.Get")
		}
		// To think about: using the id specifed as an arg, rather than the things ID
		m.thingMap[id] = thing
		err = m.component.AfterLoad(thing)
	}
	return thing, errors.Wrap(err, "cachingManager.Get")
}

// Find returns a slice of things as defined by the params. The params are specific
// to each component
func (m *CachingManager) Find(params interface{}) ([]Thing, error) {
	// create empty array so we return [] not null
	things := []Thing{}
	ids, err := m.component.Find(params)
	if err != nil {
		return nil, errors.Wrap(err, "cachingManager.Find")
	}
	for _, id := range ids {
		thing, err := m.Get(id)
		if err == nil {
			things = append(things, thing)
		} else {
			fmt.Println(id, err.Error())
		}
	}
	return things, errors.Wrap(err, "cachingManager.Find")
}

// New will attempt to find a matching existing thing based on criteria in each
// component. If there are matches it will merge it, otherwise it will save the
// thing as new
func (m *CachingManager) New(thing Thing) error {
	if err := thing.Check(); err != nil {
		return errors.Wrap(err, "cachingManger.New ("+thing.Type()+")")
	}
	existingID, err := m.component.FindExisting(thing)
	if err != nil {
		return errors.Wrap(err, "cachingManger.New ("+thing.Type()+")")
	} else if existingID > 0 {
		existing, err := m.Get(existingID)
		if err != nil {
			return errors.Wrap(err, "cachingManger.New ("+thing.Type()+")")
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
		return errors.Wrap(err, "cachingManger.New ("+thing.Type()+")")
	}
	return nil
}

// Save will attempt to save the thing. If the thing has an ID of an already existing thing
// with the same id (but a different object) the save will fail
func (m *CachingManager) Save(thing Thing) error {
	if err := thing.Check(); err != nil {
		return errors.Wrap(err, "cachingManger.Save ("+thing.Type()+")")
	}
	oldThing, err := m.Get(thing.GetID())
	if err != nil {
		return errors.Wrap(err, "cachingManger.Save ("+thing.Type()+")")
	} else if thing == oldThing {
		return m.component.Update(thing)
	} else {
		return errors.New(fmt.Sprintf("Conflicting ID %d trying to save %s", thing.GetID(), thing.Type()), nil, "cachingManager.Find", false)
	}
}

// Merge will attempt to merge two things. Outcome is component dependent
func (m *CachingManager) Merge(thing, thingToMerge Thing, params string) error {
	err := m.component.Combine(thing, thingToMerge, params)
	if err != nil {
		return errors.Wrap(err, "cachingManger.Merge ("+thing.Type()+")")
	}
	// deal with errors below
	err = m.Save(thing)
	if err != nil {
		return errors.Wrap(err, "cachingManger.Merge ("+thing.Type()+")")
	}
	err = m.Delete(thingToMerge)
	if err != nil {
		return errors.Wrap(err, "cachingManger.Merge ("+thing.Type()+")")
	}
	delete(m.thingMap, thingToMerge.GetID())
	return nil
}

// Delete remove the thing from the cache and the component is likely to destroy
// whatever representation it has of it
func (m *CachingManager) Delete(thing Thing) error {
	err := m.component.Delete(thing)
	delete(m.thingMap, thing.GetID())
	return errors.Wrap(err, "cachingManger.Delete")
}

// Overwrite the existing version of the thing with the new version provided to it
func (m *CachingManager) Overwrite(thing Thing) (Thing, error) {
	if err := thing.Check(); err != nil {
		return nil, errors.Wrap(err, "cachingManager.Overwrite")
	}
	// check is right type?
	oldThing, err := m.Get(thing.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "cachingManager.Overwrite")
	}
	oldThing.Overwrite(thing)
	return oldThing, m.Save(oldThing)
}

// NewThing returns a new thing of the type of the component
func (m *CachingManager) NewThing() Thing {
	return m.component.NewThing()
}
