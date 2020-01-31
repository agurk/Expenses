package accounts

import (
	"b2/backend"
	"b2/manager"
	"errors"
)

type AccountManager struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) manager.Manager {
	am := new(AccountManager)
	am.backend = backend
	general := new(manager.SimpleManager)
	general.Initalize(am)
	return general
}

func (am *AccountManager) Load(clid uint64) (manager.Thing, error) {
	return loadAccount(clid, am.backend.DB)
}

func (am *AccountManager) AfterLoad(account manager.Thing) error {
	return nil
}

func (am *AccountManager) Find(params interface{}) ([]uint64, error) {
	return findAccounts(am.backend.DB)
}

func (am *AccountManager) FindExisting(thing manager.Thing) (uint64, error) {
	return 0, nil
}

func (am *AccountManager) Create(cl manager.Thing) error {
	account, ok := cl.(*Account)
	if !ok {
		return errors.New("Non account passed to function")
	}
	return createAccount(account, am.backend.DB)
}

func (am *AccountManager) Update(cl manager.Thing) error {
	account, ok := cl.(*Account)
	if !ok {
		return errors.New("Non account passed to function")
	}
	return updateAccount(account, am.backend.DB)
}

func (am *AccountManager) NewThing() manager.Thing {
	return new(Account)
}

func (am *AccountManager) Combine(one, two manager.Thing, params string) error {
	return errors.New("Not implemented")
}

func (am *AccountManager) Delete(cl manager.Thing) error {
	return errors.New("Not implemented")
}

func (am *AccountManager) Process(id uint64) {
}
