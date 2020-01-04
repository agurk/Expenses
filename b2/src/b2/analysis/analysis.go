package analysis

import (
	"b2/fxrates"
	"database/sql"
	"fmt"
	"github.com/gorilla/schema"
	"net/url"
)

type totalsParams struct {
	From            string   `schema:"from"`
	To              string   `schema:"to"`
	CCY             string   `schema:"currency"`
	Classifications []uint64 `schema:"classifications"`
}

type totalsResult struct {
	Date   string
	Totals map[uint64]float64
}

func processParams(query url.Values) (*totalsParams, error) {
	params := new(totalsParams)
	decoder := schema.NewDecoder()
	err := decoder.Decode(params, query)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func totals(params *totalsParams, fx *fxrates.FxValues, db *sql.DB) (*[]*totalsResult, error) {
	results := make(map[string]*totalsResult)
	instr := "$3"
	for i := 1; i < len(params.Classifications); i++ {
		j := i + 3
		instr = fmt.Sprintf("%s, $%d", instr, j)
	}
	args := []interface{}{}
	args = append(args, params.From)
	args = append(args, params.To)
	s := make([]interface{}, len(params.Classifications))
	for i, v := range params.Classifications {
		s[i] = v
	}
	args = append(args, s...)
	rows, err := db.Query(`
		select
			amount,
			ccy,
			date,
			c.cid
		from
			expenses e,
			classifications c,
			classificationdef cd
		where
			e.eid = c.eid
			and c.cid = cd.cid
			and date >= $1
			and date <= $2
			and c.cid in(`+instr+`)`,
		args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var date, ccy string
		var amount float64
		var cid uint64
		err = rows.Scan(&amount, &ccy, &date, &cid)
		if err != nil {
			return nil, err
		}
		// todo: better date handling
		date = date[:10]
		rate, err := fx.Get(date, params.CCY, ccy)
		if err != nil {
			return nil, err
		}
		date = date[:4]
		if _, ok := results[date]; !ok {
			results[date] = new(totalsResult)
			results[date].Totals = make(map[uint64]float64)
			results[date].Date = date
		}
		results[date].Totals[cid] += amount / rate
	}
	values := []*totalsResult{}
	for _, value := range results {
		values = append(values, value)
	}
	return &values, nil
}
