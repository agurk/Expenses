package analysis

import (
	"b2/backend"
	"b2/components/managed/assets"
	"b2/components/managed/series"
	"b2/errors"
	"b2/moneyutils"
	"database/sql"
	"time"
)

type assetsResult struct {
	Name   string           `json:"name"`
	Values []*resultDetails `json:"values"`
}

type resultDetails struct {
	Amount  float64            `json:"amount"`
	Details map[string]float64 `json:"details"`
	Date    string             `json:"date"`
}

func priceAsset(asset *assets.Asset, ccy, date string, rates *moneyutils.FxValues, backend *backend.Backend) (float64, error) {
	query := new(series.Query)
	query.Date = date
	query.AssetID = asset.ID
	query.OnlyLatest = true
	s, err := backend.Series.Find(query)
	if err != nil {
		return 0, errors.Wrap(err, "analysis.priceAsset")
	}
	if len(s) != 1 {
		return 0, nil //errors.New("Wrong number of series returned", nil, "analysis.priceAsset", false)
	}
	se, ok := s[0].(*series.Series)
	if !ok {
		panic("Non series returned from function")
	}
	switch asset.Variety {
	case "cash":
		rate, err := rates.Rate(date, ccy, asset.Symbol)
		if err != nil {
			return 0, errors.Wrap(err, "analysis.priceAsset")
		}
		return se.AmountFloat() / rate, nil
	case "equity":
		// todo include date here
		price, currency, err := equityQuote(asset.Symbol, date, backend.DB)
		if err != nil {
			return 0, errors.Wrap(err, "analysis.priceAsset")
		}
		rate, err := rates.Rate(date, ccy, currency)
		if err != nil {
			return 0, errors.Wrap(err, "analysis.priceAsset")
		}
		return se.AmountFloat() * price / rate, nil
	}
	return 0, nil
}

func asts(p *totalsParams, rates *moneyutils.FxValues, backend *backend.Backend) ([]*assetsResult, error) {
	today := time.Now().Format("2006-01-02")
	week := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	month := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	year := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	dates := []string{today, week, month, year}
	results := []*assetsResult{}
	// todo don't ignore errors
	ass, _ := backend.Assets.Find(new(assets.Query))
	for _, thing := range ass {
		asset, _ := thing.(*(assets.Asset))
		for _, date := range dates {
			val, err := priceAsset(asset, p.CCY, date, rates, backend)
			if err != nil {
				return nil, errors.Wrap(err, "analyse.asts")
			}
			results = addResult(asset.Name, asset.Variety, asset.Symbol, date, val, results)
		}
	}
	return results, nil
}

func addResult(name, variety, symbol, date string, amount float64, results []*assetsResult) []*assetsResult {
	term := symbol + " " + variety
	for _, res := range results {
		if res.Name == term {
			for _, details := range res.Values {
				if details.Date == date {
					details.Amount += amount
					details.Details[name] = amount
					return results
				}
			}
			deets := new(resultDetails)
			deets.Amount = amount
			deets.Details = make(map[string]float64)
			deets.Details[name] = amount
			deets.Date = date
			res.Values = append(res.Values, deets)
			return results
		}
	}
	result := new(assetsResult)
	result.Name = term
	results = append(results, result)

	deets := new(resultDetails)
	deets.Amount = amount
	deets.Details = make(map[string]float64)
	deets.Details[name] = amount
	deets.Date = date
	result.Values = append(result.Values, deets)
	return results
}

func equityQuote(symbol, date string, db *sql.DB) (float64, string, error) {
	rows, err := db.Query(`
		select
			price,
			currency
		from
			_Quotes
		where
			ticker = $1
			and date <= $2
		order by
			date desc
		limit
			1`,
		symbol, date)
	defer rows.Close()
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
		return 0, "", errors.New("Quote not found for "+symbol+" on "+date, nil, "analysis.equityQuote", true)
	}
	return price, currency, nil
}
