package expenses

import (
	"b2/components/managed/docexmappings"
	"b2/errors"
	"b2/manager"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

type Expense struct {
	sync.RWMutex
	deleted              bool                     `json:-`
	ID                   uint64                   `json:"id"`
	TransactionReference string                   `json:"transactionReference"`
	Description          string                   `json:"description"`
	DetailedDescription  string                   `json:"detailedDescription"`
	AccountID            uint                     `json:"accountId"`
	Date                 string                   `json:"date"`
	ProcessDate          string                   `json:"processDate"`
	Currency             string                   `json:"currency"`
	Amount               int64                    `json:"-"`
	FX                   FXProperties             `json:"fx"`
	Commission           int64                    `json:"-"`
	Metadata             ExMeta                   `json:"metadata"`
	Documents            []*docexmappings.Mapping `json:"documents"`
	ExternalRecords      []*ExternalRecord        `json:"externalRecords"`
}

func (ex *Expense) Type() string {
	return "expense"
}

func (ex *Expense) GetID() uint64 {
	return ex.ID
}

func (ex *Expense) Overwrite(newThing manager.Thing) error {
	expense, ok := newThing.(*Expense)
	if !ok {
		panic("Non expense passed to overwrite function")
	}
	expense.RLock()
	ex.Lock()
	ex.TransactionReference = expense.TransactionReference
	ex.Description = expense.Description
	ex.DetailedDescription = expense.DetailedDescription
	ex.AccountID = expense.AccountID
	ex.Date = expense.Date
	ex.ProcessDate = expense.ProcessDate
	ex.Amount = expense.Amount
	ex.Currency = expense.Currency
	ex.Commission = expense.Commission
	ex.FX = expense.FX
	ex.Metadata = expense.Metadata
	//ex.Documents = expense.Documents
	expense.RUnlock()
	ex.Unlock()
	return nil
}

func (ex *Expense) Merge(newThing manager.Thing) error {
	expense, ok := newThing.(*Expense)
	if !ok {
		panic("Non expense passed to overwrite function")
	}
	ex.Lock()
	expense.RLock()
	defer ex.Unlock()
	defer expense.RUnlock()
	ex.mergeStringField(&ex.TransactionReference, &expense.TransactionReference, "Transaction Reference")
	ex.mergeStringField(&ex.Description, &expense.Description, "Description")
	ex.mergeStringField(&ex.DetailedDescription, &expense.DetailedDescription, "Detailed Description")
	// skipping date assuming ex has the correct one
	ex.mergeStringField(&ex.ProcessDate, &expense.ProcessDate, "Processed Date")
	ex.mergeStringField(&ex.Currency, &expense.Currency, "Currency")
	ex.mergeStringField(&ex.FX.Currency, &expense.FX.Currency, "FX Currency")
	ex.mergeIntField(&ex.Amount, &expense.Amount, "Amount")
	ex.mergeIntField(&ex.Commission, &expense.Commission, "Commission")
	ex.mergeFloatField(&ex.FX.Amount, &expense.FX.Amount, "FX Amount")
	ex.mergeFloatField(&ex.FX.Rate, &expense.FX.Rate, "FX Rate")
	// preserve if the expense has ever been confirmed
	if ex.Metadata.Confirmed || expense.Metadata.Confirmed {
		ex.Metadata.Confirmed = true
	}
	ex.Metadata.Temporary = expense.Metadata.Temporary
	ex.AccountID = expense.AccountID
	// todo: tagged, modified, classification
	return nil
}

func (ex *Expense) mergeStringField(oldValue, newValue *string, fieldName string) {
	if (*oldValue != "") && (*oldValue != *newValue) {
		ex.Metadata.OldValues += fieldName + " changed from " + *oldValue + "\n"
	}
	*oldValue = *newValue
}

func (ex *Expense) mergeFloatField(oldValue, newValue *float64, fieldName string) {
	if (*oldValue != 0) && (*oldValue != *newValue) {
		ex.Metadata.OldValues += fieldName + " changed from " + strconv.FormatFloat(*oldValue, 'f', -1, 64) + "\n"
	}
	*oldValue = *newValue
}

func (ex *Expense) mergeIntField(oldValue, newValue *int64, fieldName string) {
	if (*oldValue != 0) && (*oldValue != *newValue) {
		ex.Metadata.OldValues += fmt.Sprintf("%s changed from %d\n", fieldName, oldValue)
	}
	*oldValue = *newValue
}

func (ex *Expense) Check() error {
	ex.RLock()
	defer ex.RUnlock()
	if ex.deleted {
		return errors.New(fmt.Sprintf("Expense is deleted. Id: %i", ex.ID), nil, "expenses.Check")
	}
	// must have transaction reference if not temporary
	if !ex.Metadata.Temporary && ex.TransactionReference == "" {
		return errors.New(fmt.Sprintf("Transaction reference missing. Id: %i", ex.ID), nil, "expenses.Check")
	}
	// must be assigned to an account
	// todo: check if account is valid
	if ex.AccountID == 0 {
		return errors.New("Missing or invalid account id", nil, "expenses.Check")
	}
	if ex.Date == "" || ex.Description == "" {
		return errors.New("Missing date or description", nil, "expenses.Check")
	}
	return nil
}

type FXProperties struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate"`
}

type ExMeta struct {
	Confirmed      bool   `json:"confirmed"`
	Tagged         int    `json:"tagged"`
	Temporary      bool   `json:"temporary"`
	Modified       string `json:"modified"`
	Classification int64  `json:"classification"`
	OldValues      string `json:"oldValues"`
}

type ExternalRecord struct {
	Type       string `json:"type"`
	Reference  string `json:"reference"`
	FullAmount int64  `json:oldAmount"`
}

func toDisplayAmount(amount int64, ccy string) string {
	// todo: actually use ccy
	newAmount := float64(amount) / 100
	return fmt.Sprintf("%.2f", newAmount)
}

func fromDisplayAmount(amount string, oldAmount int64, ccy string) int64 {
	// we need this check, otherwise parsing a partial stream that doesn't have a matching
	// field will overwrite the existing value with 0
	if amount == "" {
		return oldAmount
	}
	// todo: actually use ccy and look to improve
	val, _ := strconv.ParseFloat(amount, 64)
	return int64(val * 100)
}

func (ex *Expense) MarshalJSON() ([]byte, error) {
	type Alias Expense
	return json.Marshal(&struct {
		Amount     string `json:"amount"`
		Commission string `json:"commission"`
		*Alias
	}{
		Amount:     toDisplayAmount(ex.Amount, ex.Currency),
		Commission: toDisplayAmount(ex.Commission, ex.Currency),
		Alias:      (*Alias)(ex),
	})
}

func (ex *Expense) UnmarshalJSON(data []byte) error {
	type Alias Expense
	aux := &struct {
		Amount     string `json:"amount"`
		Commission string `json:"commission"`
		*Alias
	}{
		Alias: (*Alias)(ex),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.Wrap(err, "expenses.UnmarshalJSON")
	}
	ex.Amount = fromDisplayAmount(aux.Amount, ex.Amount, ex.Currency)
	ex.Commission = fromDisplayAmount(aux.Commission, ex.Commission, ex.Currency)
	return nil
}
