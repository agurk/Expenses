package accounts

import (
	"database/sql"
	"errors"
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
		return nil, err
	}
	defer rows.Close()
	account := new(Account)
	if rows.Next() {
		err = rows.Scan(&account.ID,
			&account.Name,
			&account.Currency)
	} else {
		return nil, errors.New("404")
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func findAccounts(db *sql.DB) ([]uint64, error) {
	rows, err := db.Query("select aid from accountdef")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var aids []uint64
	for rows.Next() {
		var aid uint64
		err = rows.Scan(&aid)
		if err != nil {
			return nil, err
		}
		aids = append(aids, aid)
	}
	return aids, err
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
		return err
	}
	rid, err := res.LastInsertId()
	if err == nil && rid > 0 {
		account.ID = uint64(rid)
	} else {
		return errors.New("Error creating new account")
	}
	return nil
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
	return err
}
