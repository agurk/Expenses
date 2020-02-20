package assets

import (
	"b2/components/managed/series"
	"b2/errors"
	"b2/manager"
	"sync"
)

// Asset represents something of value that the system knows about
// Time series are represented by individual assets linked by references
// and with different dates
type Asset struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	// Type of asset, e.g. equities, cash, etc
	Variety string `json:"type"`
	// Symbol for the type, e.g. DKK
	Symbol string `json:"symbol"`
	// External reference
	Reference    string         `json:"reference"`
	LatestSeries *series.Series `json:"latest"`
	sync.RWMutex
}

// Type returns a string representation of the asset useful when using
// manager.Thing interfaces
func (asset *Asset) Type() string {
	return "asset"
}

// GetID returns the ID of an asset
func (asset *Asset) GetID() uint64 {
	return asset.ID
}

// Merge is a synonym for Overwrite
func (asset *Asset) Merge(newThing manager.Thing) error {
	return asset.Overwrite(newThing)
}

// Overwrite is not implemented for Asset
func (asset *Asset) Overwrite(newThing manager.Thing) error {
	return errors.New("Overwrite not implemented for asset", errors.NotImplemented, "asset.Overwrite", true)
}

// Check always returns nil errors for assets
func (asset *Asset) Check() error {
	if asset.Name == "" {
		return errors.New("Name must be specified", nil, "asset.Check", true)
	}
	if asset.Variety == "" {
		return errors.New("Type must be specified", nil, "asset.Check", true)
	}
	return nil
}
