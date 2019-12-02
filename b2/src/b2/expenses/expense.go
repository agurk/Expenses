package expenses 

import (
    "sync"
    "b2/docexmappings"
    "b2/rawexmappings"
)

type Expense struct {
    ID uint64 `json:"id"`
    TransactionReference string `json:"transactionReference"`
    Description string `json:"description"`
    DetailedDescription string `json:"detailedDescription"`
    AccountID int `json:"accountId"`
    Date string `json:"date"`
    ProcessDate string `json:"processDate"`
    Amount float64 `json:"amount"`
    Currency string `json:"currency"`
    FX FXProperties `json:"fx"`
    Commission int64 `json:"commission"`
    Metadata ExMeta `json:"metadata"`
    sync.RWMutex
    Documents []*docexmappings.Mapping `json:"documents"`
    Rawdata []*rawexmappings.Mapping `json:"raw"`
}

func (ex *Expense) Type() string {
    return "expense"
}

func (ex *Expense) GetID() uint64 {
    return ex.ID
}


type FXProperties struct {
    Amount float64 `json:"amount"`
    Currency string `json:"currency"`
    Rate float64 `json:"rate"`
}

type ExMeta struct {
    Confirmed bool `json:"confirmed"`
    Tagged int `json:"tagged"`
    Temporary bool `json:"temporary"`
    Modified string `json:"modified"`
    Classification int64 `json:"classification"`
}

/*
func (doc *Doc) MarshalJSON()  ([]byte, error) {
    buffer := bytes.NewBufferString("{")
    jsonValue, err := json.Marshal(doc.document.ID)
    if err != nil {
        return nil, err
    }
    buffer.WriteString(fmt.Sprintf("\"%s\":%s", "id", string(jsonValue)))
	buffer.WriteString("}")
    return buffer.Bytes(), nil 
}
*/
