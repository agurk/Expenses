package assets

import (
	"b2/backend"
	"b2/components/managed/series"
	"b2/errors"
	"b2/manager"
	"net/url"

	"github.com/gorilla/schema"
)

// AssetManager is a component for a Manager to control types of assets in
// the system
type AssetManager struct {
	backend *backend.Backend
}

// Query holds the values used for finding the asset
type Query struct {
	Date string `schema:"date"`
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
func (am *AssetManager) AfterLoad(thing manager.Thing) error {
	asset := Cast(thing)
	v := new(series.Query)
	v.AssetID = asset.ID
	v.OnlyLatest = true
	srs, err := am.backend.Series.Find(v)
	if err != nil {
		return errors.Wrap(err, "assets.AfterLoad")
	}

	if len(srs) > 1 {
		return errors.New("Multiple series found", nil, "assets.AfterLoad", false)
	}

	asset.Lock()
	defer asset.Unlock()

	asset.LatestSeries = nil

	if len(srs) == 1 {
		asset.LatestSeries = series.Cast(srs[0])
	}
	return nil
}

// Find returns all assets
func (am *AssetManager) Find(params interface{}) ([]uint64, error) {
	var query *Query
	switch params.(type) {
	case *Query:
		query = params.(*Query)
	case url.Values:
		query = new(Query)
		decoder := schema.NewDecoder()
		err := decoder.Decode(query, params.(url.Values))
		if err != nil {
			return nil, errors.Wrap(err, "assets.Find")
		}
	default:
		return nil, errors.New("Unexpected type passed to function", nil, "assets.Find", false)
	}
	return findAssets(query, am.backend.DB)
}

// FindExisting does nothing
func (am *AssetManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create will create a new asset in the db from the passed in asset
func (am *AssetManager) Create(thing manager.Thing) error {
	asset := Cast(thing)
	return createAsset(asset, am.backend.DB)
}

// Update will update the db of the aset if its id corresponds to
// an exsiting asset
func (am *AssetManager) Update(thing manager.Thing) error {
	asset := Cast(thing)
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
func (am *AssetManager) Delete(thing manager.Thing) error {
	asset := Cast(thing)
	return deleteAsset(asset, am.backend.DB)
}

// Process is not implemented for Assets
func (am *AssetManager) Process(id uint64) {
}
