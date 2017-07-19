package factory

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type MessengerFactory struct {
	factory             factory.Factory
	regConnections      map[cipher.PubKey]*Connection
	regConnectionsMutex sync.RWMutex
	CustomMsgHandler    func(*Connection, []byte)

	services      map[string]*ConnectionList
	servicesMutex sync.RWMutex
}

func NewMessengerFactory() *MessengerFactory {
	return &MessengerFactory{regConnections: make(map[cipher.PubKey]*Connection), services: make(map[string]*ConnectionList)}
}

func (f *MessengerFactory) Listen(address string) error {
	tcpFactory := factory.NewTCPFactory()
	f.factory = tcpFactory
	tcpFactory.AcceptedCallback = f.acceptedCallback
	return tcpFactory.Listen(address)
}

var EMPTY_KEY = cipher.PubKey{}

func (f *MessengerFactory) acceptedCallback(connection *factory.Connection) {
	go func() {
		conn := NewConnection(connection)
		conn.SetContextLogger(conn.GetContextLogger().WithField("app", "messenger"))
		defer func() {
			if err := recover(); err != nil {
				conn.GetContextLogger().Errorf("acceptedCallback err %v", err)
			}
			f.unregister(conn.GetKey(), conn)
			f.removeService(conn)
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
					if conn.GetKey() != EMPTY_KEY {
						conn.GetContextLogger().Infof("reg %s already", conn.key.Hex())
						continue
					}
					key, _ := cipher.GenerateKeyPair()
					conn.SetKey(key)
					conn.SetContextLogger(conn.GetContextLogger().WithField("pubkey", key.Hex()))
					f.register(key, conn)
					err := conn.Write(GenRegRespMsg(key))
					if err != nil {
						conn.GetContextLogger().Errorf("resp reg key %s err %v", key.Hex(), err)
						conn.Close()
					}
				case OP_SEND:
					if len(m) < SEND_MSG_TO_PUBLIC_KEY_END {
						return
					}
					key := cipher.NewPubKey(m[SEND_MSG_TO_PUBLIC_KEY_BEGIN:SEND_MSG_TO_PUBLIC_KEY_END])
					f.regConnectionsMutex.RLock()
					c, ok := f.regConnections[key]
					f.regConnectionsMutex.RUnlock()
					if !ok {
						conn.GetContextLogger().Infof("key %s not found", key.Hex())
						continue
					}
					err := c.Write(m)
					if err != nil {
						conn.GetContextLogger().Errorf("forward to key %s err %v", key.Hex(), err)
						c.GetContextLogger().Errorf("write %x err %v", m, err)
						c.Close()
					}
				case OP_CUSTOM:
					if f.CustomMsgHandler != nil {
						f.CustomMsgHandler(conn, m[MSG_HEADER_END:])
					}
				case OP_OFFER_SERVICE:
					service := string(m[MSG_HEADER_END:])
					if len(service) > 0 {
						f.addService(service, conn)
					}
				default:
					conn.GetContextLogger().Errorf("not implemented op %d", op)
				}
			}
		}
	}()
}

func (f *MessengerFactory) register(key cipher.PubKey, connection *Connection) {
	f.regConnectionsMutex.Lock()
	defer f.regConnectionsMutex.Unlock()
	c, ok := f.regConnections[key]
	if ok {
		if c == connection {
			log.Printf("reg %s %p already", key.Hex(), connection)
			return
		}
		log.Printf("reg close %s %p for %p", key.Hex(), c, connection)
		c.Close()
	}
	f.regConnections[key] = connection
	log.Printf("reg %s %p", key.Hex(), connection)
}

func (f *MessengerFactory) GetConnection(key cipher.PubKey) (c *Connection, ok bool) {
	f.regConnectionsMutex.RLock()
	c, ok = f.regConnections[key]
	f.regConnectionsMutex.RUnlock()
	return
}

func (f *MessengerFactory) unregister(key cipher.PubKey, connection *Connection) {
	f.regConnectionsMutex.Lock()
	defer f.regConnectionsMutex.Unlock()
	c, ok := f.regConnections[key]
	if ok && c == connection {
		delete(f.regConnections, key)
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
	conn = NewConnection(c)
	conn.SetContextLogger(conn.GetContextLogger().WithField("app", "messenger"))
	err = conn.Reg()
	return
}

func (f *MessengerFactory) Close() error {
	return f.factory.Close()
}

func (f *MessengerFactory) addService(service string, conn *Connection) {
	if len(conn.GetService()) < 1 {
		return
	}
	conn.SetService(service)
	f.servicesMutex.Lock()
	defer f.servicesMutex.Unlock()
	list, ok := f.services[service]
	if ok {
		list.PushBack(conn)
	} else {
		l := NewConnectionList()
		l.PushBack(conn)
		f.services[service] = l
	}
	conn.GetContextLogger().Debugf("added service %s now %v", service, f.services)
}

func (f *MessengerFactory) removeService(conn *Connection) {
	service := conn.GetService()
	if len(service) < 1 {
		return
	}
	f.servicesMutex.Lock()
	defer f.servicesMutex.Unlock()
	list, ok := f.services[service]
	if ok {
		if list.Remove(conn) < 1 {
			delete(f.services, service)
		}
		conn.GetContextLogger().Debugf("removed service %s now %v", service, f.services)
	} else {
		conn.GetContextLogger().Debugf("remove service cannot find list of %s now %v", service, f.services)
	}
}
