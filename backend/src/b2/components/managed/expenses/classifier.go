package expenses

import (
	"database/sql"
	"fmt"
)

func classifyExpense(expense *Expense, db *sql.DB) {
	// todo: add some better logic here
	expense.Lock()
	defer expense.Unlock()
	exact := getExactMatch(expense.Description, db)
	if exact > 0 {
		expense.Metadata.Classification = exact
	} else {
		expense.Metadata.Classification = 5
	}
	expense.Metadata.Confirmed = false
}

func getExactMatch(description string, db *sql.DB) int64 {
	rows, err := db.Query(`
		select
			count(*) ct,
			c.cid
		from
			expenses e,
			classifications c
		where
			e.eid = c.eid
			and e.confirmed
			and e.description = $1
		group by
			c.cid
		order by
			ct desc`,
		description)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if rows.Next() {
		var count uint64
		var classification sql.NullInt64
		err := rows.Scan(&count, &classification)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		if classification.Valid {
			return classification.Int64
		}
	}
	return 0
}
