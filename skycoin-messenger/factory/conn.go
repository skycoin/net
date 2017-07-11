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
}

func (c *Connection) GetKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.key
}

func (c *Connection) Reg(key cipher.PubKey) error {
	c.SetKey(key)
	c.SetContextLogger(c.GetContextLogger().WithField("pubkey", key.Hex()))
	return c.Write(GenRegMsg(key))
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.key, to, msg))
}
