package mappings

import (
    "database/sql"
    "net/url"
    "b2/manager"
    "errors"
    "strconv"
)

type MappingManager struct {
    db *sql.DB
}

func (mm *MappingManager) Initalize (db *sql.DB) {
    mm.db = db
}

func (mm *MappingManager) Load(dmid uint64) (manager.Thing, error) {
    return loadMapping(dmid, mm.db)
}

func (mm *MappingManager) AfterLoad(mapping manager.Thing) (error) {
    return nil
}

func (mm *MappingManager) Find(params url.Values) ([]uint64, error) {
    var id uint64
    var idType string
    for key, elem := range params {
        // Query() returns empty string as value when no value set for key
        if (len(elem) != 1 || elem[0] == "" ) {
            return nil, errors.New("Invalid query parameter " + key)
        }
        switch key {
        case "expense":
            id, _ = strconv.ParseUint(elem[0], 10, 64)
            idType = "expense"
        case "document":
            id, _ = strconv.ParseUint(elem[0], 10, 64)
            idType = "document"
        default:
            return nil, errors.New("Invalid query parameter " + key)
        }
    }

    if ( idType == "" ) {
        return nil, errors.New("Missing parameters. Expecting either document= or expense=")
    }

    return findMappings(idType, id, mm.db)
}

func (mm *MappingManager) Create(mapping manager.Thing) error {
    return errors.New("Not implemented")
}

func (mm *MappingManager) Update(mapping manager.Thing) error {
    return errors.New("Not implemented")
}

func (mm *MappingManager) Merge(from manager.Thing, to manager.Thing) error {
    return errors.New("Not implemented")
}

func (mm *MappingManager) NewThing() manager.Thing {
    return new(Mapping)
}

