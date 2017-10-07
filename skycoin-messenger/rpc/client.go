package rpc

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/skycoin/src/cipher"
)

var DefaultClient = &Client{Push: make(chan interface{}, 8)}

type Client struct {
	sync.RWMutex
	factory *factory.MessengerFactory

	Push   chan interface{}
	Logger *log.Entry
}

func (c *Client) GetFactory() *factory.MessengerFactory {
	c.RLock()
	defer c.RUnlock()
	return c.factory
}

func (c *Client) SetFactory(factory *factory.MessengerFactory) {
	c.Lock()
	if c.factory != nil {
		c.factory.Close()
	}
	c.factory = factory
	c.Unlock()
}

func (c *Client) PushLoop(conn *factory.Connection) {
	defer func() {
		if err := recover(); err != nil {
			c.Logger.Errorf("PushLoop recovered err %v", err)
		}
	}()
	key := conn.GetKey()
	c.Push <- &msg.Reg{PublicKey: key.Hex()}
	push := &msg.PushMsg{}
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok || len(m) < factory.MSG_HEADER_END {
				return
			}
			op := m[factory.MSG_OP_BEGIN]
			switch op {
			case factory.OP_SEND:
				if len(m) < factory.SEND_MSG_META_END {
					continue
				}
				key := cipher.NewPubKey(m[factory.SEND_MSG_PUBLIC_KEY_BEGIN:factory.SEND_MSG_PUBLIC_KEY_END])
				push.From = key.Hex()
				push.Msg = string(m[factory.SEND_MSG_META_END:])
				c.Push <- push
			}
		}
	}
}
