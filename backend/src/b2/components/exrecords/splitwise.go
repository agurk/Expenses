package exrecords

import (
	"b2/components/managed/expenses"
	"b2/errors"
	"b2/moneyutils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type group struct {
	Name    string            `json:"name"`
	Members map[uint64]string `json:"members"`
}

func getSplitwiseGroups(swSecret string) (*map[uint64]group, error) {
	groups := make(map[uint64]group)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://secure.splitwise.com/api/v3.0/get_groups", nil)
	if err != nil {
		return nil, errors.Wrap(err, "exrecords.getSplitwiseGroups")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", swSecret))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "exrecords.getSplitwiseGroups")
	}
	// todo: check status
	decoder := json.NewDecoder(resp.Body)
	type swMember struct {
		ID    uint64 `json:"id"`
		FName string `json:"first_name"`
		LName string `json:"last_name"`
	}
	type swGroup struct {
		ID      uint64     `json:"id"`
		Name    string     `json:"name"`
		Members []swMember `json:"members"`
	}
	type swTop struct {
		Groups []swGroup `json:"Groups"`
	}
	var swGroups swTop
	err = decoder.Decode(&swGroups)
	if err != nil {
		return nil, errors.Wrap(err, "exrecords.getSplitwiseGroups")
	}
	for _, i := range swGroups.Groups {
		var grp group
		if len(i.Name) > 20 {
			grp.Name = i.Name[:20]
			grp.Name += "..."
		} else {
			grp.Name = i.Name
		}
		grp.Members = make(map[uint64]string)
		for _, j := range i.Members {
			grp.Members[j.ID] = j.FName + " " + j.LName
		}
		groups[i.ID] = grp
	}
	return &groups, nil
}

func splitwiseData(data *postData, e *expenses.Expense, swUser uint64) (url.Values, int64) {
	formattedAmount, err := moneyutils.CurrencyAmount(e.Amount, e.Currency)
	if err != nil {
		errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
		return nil, 0
	}
	formattedAmount *= -1
	leftover := (e.Amount * -1) % int64(len(data.Members))
	amount := float64(e.Amount+leftover) / (float64(len(data.Members)) * -100)
	fraction := float64(leftover) / 100
	values := url.Values{
		"cost":          {fmt.Sprintf("%.2f", formattedAmount)},
		"currency_code": {e.Currency},
		// Timezone offset seems to be 7hrs for some reason
		"date":        {e.Date + "T07:00:00Z"},
		"group_id":    {fmt.Sprintf("%d", data.Group)},
		"description": {e.Description},
	}
	seenUser := false
	userFraction := false
	for i, user := range data.Members {
		values.Add(fmt.Sprintf("users__%d__user_id", i), fmt.Sprintf("%d", user))
		if int64(i) < leftover {
			amt, err := moneyutils.StringFloat(amount+fraction, e.Currency)
			if err != nil {
				errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
			}
			values.Add(fmt.Sprintf("users__%d__owed_share", i), amt)
			if user == swUser {
				userFraction = true
			}
		} else {
			amt, err := moneyutils.StringFloat(amount, e.Currency)
			if err != nil {
				errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
			}
			values.Add(fmt.Sprintf("users__%d__owed_share", i), amt)
		}
		paidAmount := 0.0
		if user == swUser {
			seenUser = true
			paidAmount = formattedAmount
		}
		amt, err := moneyutils.StringFloat(paidAmount, e.Currency)
		if err != nil {
			errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
		}
		values.Add(fmt.Sprintf("users__%d__paid_share", i), amt)
	}
	if !seenUser {
		amount = 0
		i := len(data.Members) + 1
		values.Add(fmt.Sprintf("users__%d__user_id", i), fmt.Sprintf("%d", swUser))
		amt, err := moneyutils.StringFloat(formattedAmount, e.Currency)
		if err != nil {
			errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
		}
		values.Add(fmt.Sprintf("users__%d__paid_share", i), amt)
		amt, err = moneyutils.StringFloat(amount, e.Currency)
		if err != nil {
			errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
		}
		values.Add(fmt.Sprintf("users__%d__owed_share", i), amt)
	}
	if userFraction {
		amt, err := moneyutils.ParseFloat(amount+fraction, e.Currency)
		if err != nil {
			errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
		}
		return values, -1 * amt
	}
	amt, err := moneyutils.ParseFloat(amount+fraction, e.Currency)
	if err != nil {
		errors.Print(errors.Wrap(err, "splitwise.splitwiseData"))
	}
	return values, -1 * amt
}

func addSplitwiseExpense(dataIn *postData, e *expenses.Expense, swSecret string, swUser uint64) error {
	client := &http.Client{}
	data, userAmount := splitwiseData(dataIn, e, swUser)
	req, err := http.NewRequest("POST", "https://secure.splitwise.com/api/v3.0/create_expense",
		strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "exrecords.addSplitwiseExpense")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", swSecret))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "exrecords.addSplitwiseExpense")
	}
	if resp.Status != "200 OK" {
		return errors.New(fmt.Sprintf("Unable to create expense on splitwise, error: %s", resp.Status), nil, "exrecord.addSplitwiseExpense")
	}
	type id struct {
		ID uint64 `json:"id"`
	}
	type base struct {
		Base []string `json"base"`
	}
	type exes struct {
		Expenses []id `json:"expenses"`
		Errors   base `json:"errors"`
	}
	response := new(exes)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return errors.Wrap(err, "exrecords.addSplitwiseExpense")
	}
	if len(response.Errors.Base) > 0 {
		// todo: what about more than one error?
		return errors.New(response.Errors.Base[0], nil, "exrecords.addSplitwiseExpense")
	}
	newRecord := new(expenses.ExternalRecord)
	newRecord.Reference = fmt.Sprintf("%d", response.Expenses[0].ID)
	newRecord.Type = "splitwise"
	newRecord.FullAmount = e.Amount
	e.Amount = userAmount
	e.ExternalRecords = append(e.ExternalRecords, newRecord)
	return nil
}
