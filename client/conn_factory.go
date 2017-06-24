package client

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/net/msg"
	"sync"
	"log"
	"errors"
)

type ClientConnection struct {
	key    cipher.PubKey
	client *Client
	In     chan []byte
	Out    chan []byte
}

func (c *ClientConnection) WriteLoop() error {
	for {
		select {
		case d := <-c.Out:
			err := c.client.conn.WriteSlice(c.key[:], d)
			if err != nil {
				return err
			}
		}
	}
}

type ClientConnectionFactory struct {
	client      *Client
	fieldsMutex *sync.RWMutex

	connections      map[string]*ClientConnection
	connectionsMutex *sync.RWMutex
}

func NewClientConnectionFactory() *ClientConnectionFactory {
	return &ClientConnectionFactory{fieldsMutex: new(sync.RWMutex), connectionsMutex: new(sync.RWMutex)}
}

var (
	ErrConnected = errors.New("factory connected")
)

func (factory *ClientConnectionFactory) Connect(network, address string, key cipher.PubKey) error {
	factory.fieldsMutex.RLock()
	if factory.client != nil {
		factory.fieldsMutex.RUnlock()
		return ErrConnected
	}
	factory.fieldsMutex.RUnlock()

	c := New()
	factory.fieldsMutex.Lock()
	factory.client = c
	factory.connections = make(map[string]*ClientConnection)
	factory.fieldsMutex.Unlock()

	err := c.Connect(network, address)
	if err != nil {
		return err
	}
	go func() {
		c.Loop()
		factory.close()
	}()

	go factory.dispatch()

	err = c.Reg(key)
	if err != nil {
		return err
	}
	return nil
}

func (factory *ClientConnectionFactory) close() {
	factory.connectionsMutex.Lock()
	for _, v := range factory.connections {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("ClientConnectionFactory close err %v", err)
				}
			}()
			close(v.In)
			close(v.Out)
		}()
	}
	factory.connectionsMutex.Unlock()

	factory.fieldsMutex.Lock()
	factory.client = nil
	factory.fieldsMutex.Unlock()
}

func (factory *ClientConnectionFactory) dispatch() {
	in := factory.client.conn.GetChanIn()
	for {
		select {
		case d, ok := <-in:
			if !ok {
				return
			}
			if len(d) < 33 {
				log.Printf("data len < 33 %x", d)
				continue
			}
			pubkey := cipher.NewPubKey(d[:msg.MSG_PUBKEY_SIZE])
			factory.connectionsMutex.RLock()
			conn, ok := factory.connections[pubkey.Hex()]
			factory.connectionsMutex.RUnlock()
			if !ok {
				log.Printf("conn not exists %s", pubkey.Hex())
				continue
			}
			conn.In <- d[msg.MSG_PUBKEY_SIZE:]
		}
	}
}

func (factory *ClientConnectionFactory) Dial(key cipher.PubKey) *ClientConnection {
	connection := &ClientConnection{key: key, client: factory.client, In: make(chan []byte), Out: make(chan []byte)}

	factory.connectionsMutex.Lock()
	factory.connections[key.Hex()] = connection
	factory.connectionsMutex.Unlock()

	go connection.WriteLoop()
	return connection
}

type ClientDirectConnectionFactory struct {
	ClientConnectionFactory
}

func NewClientDirectConnectionFactory(client *Client) *ClientDirectConnectionFactory {
	return &ClientDirectConnectionFactory{ClientConnectionFactory{client: client}}
}

func (factory *ClientDirectConnectionFactory) GetConn(key cipher.PubKey) *ClientConnection {
	return &ClientConnection{key: key, client: factory.client}
}
