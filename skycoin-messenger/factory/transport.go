package factory

import (
	"sync"

	"net"

	"io"

	"sync/atomic"

	"encoding/binary"

	log "github.com/sirupsen/logrus"
	cn "github.com/skycoin/net/conn"
	"github.com/skycoin/skycoin/src/cipher"
)

type transport struct {
	creator *MessengerFactory
	// node
	factory *MessengerFactory
	// conn between nodes
	conn *Connection
	// app
	appNet net.Listener

	fromNode, toNode cipher.PubKey
	fromApp, toApp   cipher.PubKey

	conns      map[uint32]net.Conn
	connsMutex sync.RWMutex

	fieldsMutex sync.RWMutex
}

func NewTransport(creator *MessengerFactory, fromNode, toNode, fromApp, toApp cipher.PubKey) *transport {
	t := &transport{
		creator:  creator,
		fromNode: fromNode,
		toNode:   toNode,
		fromApp:  fromApp,
		toApp:    toApp,
		factory:  NewMessengerFactory(),
		conns:    make(map[uint32]net.Conn),
	}
	return t
}

func (t *transport) ListenAndConnect(address string) (conn *Connection, err error) {
	conn, err = t.factory.connectUDPWithConfig(address, &ConnConfig{
		OnConnected: func(connection *Connection) {
			connection.Reg()
		},
		Creator: t.creator,
	})
	return
}

func (t *transport) Connect(address, appAddress string) (err error) {
	conn, err := t.factory.connectUDPWithConfig(address, &ConnConfig{
		OnConnected: func(connection *Connection) {
			connection.writeOP(OP_BUILD_APP_CONN_OK,
				&buildConnResp{
					FromNode: t.fromNode,
					Node:     t.toNode,
					FromApp:  t.fromApp,
					App:      t.toApp,
				})
		},
		Creator: t.creator,
	})
	if err != nil {
		return
	}
	t.fieldsMutex.Lock()
	t.conn = conn
	t.fieldsMutex.Unlock()

	go func() {
		var err error
		for {
			select {
			case m, ok := <-conn.GetChanIn():
				if !ok {
					log.Debugf("node conn read err %v", err)
					return
				}
				//log.Debugf("read from node udp %x", m)
				id := binary.BigEndian.Uint32(m[PKG_HEADER_ID_BEGIN:PKG_HEADER_ID_END])
				t.connsMutex.RLock()
				appConn, ok := t.conns[id]
				t.connsMutex.RUnlock()
				if !ok {
					appConn, err = net.Dial("tcp", appAddress)
					if err != nil {
						log.Debugf("app conn dial err %v", err)
						return
					}
					go t.appReadLoop(id, appConn, conn, false)
				}
				body := m[PKG_HEADER_END:]
				if len(body) < 1 {
					continue
				}
				err = writeAll(appConn, body)
				log.Debugf("send to tcp")
				if err != nil {
					log.Debugf("app conn write err %v", err)
					return
				}
			}
		}
	}()

	return
}

func (t *transport) appReadLoop(id uint32, appConn net.Conn, conn *Connection, create bool) {
	buf := make([]byte, cn.MAX_UDP_PACKAGE_SIZE-100)
	binary.BigEndian.PutUint32(buf, id)
	if create {
		conn.GetChanOut() <- buf[:PKG_HEADER_END]
	}
	t.connsMutex.Lock()
	t.conns[id] = appConn
	t.connsMutex.Unlock()
	defer func() {
		t.connsMutex.Lock()
		delete(t.conns, id)
		t.connsMutex.Unlock()
	}()
	for {
		n, err := appConn.Read(buf[PKG_HEADER_END:])
		if err != nil {
			log.Debugf("app conn read err %v, %d", err, n)
			return
		}
		pkg := make([]byte, PKG_HEADER_END+n)
		copy(pkg, buf[:PKG_HEADER_END+n])
		//log.Debugf("tcp conn read %x", pkg)
		conn.GetChanOut() <- pkg
		log.Debugf("send to node udp")
	}
}

func (t *transport) setUDPConn(conn *Connection) {
	t.fieldsMutex.Lock()
	t.conn = conn
	t.fieldsMutex.Unlock()
}

func (t *transport) ListenForApp(address string, fn func()) (err error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	t.fieldsMutex.Lock()
	t.appNet = ln
	t.fieldsMutex.Unlock()

	fn()

	go t.accept()
	return
}

const (
	PKG_HEADER_ID_SIZE = 4

	PKG_HEADER_BEGIN    = 0
	PKG_HEADER_ID_BEGIN
	PKG_HEADER_ID_END   = PKG_HEADER_ID_BEGIN + PKG_HEADER_ID_SIZE
	PKG_HEADER_END
)

func (t *transport) accept() {
	t.fieldsMutex.RLock()
	tConn := t.conn
	t.fieldsMutex.RUnlock()

	go func() {
		var err error
		for {
			select {
			case m, ok := <-tConn.GetChanIn():
				if !ok {
					log.Debug("transport closed")
					return
				}
				//log.Debugf("ListenForApp from node B udp %x", m)
				id := binary.BigEndian.Uint32(m[PKG_HEADER_ID_BEGIN:PKG_HEADER_ID_END])
				t.connsMutex.RLock()
				conn, ok := t.conns[id]
				t.connsMutex.RUnlock()
				if !ok {
					log.Debugf("node a tcp conn %d not found", id)
					continue
				}
				body := m[PKG_HEADER_END:]
				if len(body) < 1 {
					continue
				}
				err = writeAll(conn, body)
				log.Debugf("ListenForApp write to app")
				if err != nil {
					log.Debugf("ListenForApp write to app err %s", err)
					continue
				}
			}
		}
	}()

	var idSeq uint32
	for {
		conn, err := t.appNet.Accept()
		if err != nil {
			return
		}
		id := atomic.AddUint32(&idSeq, 1)
		go t.appReadLoop(id, conn, tConn, true)
	}
}

func (t *transport) Close() {
	t.fieldsMutex.Lock()
	defer t.fieldsMutex.Unlock()

	if t.factory == nil {
		return
	}

	if t.appNet != nil {
		t.appNet.Close()
		t.appNet = nil
	}
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	t.factory.Close()
	t.factory = nil

	t.connsMutex.RLock()
	for _, v := range t.conns {
		v.Close()
	}
	t.connsMutex.RUnlock()
}

func writeAll(conn io.Writer, m []byte) error {
	for i := 0; i < len(m); {
		n, err := conn.Write(m[i:])
		if err != nil {
			return err
		}
		i += n
	}
	return nil
}
