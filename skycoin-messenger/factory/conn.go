package factory

import (
	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

type Connection struct {
	*factory.Connection
	Key cipher.PubKey
}

func (c *Connection) Reg(key cipher.PubKey) error {
	c.Key = key
	return c.Write(GenRegMsg(key))
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.Key, to, msg))
}