package accounts

import (
	"b2/errors"
	"database/sql"
)

func loadAccount(aid uint64, db *sql.DB) (*Account, error) {
	rows, err := db.Query(`
        select
            aid,
            name,
            ccy 
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
			&account.Name,
			&account.Currency)
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
								name,
								ccy)
							values ($1, $2)`,
		account.Name,
		account.Currency)

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
			ccy = $2
		where
			aid = $3`,
		account.Name,
		account.Currency,
		account.ID)
	return errors.Wrap(err, "accounts.updateAccount")
}
