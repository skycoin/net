package rpc

import (
	"sync"
	"log"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

var DefaultClient = &Client{push:make(chan msg.PushMsg, 8)}

type Client struct {
	connection *factory.Connection
	sync.RWMutex

	push chan msg.PushMsg
}

func (c *Client) GetConnection() *factory.Connection{
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

func (c *Client) PushLoop(conn *factory.Connection, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{PublicKey: conn.GetKey().Hex(), Msg: string(data)}
	c.push <- *push
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			push.Msg = string(m)
			c.push <- *push
		}
	}
}