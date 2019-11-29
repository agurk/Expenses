package documents

import (
    "net/url"
    "database/sql"
    "errors"
    "b2/manager"
)

type DocManager struct {
    db *sql.DB
}

func (dm *DocManager) Initalize (db *sql.DB) {
    dm.db = db
}

func (dm *DocManager) Load(did uint64) (manager.Thing, error) {
    return loadDocument(did, dm.db)

}

func (dm *DocManager) Find(params url.Values) ([]uint64, error) {
    return findDocuments(dm.db)
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

func (dm *DocManager) Merge(from manager.Thing, to manager.Thing) error {
    document, ok := from.(*Document)
    if !ok {
        return errors.New("Non document passed to function")
    }
    oldDoc, ok := to.(*Document)
    if !ok {
        return errors.New("Non document passed to function")
    }
    document.RLock()
    oldDoc.Lock()
    oldDoc.Filename = document.Filename
    oldDoc.Deleted = document.Deleted
    oldDoc.Date = document.Date
    oldDoc.Text = document.Text
    document.RUnlock()
    oldDoc.Unlock()
    return nil
}

func (dm *DocManager) NewThing() manager.Thing {
    return new(Document)
}

