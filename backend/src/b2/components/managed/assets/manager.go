package assets

import (
	"b2/backend"
	"b2/components/managed/series"
	"b2/errors"
	"b2/manager"
)

// AssetManager is a component for a Manager to control types of assets in
// the system
type AssetManager struct {
	backend *backend.Backend
}

// Instance returns an initiated caching manager configured for assets
func Instance(backend *backend.Backend) manager.Manager {
	am := new(AssetManager)
	am.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(am)
	return general
}

// Load returns an Asset
func (am *AssetManager) Load(clid uint64) (manager.Thing, error) {
	return loadAsset(clid, am.backend.DB)
}

// AfterLoad loads the time series data into the asset
func (am *AssetManager) AfterLoad(as manager.Thing) error {
	asset, ok := as.(*Asset)
	if !ok {
		panic("Non asset passed to function")
	}
	v := new(series.Query)
	v.AssetID = asset.ID
	mapps, err := am.backend.Series.Find(v)
	asset.Lock()
	defer asset.Unlock()
	asset.Series = []*series.Series{}
	for _, thing := range mapps {
		mapping, ok := thing.(*(series.Series))
		if !ok {
			panic("Non series returned from function")
		}
		asset.Series = append(asset.Series, mapping)
	}
	return errors.Wrap(err, "asset.AfterLoad")
}

// Find returns all assets
func (am *AssetManager) Find(params interface{}) ([]uint64, error) {
	return findAssets(am.backend.DB)
}

// FindExisting does nothing
func (am *AssetManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create will create a new asset in the db from the passed in asset
func (am *AssetManager) Create(cl manager.Thing) error {
	asset, ok := cl.(*Asset)
	if !ok {
		panic("Non asset passed to function")
	}
	return createAsset(asset, am.backend.DB)
}

// Update will update the db of the aset if its id corresponds to
// an exsiting asset
func (am *AssetManager) Update(cl manager.Thing) error {
	asset, ok := cl.(*Asset)
	if !ok {
		panic("Non asset passed to function")
	}
	return updateAsset(asset, am.backend.DB)
}

// NewThing returns a new empty unsaved asset
func (am *AssetManager) NewThing() manager.Thing {
	return new(Asset)
}

// Combine is not implemented for assets
func (am *AssetManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "assets.Combine", true)
}

// Delete will delete the db representation for the asset if there are no expenses using it
func (am *AssetManager) Delete(cl manager.Thing) error {
	asset, ok := cl.(*Asset)
	if !ok {
		panic("Non asset passed to function")
	}
	return deleteAsset(asset, am.backend.DB)
}

// Process is not implemented for Assets
func (am *AssetManager) Process(id uint64) {
}
