package documents

import (
    "sync"
    "b2/mappings"
)

type Document struct {
    ID uint64
    Filename string
    Deleted bool
    Date string
    Text string
    sync.RWMutex
    Expenses []*mappings.Mapping `json:"expenses"`
}

func (doc *Document) Type() string {
    return "document"
}

func (doc *Document) GetID() uint64 {
    return doc.ID
}

