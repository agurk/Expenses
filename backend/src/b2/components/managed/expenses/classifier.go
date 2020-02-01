package expenses

import (
	"database/sql"
	"fmt"
	"strings"
)

func GetMatches(e *Expense, db *sql.DB) []int64 {
	type result struct {
		value     int64
		liklihood float64
		next      *result
	}
	resList := new(result)
	res, prob := getExactMatch(e.Description, db)
	if res > 0 {
		resList.value = res
		resList.liklihood = prob
	}
	words := wordPower(e, db)
	for _, i := range strings.Split(strings.ToLower(e.Description), " ") {
		if i == "" {
			continue
		}
		if _, ok := words[i]; !ok {
			continue
		}
		var total int64
		for _, val := range words[i] {
			total += val
		}
		for i, val := range words[i] {
			if val == 0 {
				continue
			}
			res := new(result)
			res.value = int64(i)
			res.liklihood = float64(i) / float64(total)

			var prev *result
			pos := resList
			var found bool
			for pos != nil {
				// found before inserting
				if pos.value == res.value && !found {
					break
				}
				// found old previous value
				if pos.value == res.value && found {
					// prev shouldn't be nil at this point as found cannot have been set
					prev.next = pos.next
					break
				}
				if (pos.liklihood < res.liklihood) && !found {
					found = true
					if prev == nil {
						resList = res
					} else {
						prev.next = res
					}
					res.next = pos
				}
				prev = pos
				pos = pos.next
			}
		}
	}
	var retVal []int64
	for resList != nil {
		if resList.value > 0 {
			retVal = append(retVal, resList.value)
		}
		resList = resList.next
	}
	return retVal
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
	expense.Lock()
	defer expense.Unlock()
	matches := GetMatches(expense, db)
	if len(matches) > 0 {
		expense.Metadata.Classification = matches[0]
	} else {
		expense.Metadata.Classification = getFallback(expense, db)
	}
	expense.Metadata.Confirmed = false
}

func getExactMatch(description string, db *sql.DB) (int64, float64) {
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
			ct asc`,
		description)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}
	var retVal int64
	totalCount := 0
	count := 0
	for rows.Next() {
		var classification sql.NullInt64
		err := rows.Scan(&count, &classification)
		if err != nil {
			fmt.Println(err)
			return 0, 0
		}
		if classification.Valid {
			retVal = classification.Int64
			totalCount += count
		}
	}
	return retVal, float64(count) / float64(totalCount)
}

func getFallback(e *Expense, db *sql.DB) int64 {
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
			and strftime(e.date) > date($1,'start of month','-1 month')
			and strftime(e.date) < date($2,'end of month','+1 month')
		group by
			c.cid
		order by
			ct desc`,
		e.Date, e.Date)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if rows.Next() {
		var classification sql.NullInt64
		var count int64
		err := rows.Scan(&count, &classification)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		if classification.Valid {
			return classification.Int64
		}
	}
	// todo: better fallback here if nothing found
	return 5
}
