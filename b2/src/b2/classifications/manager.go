package classifications

import (
    "database/sql"
    "b2/manager"
)

type ClassificationManager struct {
    db *sql.DB
}

func (cm *ClassificationManager) Initalize (db *sql.DB) error {
    cm.db = db
    return nil
}

func (cm *ClassificationManager) GetMultiple(from, to string) ([]manager.Thing, error) {
    return getClassifications(cm.db)
}

func (cm *ClassificationManager) Get(cid uint64) (manager.Thing, error) {
    return nil, nil
}

func (cm *ClassificationManager) Save(ex manager.Thing) error {
    return nil
}

func (cm *ClassificationManager) Overwrite(ex manager.Thing) (manager.Thing, error) {
    return nil, nil
}

