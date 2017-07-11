package rpc

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

var DefaultClient = &Client{push:make(chan interface{}, 8)}

type Client struct {
	connection *factory.Connection
	sync.RWMutex

	push chan interface{}
	logger *log.Entry
}

func (c *Client) GetConnection() *factory.Connection {
	c.RLock()
	defer c.RUnlock()
	return c.connection
}

func (c *Client) SetConnection(connection *factory.Connection) {
	c.Lock()
	if c.connection != nil {
		c.connection.Close()
	}
	c.connection = connection
	c.Unlock()
}

func (c *Client) PushLoop(conn *factory.Connection) {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Errorf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{}
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			key := cipher.NewPubKey(m[factory.MSG_PUBLIC_KEY_BEGIN:factory.MSG_PUBLIC_KEY_END])
			push.From = key.Hex()
			push.Msg = string(m[factory.MSG_META_END:])
			c.push <- push
		}
	}
}