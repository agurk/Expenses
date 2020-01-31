package expenses

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetMatches(e *Expense, db *sql.DB) []int64 {
	var result []int64
	res := getExactMatch(e.Description, db)
	if res > 0 {
		result = append(result, res)
	}
	result = append(result, 5)
	words := wordPower(e, db)
	for _, i := range strings.Split(strings.ToLower(e.Description), " ") {
		if i == "" {
			continue
		}
		if _, ok := words[i]; !ok {
			continue
		}
		for i, val := range words[i] {
			if val == 0 {
				continue
			}
			found := false
			for _, e := range result {
				if e == int64(i) {
					found = true
				}
			}
			if !found {
				result = append(result, int64(i))
			}
		}
	}
	return result
}

func wordPower(e *Expense, db *sql.DB) map[string]*[30]int64 {
	words := make(map[string]*[30]int64)
	rows, err := db.Query(`
		select
			description,
			c.cid
		from
			expenses e,
			classifications c,
			classificationdef cd
		where
			e.eid = c.eid
			and c.cid = cd.cid
			and c.confirmed
			and cd.validto = ""`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var desc string
		var clas int64
		rows.Scan(&desc, &clas)
		for _, i := range strings.Split(strings.ToLower(desc), " ") {
			if len(i) < 2 {
				continue
			}
			if _, ok := words[i]; !ok {
				array := [30]int64{}
				words[i] = &array
			}
			(*words[i])[clas]++
		}
	}
	return words
}

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
			and c.confirmed
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
