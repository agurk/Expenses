package manager

import (
	"b2/errors"
	"fmt"
)

// SimpleManager is an implementation of the Manager interface that directly passes requests
// onto the component with no caching or other processing
type SimpleManager struct {
	component Component
}

// Initalize sets up the manager
func (m *SimpleManager) Initalize(component Component) {
	m.component = component
}

// GetComponent returns the component the manager manages
func (m *SimpleManager) GetComponent() Component {
	return m.component
}

// Get a single thing by id
func (m *SimpleManager) Get(id uint64) (Thing, error) {
	thing, err := m.component.Load(id)
	if err == nil && thing != nil {
		if err := thing.Check(); err != nil {
			return nil, errors.Wrap(err, "simpleManager.Get")
		}
		err = m.component.AfterLoad(thing)
	}
	return thing, errors.Wrap(err, "simpleManager.Get")
}

// Find returns a slice of things that match the params (component dependent)
func (m *SimpleManager) Find(params interface{}) ([]Thing, error) {
	// create empty array so we return [] not null
	things := []Thing{}
	ids, err := m.component.Find(params)
	if err != nil {
		return nil, errors.Wrap(err, "simpleManager.Find")
	}
	for _, id := range ids {
		thing, err2 := m.Get(id)
		if err2 == nil {
			things = append(things, thing)
		} else {
			fmt.Println(id, err2.Error())
		}
	}
	return things, errors.Wrap(err, "simpleManager.Find")
}

// New will see if there if the component knows of an existing version of that
// thing, and if so overwrite it
func (m *SimpleManager) New(thing Thing) error {
	if err := thing.Check(); err != nil {
		return errors.Wrap(err, "simpleManger.New")
	}
	existingID, err := m.component.FindExisting(thing)
	if err != nil {
		return errors.Wrap(err, "simpleManger.New")
	} else if existingID > 0 {
		existing, err := m.Get(existingID)
		if err != nil {
			return errors.Wrap(err, "simpleManger.New")
		}
		existing.Merge(thing)
		m.Save(existing)
	} else {
		return m.component.Create(thing)
	}
	return nil
}

// Save will request the component to save the thing if it's valid
func (m *SimpleManager) Save(thing Thing) error {
	if err := thing.Check(); err != nil {
		return errors.Wrap(err, "simpleManger.Save")
	}
	_, err := m.Get(thing.GetID())
	if err != nil {
		return errors.Wrap(err, "simpleManager.Save")
	}
	return m.component.Update(thing)
}

// Merge is not implemented for simpleManager
func (m *SimpleManager) Merge(thing, thingToMerge Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "simpleManager.Merge", true)
}

// Delete invokes the delete function of the component, presumably to delete
// the thing where appropriate
func (m *SimpleManager) Delete(thing Thing) error {
	err := m.component.Delete(thing)
	return errors.Wrap(err, "simpleManager.Delete")
}

// Overwrite the existing version of the thing with the new version provided to it
func (m *SimpleManager) Overwrite(thing Thing) (Thing, error) {
	if err := thing.Check(); err != nil {
		return nil, errors.Wrap(err, "simpleManager.Overwrite")
	}
	oldThing, err := m.Get(thing.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "simpleManager.Overwrite ("+thing.Type()+")")
	}
	oldThing.Overwrite(thing)
	return oldThing, m.Save(oldThing)
}

// NewThing returns a new type from the component
func (m *SimpleManager) NewThing() Thing {
	return m.component.NewThing()
}

// Process requests the component to process that thing
func (m *SimpleManager) Process(id uint64) {
	m.component.Process(id)
}

// LoadDeps runs the AfterLoad function from the component, which is typically
// involved with dependecies
func (m *SimpleManager) LoadDeps(id uint64) {
	thing, err := m.Get(id)
	if err != nil {
		errors.Print(err)
		return
	}
	err = m.component.AfterLoad(thing)
	if err != nil {
		errors.Print(err)
	}
}
