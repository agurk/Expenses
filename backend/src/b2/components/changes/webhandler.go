package changes

import (
	"b2/backend"
	"b2/errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// Changes is a notificition system that will send a message down a
// websocket to inform the client if a change has happened on the server.
type Changes struct {
	Path    string
	backend *backend.Backend
	sync.RWMutex
	cons []*connection
}

const (
	// time in ns to keep the connection alive for before culling
	minKeepAlive = 5 * time.Minute
	minCheckTime = 5 * time.Minute
	changedMsg   = "changed"
	checkMsg     = "check"
	// ExpenseEvent is sent to the backend changes channel to notify that
	// a change has occured to an expense
	ExpenseEvent = 1
	// DocumentEvent is sent to teh backend changes channel to notify that
	// a change has occured to a document
	DocumentEvent = 2
)

type connection struct {
	conn     net.Conn
	lastSeen time.Time
	watching int
}

// Instance returns an instantiated changes backend
func Instance(path string, backend *backend.Backend) *Changes {
	changes := new(Changes)
	changes.Path = path
	changes.backend = backend
	go changes.listen()
	go changes.checkAlive()
	return changes
}

// Handle is a standard net/http handler to deal with incoming requests
// to moniter changes
func (c *Changes) Handle(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		errors.Print(errors.Wrap(err, "changes.Handle"))
	}
	conex := new(connection)
	switch r.URL.Path[len(c.Path):] {
	case "expenses":
		conex.watching = ExpenseEvent
	case "documents":
		conex.watching = DocumentEvent
	}
	conex.conn = conn
	conex.lastSeen = time.Now()
	c.registerConn(conex)
	go c.read(conex)
}

func (c *Changes) listen() {
	for {
		event := <-c.backend.Change
		c.notify(event)
	}
}

func (c *Changes) notify(event int) {
	c.RLock()
	defer c.RUnlock()
	msg := []byte(changedMsg)
	for _, conex := range c.cons {
		if conex.watching == event {
			if err := wsutil.WriteServerText(conex.conn, msg); err != nil {
				errors.Print(errors.Wrap(err, "changes.notify"))
			}
		}
	}
}

func (c *Changes) checkAlive() {
	for {
		c.RLock()
		var dead []*connection
		fmt.Println("checking", c.cons)
		msg := []byte(checkMsg)
		for _, conex := range c.cons {
			if time.Now().Sub(conex.lastSeen) > minKeepAlive {
				dead = append(dead, conex)
			} else {
				if err := wsutil.WriteServerText(conex.conn, msg); err != nil {
					errors.Print(errors.Wrap(err, "changes.checkAlive()"))
				}
			}
		}
		c.RUnlock()
		for _, conex := range dead {
			c.deRegisterConn(conex)
		}
		time.Sleep(minCheckTime)
	}
}

func (c *Changes) registerConn(conex *connection) {
	fmt.Println("Registering", conex)
	c.Lock()
	c.cons = append(c.cons, conex)
	c.Unlock()
}

func (c *Changes) deRegisterConn(conex *connection) {
	fmt.Println("Deregistering", conex)
	c.Lock()
	defer c.Unlock()
	i := -1
	for j, val := range c.cons {
		if val == conex {
			i = j
			break
		}
	}
	if i == -1 {
		return
	}
	c.cons[i] = c.cons[len(c.cons)-1]
	c.cons[len(c.cons)-1] = nil
	c.cons = c.cons[:len(c.cons)-1]
}

func (c *Changes) read(conex *connection) {
	reader := wsutil.NewReader(conex.conn, ws.StateServerSide)
	for {
		msg, err := reader.NextFrame()
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF found")
				c.deRegisterConn(conex)
				return
			}
			errors.Print(errors.Wrap(err, "changes.read (NextFrame)"))
			c.deRegisterConn(conex)
			return
		}
		if msg.OpCode == ws.OpClose {
			fmt.Println("closing conn")
			c.deRegisterConn(conex)
			return
		}
		conex.lastSeen = time.Now()
		err = reader.Discard()
		if err != nil {
			errors.Print(errors.Wrap(err, "changes.read (Discard)"))
			c.deRegisterConn(conex)
			return
		}
	}
	c.deRegisterConn(conex)
}
