package factory

import (
	"sync"

	"encoding/json"

	"time"

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

	// custom msg callback
	CustomMsgHandler func(*Connection, []byte)

	// service discovery data will update forward to this parent factory connections
	ServiceDiscoveryParent *MessengerFactory

	serviceDiscovery

	fieldsMutex sync.RWMutex
}

func NewMessengerFactory() *MessengerFactory {
	return &MessengerFactory{regConnections: make(map[cipher.PubKey]*Connection), serviceDiscovery: newServiceDiscovery()}
}

func (f *MessengerFactory) Listen(address string) error {
	tcpFactory := factory.NewTCPFactory()
	f.fieldsMutex.Lock()
	f.factory = tcpFactory
	f.fieldsMutex.Unlock()
	tcpFactory.AcceptedCallback = f.acceptedCallback
	return tcpFactory.Listen(address)
}

func (f *MessengerFactory) acceptedCallback(connection *factory.Connection) {
	var err error
	conn := newConnection(connection, f)
	conn.SetContextLogger(conn.GetContextLogger().WithField("app", "messenger"))
	defer func() {
		if e := recover(); e != nil {
			conn.GetContextLogger().Errorf("acceptedCallback recover err %v", e)
		}
		if err != nil {
			conn.GetContextLogger().Errorf("acceptedCallback err %v", err)
		}
		f.unregister(conn.GetKey(), conn)
		f.discoveryUnregister(conn)
		conn.Close()
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
			opn := m[MSG_OP_BEGIN]
			op := getOP(int(opn))
			if op == nil {
				continue
			}
			var rb []byte
			if sop, ok := op.(simpleOP); ok {
				body := m[MSG_HEADER_END:]
				if len(body) > 0 {
					err = json.Unmarshal(m[MSG_HEADER_END:], sop)
					if err != nil {
						return
					}
				}
				var r resp
				r, err = sop.Execute(f, conn)
				if err != nil {
					return
				}
				if r != nil {
					rb, err = json.Marshal(r)
				}
			} else {
				rb, err = op.RawExecute(f, conn, m)
			}
			if err != nil {
				return
			}
			if rb != nil {
				err = conn.writeOPBytes(opn|RESP_PREFIX, rb)
				if err != nil {
					return
				}
			}
			putOP(int(opn), op)
		}
	}
}

func (f *MessengerFactory) register(key cipher.PubKey, connection *Connection) {
	f.regConnectionsMutex.Lock()
	defer f.regConnectionsMutex.Unlock()
	c, ok := f.regConnections[key]
	if ok {
		if c == connection {
			log.Debugf("reg %s %p already", key.Hex(), connection)
			return
		}
		log.Debugf("reg close %s %p for %p", key.Hex(), c, connection)
		c.Close()
	}
	f.regConnections[key] = connection
	log.Debugf("reg %s %p", key.Hex(), connection)
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
		log.Debugf("unreg %s %p", key.Hex(), c)
	} else {
		log.Debugf("unreg %s %p != new %p", key.Hex(), connection, c)
	}
}

func (f *MessengerFactory) Connect(address string) (conn *Connection, err error) {
	return f.ConnectWithConfig(address, nil)
}

func (f *MessengerFactory) ConnectWithConfig(address string, config *ConnConfig) (conn *Connection, err error) {
	f.fieldsMutex.Lock()
	if f.factory == nil {
		tcpFactory := factory.NewTCPFactory()
		f.factory = tcpFactory
	}
	f.fieldsMutex.Unlock()
	c, err := f.factory.Connect(address)
	if err != nil {
		if config != nil && config.Reconnect {
			go func() {
				time.Sleep(config.ReconnectWait)
				f.ConnectWithConfig(address, config)
			}()
		}
		return nil, err
	}
	conn = newClientConnection(c, f)
	conn.SetContextLogger(conn.GetContextLogger().WithField("app", "messenger"))
	err = conn.Reg()
	if config != nil {
		conn.findServiceNodesByKeysCallback = config.FindServiceNodesByKeysCallback
		conn.findServiceNodesByAttributesCallback = config.FindServiceNodesByAttributesCallback
		if config.OnConnected != nil {
			config.OnConnected(conn)
		}
		if config.Reconnect {
			go func() {
				conn.WaitForDisconnected()
				time.Sleep(config.ReconnectWait)
				f.ConnectWithConfig(address, config)
			}()
		}
	}
	return
}

func (f *MessengerFactory) Close() error {
	return f.factory.Close()
}

func (f *MessengerFactory) ForEachConn(fn func(connection *Connection)) {
	f.factory.ForEachConn(func(conn *factory.Connection) {
		real := conn.RealObject
		if real == nil {
			return
		}
		c, ok := real.(*Connection)
		if !ok {
			return
		}
		fn(c)
	})
}

func (f *MessengerFactory) discoveryRegister(conn *Connection, ns *NodeServices) {
	f.serviceDiscovery.register(conn, ns)
	if f.ServiceDiscoveryParent != nil {
		nodeServices := f.pack()
		f.ServiceDiscoveryParent.ForEachConn(func(connection *Connection) {
			connection.UpdateServices(nodeServices)
		})
	}
}

func (f *MessengerFactory) discoveryUnregister(conn *Connection) {
	f.serviceDiscovery.unregister(conn)
	if f.ServiceDiscoveryParent != nil {
		nodeServices := f.pack()
		f.ServiceDiscoveryParent.ForEachConn(func(connection *Connection) {
			connection.UpdateServices(nodeServices)
		})
	}
}
