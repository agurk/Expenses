package docexmappings

import (
	"b2/backend"
	"b2/manager"
	"errors"
)

type Query struct {
	ExpenseId  uint64
	DocumentId uint64
}

type MappingManager struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) manager.Manager {
	mm := new(MappingManager)
	mm.initalize(backend)
	general := new(manager.CachingManager)
	general.Initalize(mm)
	return general
}

func (mm *MappingManager) initalize(backend *backend.Backend) {
	mm.backend = backend
}

func (mm *MappingManager) Load(dmid uint64) (manager.Thing, error) {
	return loadMapping(dmid, mm.backend.DB)
}

func (mm *MappingManager) AfterLoad(mapping manager.Thing) error {
	return nil
}

func (mm *MappingManager) Find(query interface{}) ([]uint64, error) {
	var search *Query
	switch query.(type) {
	case *Query:
		search = query.(*Query)
		//	case url.Values:
		//		params := query.(url.Values)
		//		for key, elem := range params {
		//			// Query() returns empty string as value when no value set for key
		//			if len(elem) != 1 || elem[0] == "" {
		//				return nil, errors.New("Invalid search parameter " + key)
		//			}
		//			switch key {
		//			case "expense":
		//				search.ExpenseId, _ = strconv.ParseUint(elem[0], 10, 64)
		//			case "document":
		//				search.DocumentId, _ = strconv.ParseUint(elem[0], 10, 64)
		//			default:
		//				return nil, errors.New("Invalid search parameter " + key)
		//			}
		//		}
	default:
		return nil, errors.New("Unknown type passed to find function")
	}

	return findMappings(search, mm.backend.DB)
}

func (mm *MappingManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (mm *MappingManager) Create(mapp manager.Thing) error {
	mapping, ok := mapp.(*Mapping)
	if !ok {
		return errors.New("Non mapping passed to function")
	}
	// todo: check it's unique?
	return createMapping(mapping, mm.backend.DB)

}

func (mm *MappingManager) Update(mp manager.Thing) error {
	mapping, ok := mp.(*Mapping)
	if !ok {
		return errors.New("Non mapping passed to function")
	}
	return updateMapping(mapping, mm.backend.DB)
}

func (mm *MappingManager) NewThing() manager.Thing {
	return new(Mapping)
}

func (mm *MappingManager) Combine(one, two manager.Thing) error {
	return errors.New("Not implemented")
}

func (mm *MappingManager) Delete(mp manager.Thing) error {
	mapping, ok := mp.(*Mapping)
	if !ok {
		return errors.New("Non mapping passed to function")
	}
	mapping.Lock()
	defer mapping.Unlock()
	err := deleteMapping(mapping, mm.backend.DB)
	if err != nil {
		return nil
	}
	mapping.deleted = true
	return nil
}
