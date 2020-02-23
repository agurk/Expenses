package analysis

import (
	"b2/errors"
	"b2/moneyutils"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/gorilla/schema"
)

type totalsParams struct {
	From            string   `schema:"from"`
	To              string   `schema:"to"`
	CCY             string   `schema:"currency"`
	Classifications []uint64 `schema:"classifications"`
	AllSpend        bool     `schema:"allSpend"`
	Grouping        string   `schema:"grouping"`
}

type totalsResult struct {
	Classifications map[uint64]float64 `json:"classifications"`
	AllSpend        float64            `json:"allSpend"`
}

func processParams(query url.Values) (*totalsParams, error) {
	params := new(totalsParams)
	decoder := schema.NewDecoder()
	err := decoder.Decode(params, query)
	if err != nil {
		return nil, errors.Wrap(err, "analysis.processParams")
	}
	return params, nil
}

func processRow(rows *sql.Rows, params *totalsParams, results *map[string]*totalsResult, fx *moneyutils.FxValues, rowType string) error {
	for rows.Next() {
		var date, ccy string
		var amount int64
		var cid uint64
		err := rows.Scan(&amount, &ccy, &date, &cid)
		if err != nil {
			return errors.Wrap(err, "analysis.processRow")
		}
		// todo: better date handling
		date = date[:10]
		rate, err := fx.Rate(date, params.CCY, ccy)
		if err != nil {
			return errors.Wrap(err, "analysis.processRow")
		}
		key := date[:4]
		switch params.Grouping {
		case "together":
			key = "total"
		}
		if _, ok := (*results)[key]; !ok {
			(*results)[key] = new(totalsResult)
			(*results)[key].Classifications = make(map[uint64]float64)
		}
		ccyAmt, err := moneyutils.CurrencyAmount(amount, ccy)
		if err != nil {
			return errors.Wrap(err, "analysis.processRow")
		}
		switch rowType {
		case "all":
			(*results)[key].AllSpend += ccyAmt / rate
		case "classifications":
			(*results)[key].Classifications[cid] += ccyAmt / rate
		}
	}
	return nil
}

func analyseAllSpend(params *totalsParams, results *map[string]*totalsResult, fx *moneyutils.FxValues, db *sql.DB) error {
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
			and cd.isExpense`,
		params.From, params.To)
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "analysis.analyseAllSpend")
	}
	return processRow(rows, params, results, fx, "all")
}

func analyseClassifications(params *totalsParams, results *map[string]*totalsResult, fx *moneyutils.FxValues, db *sql.DB) error {
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
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "analysis.analyseClassifications")
	}
	return processRow(rows, params, results, fx, "classifications")
}

func totals(params *totalsParams, fx *moneyutils.FxValues, db *sql.DB) (*map[string]*totalsResult, error) {
	results := make(map[string]*totalsResult)
	err := analyseClassifications(params, &results, fx, db)
	if err != nil {
		return nil, errors.Wrap(err, "analysis.totals")
	}
	err = analyseAllSpend(params, &results, fx, db)
	if err != nil {
		return nil, errors.Wrap(err, "analysis.totals")
	}
	return &results, nil
}
