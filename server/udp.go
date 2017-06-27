package server

import (
	"github.com/skycoin/net/conn"
	"github.com/skycoin/net/msg"
	"log"
	"net"
	"encoding/binary"
	"github.com/skycoin/skycoin/src/cipher"
)

type ServerUDPConn struct {
	conn.UDPConn
	factory *ConnectionFactory
}

func NewServerUDPConn(c *net.UDPConn, factory *ConnectionFactory) *ServerUDPConn {
	return &ServerUDPConn{UDPConn:conn.UDPConn{UdpConn:c}, factory:factory}
}

func (c *ServerUDPConn) ReadLoop() error {
	for {
		maxBuf := make([]byte, conn.MAX_UDP_PACKAGE_SIZE)
		n, addr, err := c.UdpConn.ReadFromUDP(maxBuf)
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					cc := c.factory.GetOrCreateUDPConn(c.UdpConn, addr)
					log.Println("close in")
					close(cc.In)
					continue
				}
			}
			return err
		}
		maxBuf = maxBuf[:n]
		cc := c.factory.GetOrCreateUDPConn(c.UdpConn, addr)

		t := maxBuf[msg.MSG_TYPE_BEGIN]
		switch t {
		case msg.TYPE_ACK:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.DelMsgToPendingMap(seq)
		case msg.TYPE_PING:
			log.Println("recv ping")
			err = cc.WriteBytes([]byte{msg.TYPE_PONG})
			if err != nil {
				return err
			}
		case msg.TYPE_NORMAL:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			err = cc.Ack(seq)
			if err != nil {
				return err
			}
			cc.In <- maxBuf[msg.MSG_HEADER_END:]
		case msg.TYPE_REG:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			err = cc.Ack(seq)
			if err != nil {
				return err
			}
			key := cipher.NewPubKey(maxBuf[msg.MSG_HEADER_END:])
			cc.SetPublicKey(key)
			c.factory.Register(key.Hex(), cc)
		}

		cc.UpdateLastTime()
	}
	return nil
}

