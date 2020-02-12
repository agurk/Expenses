package classifications

import (
	"b2/backend"
	"b2/errors"
	"b2/manager"
)

type ClassificationManager struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) manager.Manager {
	cm := new(ClassificationManager)
	cm.backend = backend
	general := new(manager.SimpleManager)
	general.Initalize(cm)
	return general
}

func (cm *ClassificationManager) Load(clid uint64) (manager.Thing, error) {
	return loadClassification(clid, cm.backend.DB)
}

func (cm *ClassificationManager) AfterLoad(classification manager.Thing) error {
	return nil
}

func (cm *ClassificationManager) Find(params interface{}) ([]uint64, error) {
	return findClassifications(cm.backend.DB)
}

func (cm *ClassificationManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (cm *ClassificationManager) Create(cl manager.Thing) error {
	classification, ok := cl.(*Classification)
	if !ok {
		panic("Non classification passed to function")
	}
	return createClassification(classification, cm.backend.DB)
}

func (cm *ClassificationManager) Update(cl manager.Thing) error {
	classification, ok := cl.(*Classification)
	if !ok {
		panic("Non classification passed to function")
	}
	return updateClassification(classification, cm.backend.DB)
}

func (cm *ClassificationManager) NewThing() manager.Thing {
	return new(Classification)
}

func (cm *ClassificationManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "classifications.Combine", true)
}

func (cm *ClassificationManager) Delete(cl manager.Thing) error {
	return errors.New("Not implemented", errors.NotImplemented, "classifications.Delete", true)
}

func (cm *ClassificationManager) Process(id uint64) {
}
