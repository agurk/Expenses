package analysis

import (
	"b2/errors"
	"b2/moneyutils"
	"database/sql"
)

type assetsResult struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

func assets(rates *moneyutils.FxValues, db *sql.DB) ([]*assetsResult, error) {
	results := []*assetsResult{}
	//results := new([]*assetsResult)
	rows, err := db.Query(`
		select
			a.name,
			a.type,
			a.symbol,
			s.date,
			s.amountwhole,
			s.amountfractional
		from
			assets a
		inner join
			(select
				max(date) date,
				amountwhole,
				amountfractional,
				asid
			from
				assetseries
			group by
				asid
			) s on a.asid = s.asid
		`)
	if err != nil {
		return nil, errors.Wrap(err, "analysis.assets")
	}
	for rows.Next() {
		var name, variety, symbol, date string
		var amount, fraction int64
		err = rows.Scan(&name,
			&variety,
			&symbol,
			&date,
			&amount,
			&fraction)
		if err != nil {
			return nil, errors.Wrap(err, "analysis.assets")
		}
		results = append(results, makeResult(name, variety, symbol, date, amount, rates, db))
	}
	return results, nil
}

func equityQuote(symbol string, db *sql.DB) (float64, string, error) {
	rows, err := db.Query(`
		select
			price,
			currency
		from
			_Quotes
		where
			ticker = $1
		order by
			date desc
		limit
			1`,
		symbol)
	if err != nil {
		return 0, "", errors.Wrap(err, "assetAnalysis.equityQuote")
	}
	var price float64
	var currency string
	if rows.Next() {
		err = rows.Scan(&price, &currency)
		if err != nil {
			return 0, "", errors.Wrap(err, "analysis.equityQuote")
		}

	} else {
		return 0, "", errors.New("Quote not found for "+symbol, nil, "analysis.equityQuote", true)
	}
	return price, currency, nil
}

func makeResult(name, variety, symbol, date string, amount int64, rates *moneyutils.FxValues, db *sql.DB) *assetsResult {
	result := new(assetsResult)
	result.Name = name
	switch variety {
	case "cash":
		rate, err := rates.Rate(date, "GBP", symbol)
		if err != nil {
			errors.Print(err)
		}
		result.Amount = float64(amount) / rate
	case "equity":
		price, currency, err := equityQuote(symbol, db)
		if err != nil {
			errors.Print(err)
			break
		}
		rate, err := rates.Rate(date, "GBP", currency)
		if err != nil {
			errors.Print(err)
			break
		}
		result.Amount = float64(amount) * price / rate
	}
	return result

}
