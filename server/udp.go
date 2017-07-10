package server

import (
	"github.com/skycoin/net/conn"
	"github.com/skycoin/net/msg"
	"net"
	"encoding/binary"
)

type ServerUDPConn struct {
	conn.UDPConn
}

func NewServerUDPConn(c *net.UDPConn) *ServerUDPConn {
	return &ServerUDPConn{UDPConn: conn.UDPConn{UdpConn: c}}
}

func (c *ServerUDPConn) ReadLoop(fn func(c *net.UDPConn, addr *net.UDPAddr) *conn.UDPConn) (err error) {
	defer func() {
		if err != nil {
			c.SetStatusToError(err)
		}
		c.Close()
	}()
	for {
		maxBuf := make([]byte, conn.MAX_UDP_PACKAGE_SIZE)
		n, addr, err := c.UdpConn.ReadFromUDP(maxBuf)
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					cc := fn(c.UdpConn, addr)
					c.CTXLogger.Debug("close in")
					close(cc.In)
					continue
				}
			}
			return err
		}
		maxBuf = maxBuf[:n]
		cc := fn(c.UdpConn, addr)

		t := maxBuf[msg.MSG_TYPE_BEGIN]
		switch t {
		case msg.TYPE_ACK:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.DelMsg(seq)
			c.UpdateLastAck(seq)
		case msg.TYPE_PING:
			c.CTXLogger.Debug("recv ping")
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
			func() {
				defer func() {
					if e := recover(); e != nil {
						c.CTXLogger.Debug(e)
					}
					cc.Close()
				}()
				cc.In <- maxBuf[msg.MSG_HEADER_END:]
			}()
		}

		cc.UpdateLastTime()
	}
	return nil
}
