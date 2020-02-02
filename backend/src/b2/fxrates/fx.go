package fxrates

import (
	"b2/errors"
	"database/sql"
	"time"
)

type FxValues struct {
	// map of [ccypair][date][rate]
	values         map[string]map[string]float64
	db             *sql.DB
	lookbackPeriod int
}

func (fx *FxValues) Initalize(db *sql.DB) {
	fx.db = db
	fx.values = make(map[string]map[string]float64)
	fx.loadRates()
	fx.lookbackPeriod = 30
}

func (fx *FxValues) loadRates() error {
	rows, err := fx.db.Query(`select
								date,
								ccy1,
								ccy2,
								rate
							from
								_FXRates`)
	//where
	//		strftime(date) >= date('$1','start of year')
	//		and strftime(date) <= date('$1','start of year', '+12 months')
	//		and (
	//			(ccy1 = $2 and ccy2 = $3 )
	//			or (ccy1 = $3 and ccy2 = $2)
	//		)`, year, ccy1, ccy2)
	if err != nil {
		return errors.Wrap(err, "fxrates.loadRates")
	}
	defer rows.Close()
	for rows.Next() {
		var date, ccy1, ccy2 string
		var rate float64
		err = rows.Scan(&date, &ccy1, &ccy2, &rate)
		if err != nil {
			return errors.Wrap(err, "fxrates.loadRates")
		}
		date = date[:10]
		key := ccy1 + ccy2
		if _, ok := fx.values[key]; !ok {
			if _, ok = fx.values[ccy2+ccy1]; !ok {
				fx.values[key] = make(map[string]float64)
			} else {
				key = ccy2 + ccy1
				rate = 1 / rate
			}
		}
		fx.values[key][date] = rate
	}
	return nil
}

func (fx *FxValues) Get(dateIn, ccy1, ccy2 string) (float64, error) {
	if ccy1 == ccy2 {
		return 1, nil
	}
	date, _ := time.Parse("2006-01-02", dateIn)
	for i := 0; i < fx.lookbackPeriod; i++ {
		if _, ok := fx.values[ccy1+ccy2]; ok {
			if value, ok := fx.values[ccy1+ccy2][date.Format("2006-01-02")]; ok {
				return value, nil
			}
		} else if _, ok = fx.values[ccy2+ccy1]; ok {
			if value, ok := fx.values[ccy2+ccy1][date.Format("2006-01-02")]; ok {
				return (1 / value), nil
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	// todo: try loading fx rate
	return 0, errors.New("FX rate not found", nil, "fxrates.Get")
}
