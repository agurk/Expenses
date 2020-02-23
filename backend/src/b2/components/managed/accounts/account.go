package accounts

import (
	"b2/manager"
	"sync"
)

// Account represents an account
type Account struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	sync.RWMutex
}

// Cast a thing into an *Account or panic
func Cast(thing manager.Thing) *Account {
	account, ok := thing.(*Account)
	if !ok {
		panic("Non account passed to overwrite function")
	}
	return account
}

// Type returns a string representation of the account useful when using
// manager.Thing interfaces
func (account *Account) Type() string {
	return "account"
}

// GetID returns the ID of an account
func (account *Account) GetID() uint64 {
	return account.ID
}

// Merge is a synonym for Overwrite
func (account *Account) Merge(newThing manager.Thing) error {
	return account.Overwrite(newThing)
}

// Overwrite replaces the Name of the exsiting account with
// the values in the account passed in
func (account *Account) Overwrite(newThing manager.Thing) error {
	acc := Cast(newThing)
	acc.RLock()
	account.Lock()
	defer account.Unlock()
	defer acc.RUnlock()
	account.Name = acc.Name
	return nil
}

// Check always returns nil errors for accounts
func (account *Account) Check() error {
	return nil
}
