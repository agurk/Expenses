package manager

import (
    "net/url"
    "sync"
    "errors"
)

type ManagerInterface interface {
    Load(uint64) (Thing, error)
    Find(url.Values) ([]uint64, error) 
    Create(Thing) error
    Update(Thing) error
    Merge(Thing, Thing) error
    NewThing() Thing
}

type Thing interface {
    Type() string
    RLock()
    RUnlock()
    Lock()
    Unlock()
    GetID() uint64
}

type Manager struct {
    component ManagerInterface
    thingMap map[uint64]Thing
    sync.RWMutex
}

func (m *Manager) Initalize (component ManagerInterface) {
    m.thingMap = make(map[uint64]Thing)
    m.component = component
}

func (m *Manager) Get(id uint64) (Thing, error) {
    m.RLock()
    if thing, ok := m.thingMap[id]; ok {
        m.RUnlock()
        return thing, nil
    }
    m.RUnlock()
    thing, err := m.component.Load(id)
    if err == nil && thing!= nil {
        m.Lock()
        defer m.Unlock()
        // check someone hasn't already inserted it while we were creating it
        if  newThing, ok := m.thingMap[id]; ok {
            return newThing, nil
        }
        m.thingMap[id] = thing
    }
    return thing, err
}

func (m *Manager) GetMultiple(params url.Values) ([]Thing, error) {
    // create empty array so we return [] not null
    things := []Thing{}
    ids, err := m.component.Find(params)
    for _, id := range ids {
        thing, err := m.Get(id)
        if (err == nil ) {
            things = append (things, thing)
        }
    }
    return things, err
}

func (m *Manager) Save(thing Thing) error {
    oldThing, err := m.Get(thing.GetID())
    if err != nil {
        if err.Error() == "404" {
            err := m.component.Create(thing)
            if err != nil && thing.GetID() > 0 {
                m.Lock();
                defer m.Unlock()
                m.thingMap[thing.GetID()] = thing 
            }
            return err
        }
        return errors.New("Error loading existing " + thing.Type())
    } else if thing == oldThing {
        return m.component.Update(thing)
    } else {
        return errors.New(thing.Type() + " pointer different to one in manager")
    }
}

func (m *Manager) Overwrite(thing Thing) (Thing, error) {
    // check is right type?
    oldThing, err := m.Get(thing.GetID())
    if err != nil {
        return nil, errors.New("Error loading existing " + thing.Type())
    }
    m.component.Merge(thing, oldThing)
    return oldThing, m.Save(oldThing)
}

func (m *Manager) NewThing() Thing {
    return m.component.NewThing()
}

