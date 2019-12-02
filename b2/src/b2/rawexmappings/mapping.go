package rawexmappings

import "sync"

type Mapping struct {
    ID uint64 `json:"id"`
    EID uint64 `json:"expenseId"`
    RID uint64 `json:"rawId"`
    sync.RWMutex
}

func (mapping *Mapping) Type() string {
    return "rawexpensemapping"
}

func (mapping *Mapping) GetID() uint64 {
    return mapping.ID
}

