package expenses 

import "sync"

type Expense struct {
    ID uint64
    TransactionReference string
    Description string
    DetailedDescription string
    AccountID int
    Date string
    ProcessDate string
    Amount float64
    Currency string
    FX FXProperties
    Commission int64
    Metadata ExMeta
    sync.RWMutex
    Documents []*Doc
}

func (ex *Expense) Type() string {
    return "expense"
}

func (ex *Expense) GetID() uint64 {
    return ex.ID
}


type FXProperties struct {
    Amount float64
    Currency string
    Rate float64
}

type ExMeta struct {
    Confirmed bool
    Tagged int
    Temporary bool
    Modified string
    Classification string
}

type Doc struct {
    ID uint64
    Confirmed bool
    Filename string
}

