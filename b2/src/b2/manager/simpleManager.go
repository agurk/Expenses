package manager

import (
	"errors"
	"fmt"
	"net/url"
)

type SimpleManager struct {
	component ManagerComponent
}

func (m *SimpleManager) Initalize(component ManagerComponent) {
	m.component = component
}

// Get a single thing by id
func (m *SimpleManager) Get(id uint64) (Thing, error) {
	thing, err := m.component.Load(id)
	if err == nil && thing != nil {
		if err := thing.Check(); err != nil {
			return nil, err
		}
		err = m.component.AfterLoad(thing)
	}
	return thing, err
}

func (m *SimpleManager) GetMultiple(params url.Values) ([]Thing, error) {
	// create empty array so we return [] not null
	things := []Thing{}
	ids, err := m.component.FindFromUrl(params)
	for _, id := range ids {
		thing, err2 := m.Get(id)
		if err2 == nil {
			things = append(things, thing)
		} else {
			// todo: better logging
			fmt.Println(id, err2.Error())
		}
	}
	return things, err
}

func (m *SimpleManager) New(thing Thing) error {
	return errors.New("Not implemented")
}

func (m *SimpleManager) Save(thing Thing) error {
	return errors.New("Not implemented")
}

func (m *SimpleManager) Merge(thing, thingToMerge Thing) error {
	return errors.New("Not implemented")
}

func (m *SimpleManager) Delete(thing Thing) error {
	return errors.New("Not implemented")
}

// overwrite the existing version of the thing with the new version provided to it
func (m *SimpleManager) Overwrite(thing Thing) (Thing, error) {
	return nil, errors.New("Not implemented")
}

func (m *SimpleManager) NewThing() Thing {
	return m.component.NewThing()
}
