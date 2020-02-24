package classifications

import (
	"b2/backend"
	"b2/errors"
	"b2/manager"
	"net/url"

	"github.com/gorilla/schema"
)

// ClassificationManager implements the Component interface to manager Classifications
// generally this would be expected to be run through a simple (non-caching) manager as
// the number of classifications is small
type ClassificationManager struct {
	backend *backend.Backend
}

// Query holds the fields that can be used to find a classification. Empty fields will be
// ignored. The schema fields are parsed from a url query string
//
// Specifying a from and to will give classifications that are valid for
// that entire period, and any valid for only part of that period will
// be excluded
type Query struct {
	From string `schema:"from"`
	To   string `schema:"to"`
	Date string `schema:"date"`
}

// Instance returns a configured manager configured to handle classifications. This
// is managed by a simple manager as there are few classifications
func Instance(backend *backend.Backend) manager.Manager {
	cm := new(ClassificationManager)
	cm.backend = backend
	general := new(manager.SimpleManager)
	general.Initalize(cm)
	return general
}

// Load returns a classification
func (cm *ClassificationManager) Load(clid uint64) (manager.Thing, error) {
	return loadClassification(clid, cm.backend.DB)
}

// AfterLoad performs no fuction for classifications
func (cm *ClassificationManager) AfterLoad(classification manager.Thing) error {
	return nil
}

// Find returns a list of classifications based on the possible search criteria
// specified in a Query or by parsing a url.Values struct
func (cm *ClassificationManager) Find(params interface{}) ([]uint64, error) {
	var search *Query
	switch params.(type) {
	case *Query:
		search = params.(*Query)
	case url.Values:
		search = new(Query)
		decoder := schema.NewDecoder()
		err := decoder.Decode(search, params.(url.Values))
		if err != nil {
			return nil, errors.Wrap(err, "classifications.Find")
		}
	default:
		return nil, errors.New("Unexpected type passed to function", nil, "classifications.Find", false)
	}
	return findClassifications(search, cm.backend.DB)
}

// FindExisting is not implemented for classifications
func (cm *ClassificationManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create saves a new instance of a classification passed to the function. The ID
// of the newly created object will be populated into the classification
func (cm *ClassificationManager) Create(cl manager.Thing) error {
	classification := Cast(cl)
	return createClassification(classification, cm.backend.DB)
}

// Update will save any changes to the classification passed to the function. This
// will not create a new classification if the ID is 0 or set to a non-existing classification
func (cm *ClassificationManager) Update(cl manager.Thing) error {
	classification := Cast(cl)
	return updateClassification(classification, cm.backend.DB)
}

// NewThing returns a newly created empty classification
func (cm *ClassificationManager) NewThing() manager.Thing {
	return new(Classification)
}

// Combine is not implemented for classifications
func (cm *ClassificationManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "classifications.Combine", true)
}

// Delete the classification if there are no expenses using it
func (cm *ClassificationManager) Delete(cl manager.Thing) error {
	classification := Cast(cl)
	return deleteClassification(classification, cm.backend.DB)
}
