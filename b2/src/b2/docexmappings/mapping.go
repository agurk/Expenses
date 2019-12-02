package docexmappings

import "sync"

type Mapping struct {
    ID uint64 `json:"id"`
    EID uint64 `json:"expenseId"`
    DID uint64 `json:"documentId"`
    Confirmed bool `json:"confirmed"`
    sync.RWMutex
}

func (mapping *Mapping) Type() string {
    return "docexpensemapping"
}

func (mapping *Mapping) GetID() uint64 {
    return mapping.ID
}

