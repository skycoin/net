package conn

import (
	"net"
	"time"
	"sync"
	"log"
)

var (
	once = &sync.Once{}
)

type ConnectionFactory struct {
	udpConnsMutex *sync.RWMutex
	udpConns      map[string]*UDPConn

	TCPClientHandler func(connection *TCPConn)
	UDPClientHandler func(connection *UDPConn)
}

func NewFactory() *ConnectionFactory {
	return &ConnectionFactory{udpConnsMutex: new(sync.RWMutex), udpConns: make(map[string]*UDPConn)}
}

func (factory *ConnectionFactory) CreateTCPConn(c *net.TCPConn) *TCPConn {
	cc := NewTCPConn(c)
	go factory.TCPClientHandler(cc)
	return cc
}

func (factory *ConnectionFactory) GetOrCreateUDPConn(c *net.UDPConn, addr *net.UDPAddr) *UDPConn {
	log.Println(addr.String())
	factory.udpConnsMutex.RLock()
	if cc, ok := factory.udpConns[addr.String()]; ok {
		factory.udpConnsMutex.RUnlock()
		return cc
	}
	factory.udpConnsMutex.RUnlock()

	log.Println("new udp")
	once.Do(func() {
		go factory.GC()
	})

	udpConn := NewUDPConn(c, addr)
	factory.udpConnsMutex.Lock()
	factory.udpConns[addr.String()] = udpConn
	factory.udpConnsMutex.Unlock()

	go factory.UDPClientHandler(udpConn)
	return udpConn
}

const UDP_GC_PERIOD = 90

func (factory *ConnectionFactory) GC() {
	ticker := time.NewTicker(time.Second * UDP_GC_PERIOD)
	for {
		select {
		case <-ticker.C:
			nowUnix := time.Now().Unix()
			closed := []string{}
			factory.udpConnsMutex.RLock()
			for k, udp := range factory.udpConns {
				if nowUnix-udp.lastTime >= UDP_GC_PERIOD {
					udp.close()
					closed = append(closed, k)
				}
			}
			factory.udpConnsMutex.RUnlock()
			if len(closed) < 1 {
				continue
			}
			factory.udpConnsMutex.Lock()
			for _, u := range closed {
				delete(factory.udpConns, u)
			}
			factory.udpConnsMutex.Unlock()
		}
	}
}
