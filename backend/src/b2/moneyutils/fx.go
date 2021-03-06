package moneyutils

import (
	"b2/errors"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

// FxValues is designed to allow money to be converted at the closing FX rate
// of the provided date
type FxValues struct {
	sync.RWMutex
	// map of [ccypair][date][rate]
	values         map[string]map[string]float64
	db             *sql.DB
	lookbackPeriod int
	lastLoad       time.Time
	checkRates     chan bool
}

const maxLoadTime = 6 * time.Hour

// Initalize loads the known fx values into the object, and other loading configuration
func (fx *FxValues) Initalize(db *sql.DB) {
	fx.db = db
	fx.values = make(map[string]map[string]float64)
	fx.checkRates = make(chan bool, 1000)
	fx.loadRates()
	fx.lookbackPeriod = 30
	go fx.ratesChecker()
}

func (fx *FxValues) ratesChecker() {
	for {
		_ = <-fx.checkRates
		if time.Now().Sub(fx.lastLoad) > maxLoadTime {
			err := fx.loadRates()
			errors.Print(err)
		}
	}
}

func (fx *FxValues) loadRates() error {
	fx.Lock()
	defer fx.Unlock()
	fmt.Println("rates: loading rates from db")
	fx.lastLoad = time.Now()
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
	defer rows.Close()
	if err != nil {
		return errors.Wrap(err, "fxrates.loadRates")
	}
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

// Rate takes in a date, currency from and currency to and returns the
// amount from that day
func (fx *FxValues) Rate(dateIn, ccy1in, ccy2in string) (float64, error) {
	fx.checkRates <- true
	ccy1 := strings.ToUpper(ccy1in)
	ccy2 := strings.ToUpper(ccy2in)
	if ccy1 == ccy2 {
		return 1, nil
	}

	fx.RLock()
	date, _ := time.Parse("2006-01-02", dateIn)
	for i := 0; i < fx.lookbackPeriod; i++ {
		if _, ok := fx.values[ccy1+ccy2]; ok {
			if value, ok := fx.values[ccy1+ccy2][date.Format("2006-01-02")]; ok {
				fx.RUnlock()
				return value, nil
			}
		} else if _, ok = fx.values[ccy2+ccy1]; ok {
			if value, ok := fx.values[ccy2+ccy1][date.Format("2006-01-02")]; ok {
				fx.RUnlock()
				return (1 / value), nil
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	fx.RUnlock()

	// todo: try loading fx rate
	if ccy1 != "USD" && ccy2 != "USD" {
		usdrate, err := fx.Rate(dateIn, ccy1, "USD")
		if err != nil {
			return 0, errors.New("FX rate not found for "+ccy1+ccy2+" on "+dateIn+" (including USD variant)", nil, "fxrates.Rate", true)
		}
		usdExchange, err := fx.Rate(dateIn, ccy2, "USD")
		if err != nil {
			return 0, errors.New("FX rate not found for "+ccy1+ccy2+" on "+dateIn+" (including USD variant)", nil, "fxrates.Rate", true)
		}
		return usdrate / usdExchange, nil
	}
	return 0, errors.New("FX rate not found for "+ccy1+ccy2+" on "+dateIn, nil, "fxrates.Rate", true)
}
