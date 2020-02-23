package series

import (
	"b2/backend"
	"b2/errors"
	"b2/manager"
)

// Query represents the fields that can be used to find specific series
type Query struct {
	AssetID    uint64
	OnlyLatest bool
	Date       string
}

// SManager is a component for a Manager to control types of series in
// the system
type SManager struct {
	backend *backend.Backend
}

// Instance returns an initiated caching manager configured for series
func Instance(backend *backend.Backend) manager.Manager {
	sm := new(SManager)
	sm.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(sm)
	return general
}

// Load returns an Series
func (sm *SManager) Load(clid uint64) (manager.Thing, error) {
	return loadSeries(clid, sm.backend.DB)
}

// AfterLoad  is not implemented for Series
func (sm *SManager) AfterLoad(as manager.Thing) error {
	return nil
}

// Find returns series that match the passed in Query
func (sm *SManager) Find(params interface{}) ([]uint64, error) {
	var query *Query
	switch params.(type) {
	case *Query:
		query = params.(*Query)
	default:
		return nil, errors.New("Unexpected type passed to function", nil, "assets.Find", false)
	}
	return findSeries(query, sm.backend.DB)
}

// FindExisting will find series that match on an assetid and date and return
// the corresponding id
func (sm *SManager) FindExisting(thing manager.Thing) (uint64, error) {
	series := Cast(thing)
	return findExistingSeries(series, sm.backend.DB)
}

// Create will create a new series in the db from the passed in series
func (sm *SManager) Create(thing manager.Thing) error {
	series := Cast(thing)
	err := createSeries(series, sm.backend.DB)
	sm.backend.ReloadAssetSeries <- series.AssetID
	return err
}

// Update will update the db of the aset if its id corresponds to
// an exsiting series
func (sm *SManager) Update(thing manager.Thing) error {
	series := Cast(thing)
	err := updateSeries(series, sm.backend.DB)
	sm.backend.ReloadAssetSeries <- series.AssetID
	return err
}

// NewThing returns a new empty unsaved series
func (sm *SManager) NewThing() manager.Thing {
	return new(Series)
}

// Combine is not implemented for series
func (sm *SManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "series.Combine", true)
}

// Delete will delete the db representation for the series if there are no expenses using it
func (sm *SManager) Delete(thing manager.Thing) error {
	series := Cast(thing)
	return deleteSeries(series, sm.backend.DB)
}
