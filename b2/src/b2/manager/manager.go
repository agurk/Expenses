package manager

type ManagerInterface interface {
    Get(uint64) (Thing, error)
    Save(Thing) error
    Overwrite(Thing) (Thing, error)
    GetMultiple(string, string) ([]Thing, error)
}

type Thing interface {
    Type() string
    RLock()
    RUnlock()
    Lock()
    Unlock()
    GetID() uint64
}

