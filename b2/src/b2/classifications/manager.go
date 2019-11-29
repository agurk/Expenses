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

func (cm *ClassificationManager) Initalize (db *sql.DB) {
    cm.db = db
}

func (cm *ClassificationManager) Load(clid uint64) (manager.Thing, error) {
    return loadClassification(clid, cm.db)
}

func (cm *ClassificationManager) Find(params url.Values) ([]uint64, error) {
    return findClassifications(cm.db)
}

func (cm *ClassificationManager) Create(cl manager.Thing) error {
    return errors.New("Not implemented")
}

func (cm *ClassificationManager) Update(cl manager.Thing) error {
    return errors.New("Not implemented")
}

func (cm *ClassificationManager) Merge(from manager.Thing, to manager.Thing) error {
    return errors.New("Not implemented")
}

func (cm *ClassificationManager) NewThing() manager.Thing {
    return new(Classification)
}

