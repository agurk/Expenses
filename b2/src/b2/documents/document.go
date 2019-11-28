package documents

import "sync"

type Document struct {
    ID uint64
    Filename string
    Deleted bool
    Date string
    Text string
    sync.RWMutex
    Expenses []*Expense
}

type Expense struct {
    ID uint64
    Confirmed bool
}

