package series

import (
	"b2/errors"
	"b2/manager"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

// Series represents a valuation of an series at a single
// point in time. Internally the amount is represented as a
// pair of ints for the Whole and Fractional amount. There is also
// the fractional carrier that's the value added to the fractional amount
// to preserve any leading 0s.
//
// For example 2.03 would be parsed to:
//   WholeAmount = 2
//   FractionalAmount = 103
//   FractionalCarrier = 100
//
// The fractionalcarrier will always be 1 order of magnitude more than the
// fractionalamount
type Series struct {
	sync.RWMutex
	ID                uint64 `json:"id"`
	AssetID           uint64 `json:"assetid"`
	Date              string `json:"date"`
	WholeAmount       int64  `json:"-"`
	FractionalAmount  int64  `json:"-"`
	FractionalCarrier int64  `json:"-"`
}

// Cast a manager.Thing into a *Series or panic
func Cast(thing manager.Thing) *Series {
	series, ok := thing.(*Series)
	if !ok {
		panic("Non series passed to function")
	}
	return series
}

// Type returns a string representation of the series useful when using
// manager.Thing interfaces
func (series *Series) Type() string {
	return "series"
}

// GetID returns the ID of an series
func (series *Series) GetID() uint64 {
	return series.ID
}

// Merge is overwriting with the new values
func (series *Series) Merge(newThing manager.Thing) error {
	newSeries := Cast(newThing)
	series.WholeAmount = newSeries.WholeAmount
	series.FractionalAmount = newSeries.FractionalAmount
	series.FractionalCarrier = newSeries.FractionalCarrier
	return nil
}

// Overwrite is not implemented for Series
func (series *Series) Overwrite(newThing manager.Thing) error {
	return errors.New("Overwrite not implemented for series", errors.NotImplemented, "series.Overwrite", true)
}

// Check always returns nil errors for seriess
func (series *Series) Check() error {
	if series.Date == "" {
		return errors.New("Date must be specified", nil, "series.Check", true)
	}
	if series.AssetID == 0 {
		return errors.New("AssetID must be specified", nil, "series.Check", true)
	}
	if series.FractionalCarrier == 0 && series.FractionalAmount != 0 {
		return errors.New("Cannot have zero fractional carrier when there is an amount", nil, "series.Check", true)
	}
	return nil
}

func (series *Series) parseAmount(amount string) error {
	if amount == "" {
		return nil
	}
	parts := strings.Split(amount, ".")
	var err error
	switch len(parts) {
	case 1:
		series.WholeAmount, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "series.parseAmount")
		}
	case 2:
		series.WholeAmount, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return errors.Wrap(err, "series.parseAmount")
		}
		// 1 is added to the beginning of the string to allow for the carrier
		series.FractionalAmount, err = strconv.ParseInt("1"+parts[1], 10, 64)
		if err != nil {
			return errors.Wrap(err, "series.parseAmount")
		}
		// 10 as it's the order of magnitude bigger (as set by the 1 above)
		series.FractionalCarrier = int64(math.Pow10(len(parts[1])))
	default:
		return errors.New("Badly formatted amount", nil, "series.parseAmount", true)
	}
	return nil
}

func (series *Series) amountString() string {
	if series.FractionalAmount == 0 {
		return fmt.Sprintf("%d", series.WholeAmount)
	}
	fAmt := fmt.Sprintf("%d", series.FractionalAmount)
	return fmt.Sprintf("%d.%s", series.WholeAmount, fAmt[1:])
}

// AmountFloat returns the float representation from the series
// This amount is not linked to the series, so altering it will
// have no effect on the underlying series
func (series *Series) AmountFloat() float64 {
	if series.FractionalAmount == 0 {
		return float64(series.WholeAmount)
	}
	// minus 1 as that will be the residual amount of the carrier
	return float64(series.WholeAmount) + (float64(series.FractionalAmount) / float64(series.FractionalCarrier)) - 1
}

// MarshalJSON is to deal with amounts having decimal points in the real world
func (series *Series) MarshalJSON() ([]byte, error) {
	type Alias Series
	return json.Marshal(&struct {
		Amount string `json:"amount"`
		*Alias
	}{
		Amount: series.amountString(),
		Alias:  (*Alias)(series),
	})
}

// UnmarshalJSON is to deal with amounts having decimal points in the real world
func (series *Series) UnmarshalJSON(data []byte) error {
	type Alias Series
	aux := &struct {
		Amount string `json:"amount"`
		*Alias
	}{
		Alias: (*Alias)(series),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "series.UnmarshallJSON")
	}
	err := series.parseAmount(aux.Amount)
	if err != nil {
		return errors.Wrap(err, "series.UnmarshallJSON(parseAmount)")
	}
	return nil
}
