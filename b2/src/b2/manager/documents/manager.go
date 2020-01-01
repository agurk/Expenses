package documents

import (
	"b2/manager"
	"b2/manager/docexmappings"
	"database/sql"
	"errors"
	"net/url"
	"strconv"
)

type DocManager struct {
	db *sql.DB
	mm manager.Manager
}

func Instance(db *sql.DB, mm manager.Manager) manager.Manager {
	dm := new(DocManager)
	dm.initalize(db, mm)
	general := new(manager.CachingManager)
	general.Initalize(dm)
	return general
}

func (dm *DocManager) initalize(db *sql.DB, mm manager.Manager) {
	dm.db = db
	dm.mm = mm
}

func (dm *DocManager) Load(did uint64) (manager.Thing, error) {
	return loadDocument(did, dm.db)

}

func (dm *DocManager) AfterLoad(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		return errors.New("Non document passed to function")
	}
	v := url.Values{}
	v.Set("document", strconv.FormatUint(document.ID, 10))
	mapps, err := dm.mm.GetMultiple(v)
	for _, thing := range mapps {
		mapping, ok := thing.(*(docexmappings.Mapping))
		if !ok {
			return errors.New("Non mapping returned from function")
		}
		document.Expenses = append(document.Expenses, mapping)
	}
	return err
}

func (dm *DocManager) FindFromUrl(params url.Values) ([]uint64, error) {
	return findDocuments(dm.db)
}

func (dm *DocManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (dm *DocManager) Create(doc manager.Thing) error {
	document, ok := doc.(*Document)
	if !ok {
		return errors.New("Non document passed to function")
	}
	return createDocument(document, dm.db)
}

func (dm *DocManager) Update(doc manager.Thing) error {
	return errors.New("Method not implemented")
}

func (dm *DocManager) NewThing() manager.Thing {
	return new(Document)
}

func (dm *DocManager) Combine(one, two manager.Thing) error {
	return errors.New("Not implemented")
}

func (dm *DocManager) Delete(doc manager.Thing) error {
	return errors.New("Not implemented")
}
