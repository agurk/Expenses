package docexmappings

import (
	"b2/backend"
	"b2/components/changes"
	"b2/errors"
	"b2/manager"
)

// Query is passed to the Find function to find mappings that match the
// criteria in the struct
type Query struct {
	ExpenseID  uint64
	DocumentID uint64
}

// MappingManager is a component manager for mappings between documents and expenses
type MappingManager struct {
	backend *backend.Backend
}

// Instance retutrns a caching manager configured to manage document/expense mappings
func Instance(backend *backend.Backend) manager.Manager {
	mm := new(MappingManager)
	mm.backend = backend
	general := new(manager.CachingManager)
	general.Initalize(mm)
	return general
}

// Load returns the mapping matching the mapping ID argument
func (mm *MappingManager) Load(dmid uint64) (manager.Thing, error) {
	return loadMapping(dmid, mm.backend.DB)
}

// AfterLoad does nothing for mappings
func (mm *MappingManager) AfterLoad(mapping manager.Thing) error {
	return nil
}

// Find returns a slice of IDs of mappings that match the criteria in a Query
func (mm *MappingManager) Find(query interface{}) ([]uint64, error) {
	var search *Query
	switch query.(type) {
	case *Query:
		search = query.(*Query)
	default:
		panic("Unknown type passed to find function")
	}

	return findMappings(search, mm.backend.DB)
}

// FindExisting is not implemented for mappings
func (mm *MappingManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create saves a new mapping from the one passed in
func (mm *MappingManager) Create(thing manager.Thing) error {
	mapping := Cast(thing)
	_, err := mm.backend.Expenses.Get(mapping.EID)
	if err != nil {
		return errors.Wrap(err, "mapping.Create")
	}
	_, err = mm.backend.Documents.Get(mapping.DID)
	if err != nil {
		return errors.Wrap(err, "mapping.Create")
	}
	err = createMapping(mapping, mm.backend.DB)
	if err != nil {
		return errors.Wrap(err, "mapping.Create")
	}
	mm.backend.ReloadExpenseMappings <- mapping.EID
	mm.backend.ReloadDocumentMappings <- mapping.DID
	mm.backend.Change <- changes.ExpenseEvent
	mm.backend.Change <- changes.DocumentEvent

	return errors.Wrap(err, "mapping.Create")
}

// Update updates the db of any changes made to the mapping
func (mm *MappingManager) Update(thing manager.Thing) error {
	mapping := Cast(thing)
	return updateMapping(mapping, mm.backend.DB)
}

// NewThing returns a new instatiated unsaved empty mapping
func (mm *MappingManager) NewThing() manager.Thing {
	return new(Mapping)
}

// Combine is not implemented for mappings
func (mm *MappingManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "mapping.Combine", true)
}

// Delete removes a mapping from the db and alerts any expenses and documents that this
// has occured by requesting their mappings to be updated
func (mm *MappingManager) Delete(thing manager.Thing) error {
	mapping := Cast(thing)
	err := deleteMapping(mapping, mm.backend.DB)
	if err != nil {
		return nil
	}
	// alert the dependecies that their mappings have changed
	mm.backend.ReloadExpenseMappings <- mapping.EID
	mm.backend.ReloadDocumentMappings <- mapping.DID
	mapping.deleted = true
	mm.backend.Change <- changes.ExpenseEvent
	mm.backend.Change <- changes.DocumentEvent
	return nil
}
