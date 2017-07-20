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
	coworkers   []cipher.PubKey
	fieldsMutex sync.RWMutex

	in chan []byte
}

func NewConnection(c *factory.Connection) *Connection {
	return &Connection{Connection: c}
}

func NewClientConnection(c *factory.Connection) *Connection {
	connection := &Connection{Connection: c, in: make(chan []byte)}
	go connection.preprocessor()
	return connection
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

func (c *Connection) GetCoworkers() []cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.coworkers
}

func (c *Connection) Reg() error {
	return c.Write(GenRegMsg())
}

func (c *Connection) OfferService(service string) error {
	c.SetService(service)
	return c.Write(GenOfferServiceMsg(service))
}

func (c *Connection) Send(to cipher.PubKey, msg []byte) error {
	return c.Write(GenSendMsg(c.key, to, msg))
}

func (c *Connection) SendCustom(msg []byte) error {
	return c.Write(GenCustomMsg(msg))
}

func (c *Connection) preprocessor() error {
	for {
		select {
		case m, ok := <-c.Connection.GetChanIn():
			if !ok {
				close(c.in)
				return nil
			}
			c.GetContextLogger().Debugf("read %x", m)
			if len(m) >= MSG_HEADER_END {
				switch m[MSG_OP_BEGIN] {
				case OP_REG:
					reg := m[MSG_HEADER_END:]
					if len(reg) < MSG_PUBLIC_KEY_SIZE {
						continue
					}
					key := cipher.NewPubKey(reg[:MSG_PUBLIC_KEY_SIZE])
					c.SetKey(key)
				case OP_OFFER_SERVICE:
					ks := m[MSG_HEADER_END:]
					kc := len(ks) / MSG_PUBLIC_KEY_SIZE
					if len(ks)%MSG_PUBLIC_KEY_SIZE != 0 || kc < 1 {
						continue
					}
					coworkers := make([]cipher.PubKey, kc)
					for i := 0; i < kc; i++ {
						key := cipher.NewPubKey(ks[i*MSG_PUBLIC_KEY_SIZE:(i+1)*MSG_PUBLIC_KEY_SIZE])
						coworkers[i] = key
					}
					c.fieldsMutex.Lock()
					c.coworkers = coworkers
					c.fieldsMutex.Unlock()
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