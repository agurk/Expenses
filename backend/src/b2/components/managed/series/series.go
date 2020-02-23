package series

import (
	"b2/errors"
	"b2/manager"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Series represents a valuation of an series at a single
// point in time
type Series struct {
	sync.RWMutex
	ID               uint64 `json:"id"`
	AssetID          uint64 `json:"assetid"`
	Date             string `json:"date"`
	WholeAmount      int64  `json:"-"`
	FractionalAmount int64  `json:"-"`
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

// Merge is not implemented for Series
func (series *Series) Merge(newThing manager.Thing) error {
	newSeries := Cast(newThing)
	series.WholeAmount = newSeries.WholeAmount
	series.FractionalAmount = newSeries.FractionalAmount
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
		series.FractionalAmount, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return errors.Wrap(err, "series.parseAmount")
		}
	default:
		return errors.New("Badly formatted amount", nil, "series.parseAmount", true)
	}
	return nil
}

func (series *Series) amountString() string {
	if series.FractionalAmount == 0 {
		return fmt.Sprintf("%d", series.WholeAmount)
	}
	return fmt.Sprintf("%d.%d", series.WholeAmount, series.FractionalAmount)
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
