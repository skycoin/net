package factory

import (
	"sync"

	"fmt"

	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

type Connection struct {
	*factory.Connection
	key         cipher.PubKey
	service     string
	fieldsMutex sync.RWMutex

	in chan []byte
}

func NewConnection(c *factory.Connection) *Connection {
	return &Connection{Connection: c, in: make(chan []byte)}
}

func (c *Connection) SetKey(key cipher.PubKey) {
	c.fieldsMutex.Lock()
	c.key = key
	c.fieldsMutex.Unlock()
	c.SetContextLogger(c.GetContextLogger().WithField("pubkey", key.Hex()))
}

func (c *Connection) GetKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.key
}

func (c *Connection) SetService(s string) {
	c.fieldsMutex.Lock()
	c.service = s
	c.fieldsMutex.Unlock()
	c.SetContextLogger(c.GetContextLogger().WithField("service", fmt.Sprintf("%s %x", s, s)))
}

func (c *Connection) GetService() string {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.service
}

func (c *Connection) Reg() error {
	return c.Write(GenRegMsg())
}

func (c *Connection) OfferService(service string) error {
	return c.Write(GenOfferServiceMsg(service))
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.key, to, msg))
}

func (c *Connection) SendCustom(msg []byte) error {
	return c.Write(GenCustomMsg(msg))
}

func (c *Connection) ReadLoop() error {
	go c.Connection.ReadLoop()
	for {
		select {
		case m, ok := <-c.Connection.GetChanIn():
			if !ok {
				return nil
			}

			c.in <- m
		}
	}
}
