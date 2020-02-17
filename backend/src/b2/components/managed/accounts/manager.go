package accounts

import (
	"b2/backend"
	"b2/errors"
	"b2/manager"
)

// AccountManager is a managed component designed to managed the accounts used
// in this application
type AccountManager struct {
	backend *backend.Backend
}

// Instance returns an initiated simple manager configured for accounts
func Instance(backend *backend.Backend) manager.Manager {
	am := new(AccountManager)
	am.backend = backend
	general := new(manager.SimpleManager)
	general.Initalize(am)
	return general
}

// Load returns an Account
func (am *AccountManager) Load(clid uint64) (manager.Thing, error) {
	return loadAccount(clid, am.backend.DB)
}

// AfterLoad does nothing
func (am *AccountManager) AfterLoad(account manager.Thing) error {
	return nil
}

// Find returns all accounts
func (am *AccountManager) Find(params interface{}) ([]uint64, error) {
	return findAccounts(am.backend.DB)
}

// FindExisting does nothing
func (am *AccountManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

// Create will create a new account in the db from the passed in account
func (am *AccountManager) Create(cl manager.Thing) error {
	account, ok := cl.(*Account)
	if !ok {
		panic("Non account passed to function")
	}
	return createAccount(account, am.backend.DB)
}

// Update will update the db of the account if its id corresponds to
// an exsiting account
func (am *AccountManager) Update(cl manager.Thing) error {
	account, ok := cl.(*Account)
	if !ok {
		panic("Non account passed to function")
	}
	return updateAccount(account, am.backend.DB)
}

// NewThing returns a new empty unsaved account
func (am *AccountManager) NewThing() manager.Thing {
	return new(Account)
}

// Combine is not implemented for accounts
func (am *AccountManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented", errors.NotImplemented, "accounts.Combine", true)
}

// Delete will delete the db representation for the account if there are no expenses using it
func (am *AccountManager) Delete(cl manager.Thing) error {
	account, ok := cl.(*Account)
	if !ok {
		panic("Non account passed to function")
	}
	return deleteAccount(account, am.backend.DB)
}
