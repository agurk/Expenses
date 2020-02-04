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
	backend *backend.Backend
	sync.RWMutex
	cons []*connection
}

const (
	// time in ns to keep the connection alive for before culling
	MinKeepAlive = 5 * time.Minute
	MinCheckTime = 5 * time.Minute
	ChangedMsg   = "changed"
	CheckMsg     = "check"
)

type connection struct {
	conn     net.Conn
	lastSeen time.Time
}

func Instance(backend *backend.Backend) *Changes {
	changes := new(Changes)
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
	conex.conn = conn
	conex.lastSeen = time.Now()
	c.registerConn(conex)
	go c.read(conex)
}

func (c *Changes) Listen() {
	for {
		_ = <-c.backend.Change
		fmt.Println("Got change", c.cons)
		c.notify()
	}
}

func (c *Changes) notify() {
	c.RLock()
	defer c.RUnlock()
	msg := []byte(ChangedMsg)
	for _, conex := range c.cons {
		fmt.Println("change happens")
		if err := wsutil.WriteServerText(conex.conn, msg); err != nil {
			errors.Print(errors.Wrap(err, "changes.notify"))
		}
	}
}

func (c *Changes) checkAlive() {
	for {
		c.RLock()
		fmt.Println("checking", c.cons)
		msg := []byte(CheckMsg)
		for _, conex := range c.cons {
			if time.Now().Sub(conex.lastSeen) > MinKeepAlive {
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
		time.Sleep(MinCheckTime)
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
		}
	}
}
