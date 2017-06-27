package server

import (
	"net"
	"time"
	"sync"
	"log"
	"github.com/skycoin/net/conn"
)

var (
	once = &sync.Once{}
)

type ConnectionFactory struct {
	connMapMutex sync.RWMutex
	connMap      map[string]conn.Connection

	udpConnMapMutex sync.RWMutex
	udpConnMap      map[string]*conn.UDPConn

	ConnHandler func(connection conn.Connection)
}

func NewFactory() *ConnectionFactory {
	return &ConnectionFactory{connMap: make(map[string]conn.Connection),
		udpConnMap: make(map[string]*conn.UDPConn)}
}

func (factory *ConnectionFactory) CreateTCPConn(c *net.TCPConn) conn.Connection {
	cc := NewServerTCPConn(c, factory)
	go factory.ConnHandler(cc)
	go func() {
		cc.WriteLoop()
		factory.UnRegister(cc.GetPublicKey().Hex(), cc)
	}()
	return cc
}

func (factory *ConnectionFactory) GetOrCreateUDPConn(c *net.UDPConn, addr *net.UDPAddr) *conn.UDPConn {
	log.Println(addr.String())
	factory.udpConnMapMutex.RLock()
	if cc, ok := factory.udpConnMap[addr.String()]; ok {
		factory.udpConnMapMutex.RUnlock()
		return cc
	}
	factory.udpConnMapMutex.RUnlock()

	log.Println("new udp")
	once.Do(func() {
		go factory.GC()
	})

	udpConn := conn.NewUDPConn(c, addr)
	factory.udpConnMapMutex.Lock()
	factory.udpConnMap[addr.String()] = udpConn
	factory.udpConnMapMutex.Unlock()

	go factory.ConnHandler(udpConn)
	go func() {
		udpConn.WriteLoop()
		factory.UnRegister(udpConn.GetPublicKey().Hex(), udpConn)
	}()
	return udpConn
}

func (factory *ConnectionFactory) Register(pubkey string, conn conn.Connection) {
	factory.connMapMutex.Lock()
	defer factory.connMapMutex.Unlock()
	factory.connMap[pubkey] = conn
	log.Printf("regsiter %s %v", pubkey, conn)
}

func (factory *ConnectionFactory) UnRegister(pubkey string, conn conn.Connection) {
	factory.connMapMutex.Lock()
	defer factory.connMapMutex.Unlock()
	if c, ok := factory.connMap[pubkey]; ok && c == conn {
		delete(factory.connMap, pubkey)
		log.Printf("UnRegister %s %v ok", pubkey, conn)
	}
	log.Printf("UnRegister %s %v", pubkey, conn)
}

func (factory *ConnectionFactory) GetConn(pubkey string) conn.Connection {
	factory.connMapMutex.RLock()
	defer factory.connMapMutex.RUnlock()
	if c, ok := factory.connMap[pubkey]; ok {
		return c
	}
	return nil
}

const UDP_GC_PERIOD = 90

func (factory *ConnectionFactory) GC() {
	ticker := time.NewTicker(time.Second * UDP_GC_PERIOD)
	for {
		select {
		case <-ticker.C:
			nowUnix := time.Now().Unix()
			closed := []string{}
			factory.udpConnMapMutex.RLock()
			for k, udp := range factory.udpConnMap {
				if nowUnix-udp.GetLastTime() >= UDP_GC_PERIOD {
					udp.Close()
					closed = append(closed, k)
				}
			}
			factory.udpConnMapMutex.RUnlock()
			if len(closed) < 1 {
				continue
			}
			factory.udpConnMapMutex.Lock()
			for _, u := range closed {
				delete(factory.udpConnMap, u)
			}
			factory.udpConnMapMutex.Unlock()
		}
	}
}
