package rawprocessor

import (
    "fmt"
)

type RawProcessor struct {
    processId chan uint64
}

func (rp *RawProcessor) Listen() {
    for {
        fmt.Println(<- rp.processId)
    }
}

func (rp *RawProcessor) Channel() chan uint64 {
    rp.processId = make (chan uint64)
    go rp.Listen()
    return rp.processId
}
