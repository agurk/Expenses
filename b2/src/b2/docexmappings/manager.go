package docexmappings

import (
	"b2/manager"
	"database/sql"
	"errors"
	"net/url"
	"strconv"
)

type Query struct {
	ExpenseId  uint64
	DocumentId uint64
}

type MappingManager struct {
	db *sql.DB
}

func Instance(db *sql.DB) *manager.Manager {
	mm := new(MappingManager)
	mm.initalize(db)
	general := new(manager.Manager)
	general.Initalize(mm)
	return general
}

func (mm *MappingManager) initalize(db *sql.DB) {
	mm.db = db
}

func (mm *MappingManager) Load(dmid uint64) (manager.Thing, error) {
	return loadMapping(dmid, mm.db)
}

func (mm *MappingManager) AfterLoad(mapping manager.Thing) error {
	return nil
}

func (mm *MappingManager) FindFromUrl(params url.Values) ([]uint64, error) {
	var query Query
	for key, elem := range params {
		// Query() returns empty string as value when no value set for key
		if len(elem) != 1 || elem[0] == "" {
			return nil, errors.New("Invalid query parameter " + key)
		}
		switch key {
		case "expense":
			query.ExpenseId, _ = strconv.ParseUint(elem[0], 10, 64)
		case "document":
			query.DocumentId, _ = strconv.ParseUint(elem[0], 10, 64)
		default:
			return nil, errors.New("Invalid query parameter " + key)
		}
	}

	// todo: error checking
	//if ( idType == "" ) {
	//    return nil, errors.New("Missing parameters. Expecting either document= or expense=")
	//}
	return mm.Find(&query)
}

func (mm *MappingManager) Find(query *Query) ([]uint64, error) {
	return findMappings(query, mm.db)
}

func (mm *MappingManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (mm *MappingManager) Create(mapping manager.Thing) error {
	return errors.New("Not implemented")
}

func (mm *MappingManager) Update(mp manager.Thing) error {
	mapping, ok := mp.(*Mapping)
	if !ok {
		return errors.New("Non mapping passed to function")
	}
	return updateMapping(mapping, mm.db)
}

func (mm *MappingManager) NewThing() manager.Thing {
	return new(Mapping)
}

func (mm *MappingManager) Combine(one, two manager.Thing) error {
	return errors.New("Not implemented")
}

func (mm *MappingManager) Delete(cl manager.Thing) error {
	return errors.New("Not implemented")
}
