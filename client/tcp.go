package client

import (
	"time"
	"github.com/skycoin/net/conn"
	"net"
	"log"
)

type ClientTCPConn struct {
	*conn.TCPConn
}

func NewClientTCPConn(c net.Conn) *ClientTCPConn {
	return &ClientTCPConn{&conn.TCPConn{TcpConn: c, In: make(chan []byte), Out: make(chan []byte), ConnCommonFields:conn.NewConnCommonFileds()}}
}

func (c *ClientTCPConn) WriteLoop() (err error) {
	ticker := time.NewTicker(time.Second * TICK_PERIOD)
	defer func() {
		ticker.Stop()
		if err != nil {
			c.SetStatusToError(err)
		}
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

