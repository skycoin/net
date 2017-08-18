package factory

import (
	"sync"

	"encoding/json"
	"errors"

	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

type Connection struct {
	*factory.Connection
	factory     *MessengerFactory
	key         cipher.PubKey
	keySetCond  *sync.Cond
	keySet      bool
	services    *NodeServices
	fieldsMutex sync.RWMutex

	in           chan []byte
	disconnected chan struct{}

	// callbacks

	// call after received response for FindServiceNodesByKeys
	findServiceNodesByKeysCallback func(map[string][]string)

	// call after received response for FindServiceNodesByAttributes
	findServiceNodesByAttributesCallback func([]cipher.PubKey)
}

// Used by factory to spawn connections for server side
func newConnection(c *factory.Connection, factory *MessengerFactory) *Connection {
	connection := &Connection{Connection: c, factory: factory}
	c.RealObject = connection
	connection.keySetCond = sync.NewCond(connection.fieldsMutex.RLocker())
	return connection
}

// Used by factory to spawn connections for client side
func newClientConnection(c *factory.Connection, factory *MessengerFactory) *Connection {
	connection := &Connection{Connection: c, factory: factory, in: make(chan []byte), disconnected: make(chan struct{})}
	c.RealObject = connection
	connection.keySetCond = sync.NewCond(connection.fieldsMutex.RLocker())
	go func() {
		connection.preprocessor()
		close(connection.disconnected)
	}()
	return connection
}

func (c *Connection) WaitForDisconnected() {
	<-c.disconnected
}

func (c *Connection) SetKey(key cipher.PubKey) {
	c.fieldsMutex.Lock()
	c.key = key
	c.keySet = true
	c.fieldsMutex.Unlock()
	c.keySetCond.Broadcast()
}

func (c *Connection) IsKeySet() bool {
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	return c.keySet
}

func (c *Connection) GetKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	if !c.keySet {
		c.keySetCond.Wait()
	}
	return c.key
}

func (c *Connection) setServices(s *NodeServices) {
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	c.services = s
}

func (c *Connection) GetServices() *NodeServices {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.services
}

func (c *Connection) Reg() error {
	return c.Write(GenRegMsg())
}

func (c *Connection) UpdateServices(ns *NodeServices) error {
	if ns == nil || len(ns.Services) < 1 {
		return errors.New("invalid arguments")
	}
	js, err := json.Marshal(ns)
	if err != nil {
		return err
	}
	err = c.WriteOP(OP_OFFER_SERVICE, js)
	if err != nil {
		return err
	}
	c.setServices(ns)
	return nil
}

func (c *Connection) OfferService(attr string) error {
	return c.UpdateServices(&NodeServices{Services: []*Service{{Key: c.GetKey(), Attributes: []string{attr}}}})
}

func (c *Connection) FindServiceNodesByAttributes(attrs ...string) (err error) {
	js, err := json.Marshal(&queryByAttrs{Attrs: attrs})
	if err != nil {
		return
	}
	err = c.WriteOP(OP_QUERY_BY_ATTRS, js)
	return
}

func (c *Connection) FindServiceNodesByKeys(keys []cipher.PubKey) (err error) {
	js, err := json.Marshal(&query{Keys: keys})
	if err != nil {
		return
	}
	err = c.WriteOP(OP_QUERY_SERVICE_NODES, js)
	return
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.key, to, msg))
}

func (c *Connection) SendCustom(msg []byte) error {
	return c.WriteOP(OP_CUSTOM, msg)
}

func (c *Connection) preprocessor() (err error) {
	defer func() {
		if e := recover(); e != nil {
			c.GetContextLogger().Debugf("panic in preprocessor %v", e)
		}
		if err != nil {
			c.GetContextLogger().Debugf("preprocessor err %v", err)
		}
		c.Close()
	}()
	for {
		select {
		case m, ok := <-c.Connection.GetChanIn():
			if !ok {
				return
			}
			c.GetContextLogger().Debugf("read %x", m)
			if len(m) < MSG_HEADER_END {
				return
			}
			opn := m[MSG_OP_BEGIN]
			if opn&RESP_PREFIX > 0 {
				i := int(opn &^ RESP_PREFIX)
				r := getResp(i)
				if r != nil {
					body := m[MSG_HEADER_END:]
					if len(body) > 0 {
						err = json.Unmarshal(body, r)
						if err != nil {
							return
						}
					}
					err = r.Execute(c)
					if err != nil {
						return
					}
					putResp(i, r)
					continue
				}
			}

			c.in <- m
		}
	}
}

func (c *Connection) GetChanIn() <-chan []byte {
	if c.in == nil {
		return c.Connection.GetChanIn()
	}
	return c.in
}

func (c *Connection) Close() {
	c.keySetCond.Broadcast()
	c.fieldsMutex.Lock()
	defer c.fieldsMutex.Unlock()
	if c.keySet {
		c.factory.unregister(c.key, c)
	}
	if c.in != nil {
		close(c.in)
		c.in = nil
	}
	c.Connection.Close()
}

func (c *Connection) WriteOP(op byte, body []byte) error {
	data := make([]byte, MSG_HEADER_END+len(body))
	data[MSG_OP_BEGIN] = op
	copy(data[MSG_HEADER_END:], body)
	return c.Write(data)
}
