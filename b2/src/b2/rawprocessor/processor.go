package rawprocessor

import (
    "fmt"
    "strings"
    "strconv"
    "b2/manager"
    "b2/rawexpenses"
    "b2/expenses"
)

type RawProcessor struct {
    processId chan *rawexpenses.RawExpense
    m *manager.Manager
}

func (rp *RawProcessor) Channel() chan *rawexpenses.RawExpense {
    rp.processId = make (chan *rawexpenses.RawExpense)
    go rp.Listen()
    return rp.processId
}

func (rp *RawProcessor) AttachExManager(m *manager.Manager) {
    rp.m = m
}

func (rp *RawProcessor) Listen() {
    for {
        raw := <- rp.processId
        fmt.Println(raw)
        ex := processGeneric(raw)
        rp.m.Save(ex)
    }
}

// Processing a "Generic raw line" as defined in the original version of this 
//     0 TransactionDate
//     1 ProcessedDate
//     2 Description
//     3 Amount
//     4 DebitCredit
//     5 FXAmount
//     6 FXCCY
//     7 FXRate
//     8 Commission
//     9 RefID
//    10 Temporary
//    11 ExtraText
func processGeneric(raw *rawexpenses.RawExpense) *expenses.Expense  {
    // todo check for errors with parsing
    s := strings.Split(raw.Data, ";")
    expense := new (expenses.Expense)
    expense.Date = s[0]
    expense.ProcessDate = s[1]
    expense.Description = s[2]
    expense.Amount, _ = strconv.ParseFloat(s[3], 64)
    expense.FX.Amount, _ = strconv.ParseFloat(s[5], 64)
    expense.FX.Currency = s[6]
    expense.FX.Rate, _ = strconv.ParseFloat(s[7], 64)
    expense.Commission, _ = strconv.ParseFloat(s[8], 64)
    expense.TransactionReference = s[9]
    expense.Metadata.Temporary, _ = strconv.ParseBool(s[10])
    expense.DetailedDescription = s[11]
    if (s[4] == "DR") {
        expense.Amount *= -1
    }
    expense.Metadata.Classification = 1
    return expense
}

