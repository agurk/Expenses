package changes

import (
	"b2/backend"
	"fmt"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Changes struct {
	backend *backend.Backend
}

func Instance(backend *backend.Backend) *Changes {
	changes := new(Changes)
	changes.backend = backend
	return changes
}

func (changes *Changes) Handle(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	msg := []byte("changed")
	for {
		_ = <-changes.backend.Change
		fmt.Println("change happens")
		if err = wsutil.WriteServerText(conn, msg); err != nil {
			// handle error
		}
	}
}
