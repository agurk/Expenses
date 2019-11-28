package classifications 

import "sync"

type Classification struct {
    ID uint64
    Description string
    Hidden bool
    From string
    To string
    sync.RWMutex
}

func (classification *Classification) Type() string {
    return "classification"
}

func (classification *Classification) GetID() uint64 {
    return classification.ID
}

