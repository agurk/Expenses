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
    Documents []*mappings.Mapping `json:"documents"`
}

func (doc *Document) Type() string {
    return "document"
}

func (doc *Document) GetID() uint64 {
    return doc.ID
}

