package suggestions

import (
	"b2/backend"
	"b2/components/managed/expenses"
	"fmt"
)

type suggestion struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func getSuggestions(id uint64, b *backend.Backend) ([]*suggestion, error) {
	var result []*suggestion
	e, err := b.Expenses.Get(id)
	if err != nil {
		return nil, err
	}
	expense := e.(*expenses.Expense)
	for _, i := range expenses.Matches(expense, b.DB) {
		if i == expense.Metadata.Classification {
			continue
		}
		s := new(suggestion)
		s.Type = "classification"
		s.Value = fmt.Sprintf("%d", i)
		result = append(result, s)
	}
	return result, nil
}
