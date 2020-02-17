package accounts

import (
	"b2/errors"
	"database/sql"
	"fmt"
)

func loadAccount(aid uint64, db *sql.DB) (*Account, error) {
	rows, err := db.Query(`
        select
            aid,
            name,
        from
            accountdef
        where
            aid = $1`,
		aid)
	if err != nil {
		return nil, errors.Wrap(err, "accounts.loadAccount")
	}
	defer rows.Close()
	account := new(Account)
	if rows.Next() {
		err = rows.Scan(&account.ID,
			&account.Name)
	} else {
		return nil, errors.New("Account not found", errors.ThingNotFound, "accounts.loadAccount", true)
	}
	if err != nil {
		return nil, errors.Wrap(err, "accounts.loadAccount")
	}
	return account, nil
}

func findAccounts(db *sql.DB) ([]uint64, error) {
	rows, err := db.Query("select aid from accountdef")
	if err != nil {
		return nil, errors.Wrap(err, "accounts.findAccounts")
	}
	defer rows.Close()
	var aids []uint64
	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid)
		if err != nil {
			return nil, errors.Wrap(err, "accounts.findAccounts")
		}
		aids = append(aids, aid)
	}
	return aids, errors.Wrap(err, "accounts.findAccounts")
}

func createAccount(account *Account, db *sql.DB) error {
	account.Lock()
	defer account.Unlock()
	res, err := db.Exec(`insert into
							accountdef (
								name)
							values ($1)`,
		account.Name)

	if err != nil {
		return errors.Wrap(err, "accounts.createAccount")
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		account.ID = uint64(rid)
	} else {
		return errors.New("Error creating new account", errors.InternalError, "accounts.createAccount", false)
	}
	return errors.Wrap(err, "accounts.createAccount")
}

func updateAccount(account *Account, db *sql.DB) error {
	account.RLock()
	defer account.RUnlock()
	_, err := db.Exec(`
		update
			accountdef
		set
			name = $1,
		where
			aid = $2`,
		account.Name,
		account.ID)
	return errors.Wrap(err, "accounts.updateAccount")
}

func deleteAccount(account *Account, db *sql.DB) error {
	rows, err := db.Query("select count(*) from expenses where aid = $1", account.ID)
	if err != nil {
		return errors.Wrap(err, "account.deleteAccount (count)")
	}
	defer rows.Close()
	for rows.Next() {
		var count uint64
		err = rows.Scan(&count)
		if err != nil {
			return errors.Wrap(err, "account.deleteAccount(count rows.Scan)")
		}
		if count > 0 {
			return errors.New(fmt.Sprintf("Cannot delete account as it's being used by %d expenses", count),
				nil, "account.deleteAccount", true)
		}
	}
	_, err = db.Exec(`
        delete from
			accountdef
        where
            aid = $1`,
		account.ID)
	return errors.Wrap(err, "account.deleteAccount(delete)")

}
