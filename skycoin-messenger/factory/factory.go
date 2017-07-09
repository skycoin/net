package factory

import (
	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
	"log"
)

type MessengerFactory struct {
	factory factory.Factory
	regConnections      map[string]*Connection
	regConnectionsMutex sync.RWMutex
}

func NewMessengerFactory() *MessengerFactory {
	return &MessengerFactory{regConnections:make(map[string]*Connection)}
}

func (f *MessengerFactory) Listen(address string) error {
	tcpFactory := factory.NewTCPFactory()
	f.factory = tcpFactory
	tcpFactory.AcceptedCallback = f.acceptedCallback
	return tcpFactory.Listen(address)
}

func (f *MessengerFactory) acceptedCallback(connection *factory.Connection) {
	go func() {
		conn := &Connection{Connection:connection}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("acceptedCallback err %v", err)
			}
			f.unregister(conn.GetKey(), conn)
		}()
		for {
			select {
			case m, ok := <-conn.GetChanIn():
				if !ok {
					return
				}
				if len(m) < MSG_HEADER_END {
					return
				}
				op := m[MSG_OP_BEGIN]
				switch op {
				case OP_REG:
					if len(m) < MSG_PUBLIC_KEY_END {
						return
					}
					key := cipher.NewPubKey(m[MSG_PUBLIC_KEY_BEGIN:MSG_PUBLIC_KEY_END])
					conn.SetKey(key)
					f.register(key, conn)
				case OP_SEND:
					if len(m) < MSG_TO_PUBLIC_KEY_END {
						return
					}
					key := cipher.NewPubKey(m[MSG_TO_PUBLIC_KEY_BEGIN:MSG_TO_PUBLIC_KEY_END])
					f.regConnectionsMutex.RLock()
					c, ok := f.regConnections[key.Hex()]
					f.regConnectionsMutex.RUnlock()
					if !ok {
						log.Printf("key %s not found", key.Hex())
						continue
					}
					err := c.Write(m)
					if err != nil {
						log.Printf("forward err %v", err)
						c.Close()
					}
				}
			}
		}
	}()
}

func (f *MessengerFactory) register(key cipher.PubKey, connection *Connection) {
	f.regConnectionsMutex.Lock()
	defer f.regConnectionsMutex.Unlock()
	c, ok := f.regConnections[key.Hex()]
	if ok {
		if c == connection {
			log.Printf("reg %s %p already", key.Hex(), connection)
			return
		}
		log.Printf("reg close %s %p for %p", key.Hex(), c, connection)
		c.Close()
	}
	f.regConnections[key.Hex()] = connection
	log.Printf("reg %s %p", key.Hex(), connection)
}

func (f *MessengerFactory) unregister(key cipher.PubKey, connection *Connection) {
	f.regConnectionsMutex.Lock()
	defer f.regConnectionsMutex.Unlock()
	c, ok := f.regConnections[key.Hex()]
	if ok && c == connection{
		delete(f.regConnections, key.Hex())
		log.Printf("unreg %s %p", key.Hex(), c)
	} else {
		log.Printf("unreg %s %p != new %p", key.Hex(), connection, c)
	}
}

func (f *MessengerFactory) Connect(address string) (conn *Connection, err error) {
	tcpFactory := factory.NewTCPFactory()
	c, err := tcpFactory.Connect(address)
	if err != nil {
		return nil, err
	}
	conn = &Connection{Connection:c}
	return
}

func (f *MessengerFactory) Close() error {
	return f.factory.Close()
}
