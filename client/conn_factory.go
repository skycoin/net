package client

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/net/msg"
	"sync"
	"log"
	"errors"
)

type ClientConnection struct {
	Key    cipher.PubKey
	client *Client
	In     chan []byte
	Out    chan []byte
}

func NewClientConnection(key cipher.PubKey, client *Client) *ClientConnection {
	return &ClientConnection{Key: key, client: client, In: make(chan []byte), Out: make(chan []byte)}
}

func (c *ClientConnection) WriteLoop() error {
	for {
		select {
		case d := <-c.Out:
			err := c.client.conn.WriteSlice(c.Key[:], d)
			if err != nil {
				return err
			}
		}
	}
}

// func(conn *ClientConnection, data []byte) bool
//
// return true for save this conn in factory so can use conn.Out for resp something
// otherwise conn.Out can not be used, because no receiver goroutine exists
type IncomingCallbackType func(conn *ClientConnection, data []byte) bool

type ClientConnectionFactory struct {
	client           *Client
	incomingCallback IncomingCallbackType
	fieldsMutex      *sync.RWMutex

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

func (factory *ClientConnectionFactory) SetIncomingCallback(fn IncomingCallbackType) {
	factory.fieldsMutex.Lock()
	defer factory.fieldsMutex.Unlock()
	factory.incomingCallback = fn
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
			key := cipher.NewPubKey(d[:msg.MSG_PUBKEY_SIZE])
			factory.connectionsMutex.RLock()
			conn, ok := factory.connections[key.Hex()]
			factory.connectionsMutex.RUnlock()
			data := d[msg.MSG_PUBKEY_SIZE:]
			if !ok {
				log.Printf("conn not exists %s", key.Hex())
				connection := NewClientConnection(key, factory.client)
				if factory.incomingCallback != nil && factory.incomingCallback(connection, data) {
					factory.connectionsMutex.Lock()
					factory.connections[key.Hex()] = connection
					factory.connectionsMutex.Unlock()

					go connection.WriteLoop()
				}
				continue
			}
			conn.In <- data
		}
	}
}

func (factory *ClientConnectionFactory) Dial(key cipher.PubKey) *ClientConnection {
	connection := NewClientConnection(key, factory.client)

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
	return &ClientConnection{Key: key, client: factory.client}
}
