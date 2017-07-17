package factory

import (
	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type Connection struct {
	*factory.Connection
	key         cipher.PubKey
	fieldsMutex sync.RWMutex
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

func (c *Connection) Reg() error {
	return c.Write(GenRegMsg())
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.key, to, msg))
}
