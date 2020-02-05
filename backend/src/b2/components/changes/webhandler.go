package changes

import (
	"b2/backend"
	"b2/errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Changes struct {
	Path    string
	backend *backend.Backend
	sync.RWMutex
	cons []*connection
}

const (
	// time in ns to keep the connection alive for before culling
	minKeepAlive  = 5 * time.Minute
	minCheckTime  = 5 * time.Minute
	changedMsg    = "changed"
	checkMsg      = "check"
	ExpenseEvent  = 1
	DocumentEvent = 2
)

type connection struct {
	conn     net.Conn
	lastSeen time.Time
	watching int
}

func Instance(path string, backend *backend.Backend) *Changes {
	changes := new(Changes)
	changes.Path = path
	changes.backend = backend
	go changes.Listen()
	go changes.checkAlive()
	return changes
}

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
	fmt.Println("setting up:", conex)
	go c.read(conex)
}

func (c *Changes) Listen() {
	for {
		event := <-c.backend.Change
		fmt.Println("Got change", c.cons)
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
		fmt.Println("checking", c.cons)
		msg := []byte(checkMsg)
		for _, conex := range c.cons {
			if time.Now().Sub(conex.lastSeen) > minKeepAlive {
				c.RUnlock()
				c.deRegisterConn(conex)
				c.RLock()
			} else {
				if err := wsutil.WriteServerText(conex.conn, msg); err != nil {
					errors.Print(errors.Wrap(err, "changes.checkAlive()"))
				}
			}
		}
		c.RUnlock()
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
	for {
		reader := wsutil.NewReader(conex.conn, ws.StateServerSide)
		msg, err := reader.NextFrame()
		if err != nil {
			errors.Print(errors.Wrap(err, "changes.read"))
		}
		if msg.OpCode == ws.OpClose {
			fmt.Println("closing conn")
			c.deRegisterConn(conex)
			return
		}
		conex.lastSeen = time.Now()
		_, err = ioutil.ReadAll(reader)
		if err != nil {
			errors.Print(errors.Wrap(err, "changes.read"))
			c.deRegisterConn(conex)
			return
		}
	}
}
