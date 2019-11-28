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

func (doc *Document) Type() string {
    return "document"
}

func (doc *Document) GetID() uint64 {
    return doc.ID
}

type Expense struct {
    ID uint64
    Confirmed bool
    Date string
    Description string
}

