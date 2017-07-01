package client

import (
	"time"
	"github.com/skycoin/net/conn"
	"net"
	"log"
)

type ClientTCPConn struct {
	conn.TCPConn
}

func NewClientTCPConn(c *net.TCPConn) *ClientTCPConn {
	return &ClientTCPConn{conn.TCPConn{TcpConn: c, In: make(chan []byte), Out: make(chan []byte), PendingMap: conn.PendingMap{Pending: make(map[uint32]interface{})}}}
}

func (c *ClientTCPConn) WriteLoop() error {
	ticker := time.NewTicker(time.Second * TICK_PERIOD)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			log.Println("ping out")
			err := c.Ping()
			if err != nil {
				return err
			}
		case m, ok := <-c.Out:
			if !ok {
				log.Println("conn closed")
				return nil
			}
			log.Printf("msg Out %x", m)
			err := c.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return err
			}
		}
	}
}

