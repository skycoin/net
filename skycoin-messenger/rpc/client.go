package rpc

import (
	"github.com/skycoin/net/client"
	"sync"
	"log"
	"github.com/skycoin/net/skycoin-messenger/msg"
)

var DefaultClient = &Client{push:make(chan msg.PushMsg, 8)}

type Client struct {
	factory *client.ClientConnectionFactory
	sync.RWMutex

	push chan msg.PushMsg
}

func (c *Client) GetFactory() *client.ClientConnectionFactory {
	c.RLock()
	defer c.RUnlock()
	return c.factory
}

func (c *Client) SetFactory(factory *client.ClientConnectionFactory) {
	c.Lock()
	if c.factory != nil {
		c.factory.Close()
	}
	c.factory = factory
	c.Unlock()
}

func (c *Client) PushLoop(conn *client.ClientConnection, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{PublicKey: conn.Key.Hex(), Msg: string(data)}
	c.push <- *push
	for {
		select {
		case m, ok := <-conn.In:
			if !ok {
				return
			}
			push.Msg = string(m)
			c.push <- *push
		}
	}
}