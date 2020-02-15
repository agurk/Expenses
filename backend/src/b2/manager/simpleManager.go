package manager

import (
	"b2/errors"
	"fmt"
)

type SimpleManager struct {
	component Component
}

func (m *SimpleManager) Initalize(component Component) {
	m.component = component
}

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

func (m *SimpleManager) Merge(thing, thingToMerge Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "simpleManager.Merge", true)
}

// Delete invokes the delete function of the component, presumably to delete
// the thing where appropriate
func (m *SimpleManager) Delete(thing Thing) error {
	err := m.component.Delete(thing)
	return errors.Wrap(err, "simpleManager.Delete")
}

// overwrite the existing version of the thing with the new version provided to it
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

func (m *SimpleManager) NewThing() Thing {
	return m.component.NewThing()
}

func (m *SimpleManager) Process(id uint64) {
	m.component.Process(id)
}

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
