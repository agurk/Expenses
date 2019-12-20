package classifications

import (
    "database/sql"
    "net/url"
    "b2/manager"
    "errors"
)

type ClassificationManager struct {
    db *sql.DB
}

func Instance(db *sql.DB) *manager.Manager {
    cm := new (ClassificationManager)
    cm.initalize(db)
    general := new (manager.Manager)
    general.Initalize(cm)
    return general
}

func (cm *ClassificationManager) initalize (db *sql.DB) {
    cm.db = db
}

func (cm *ClassificationManager) Load(clid uint64) (manager.Thing, error) {
    return loadClassification(clid, cm.db)
}

func (cm *ClassificationManager) AfterLoad(classification manager.Thing) (error) {
    return nil
}

func (cm *ClassificationManager) FindFromUrl(params url.Values) ([]uint64, error) {
    return findClassifications(cm.db)
}

func (cm *ClassificationManager) FindExisting(thing manager.Thing) (uint64, error) {
    return 0, nil
}

func (cm *ClassificationManager) Create(cl manager.Thing) error {
    return errors.New("Not implemented")
}

func (cm *ClassificationManager) Update(cl manager.Thing) error {
    return errors.New("Not implemented")
}

func (cm *ClassificationManager) NewThing() manager.Thing {
    return new(Classification)
}

