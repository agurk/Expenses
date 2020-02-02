package accounts

import (
	"b2/manager"
	"sync"
)

type Account struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	sync.RWMutex
}

func (account *Account) Type() string {
	return "account"
}

func (account *Account) GetID() uint64 {
	return account.ID
}

func (account *Account) Merge(newThing manager.Thing) error {
	return account.Overwrite(newThing)
}

func (account *Account) Overwrite(newThing manager.Thing) error {
	acc, ok := newThing.(*Account)
	if !ok {
		panic("Non account passed to overwrite function")
	}
	acc.RLock()
	account.Lock()
	defer account.Unlock()
	defer acc.RUnlock()
	account.Name = acc.Name
	account.Currency = acc.Currency
	return nil
}

func (account *Account) Check() error {
	return nil
}
