package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/skycoin/net/conn"
	"github.com/skycoin/net/msg"
)

type ClientUDPConn struct {
	conn.UDPConn
}

func NewClientUDPConn(c *net.UDPConn) *ClientUDPConn {
	return &ClientUDPConn{
		UDPConn: conn.UDPConn{
			UdpConn:          c,
			ConnCommonFields: conn.NewConnCommonFileds(),
			AckChan:          make(chan struct{}, 1),
		},
	}
}

func (c *ClientUDPConn) ReadLoop() (err error) {
	defer func() {
		if e := recover(); e != nil {
			c.CTXLogger.Debug(e)
			err = fmt.Errorf("readloop panic err:%v", e)
		}
		if err != nil {
			c.SetStatusToError(err)
		}
		c.Close()
	}()
	for {
		maxBuf := make([]byte, conn.MAX_UDP_PACKAGE_SIZE)
		n, err := c.UdpConn.Read(maxBuf)
		if err != nil {
			return err
		}
		maxBuf = maxBuf[:n]

		t := maxBuf[msg.MSG_TYPE_BEGIN]
		switch t {
		case msg.TYPE_PONG:
		case msg.TYPE_ACK:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.CTXLogger.Debugf("recv ack %d", seq)
			if c.DelMsg(seq) {
				c.AckChan <- struct{}{}
			}
			c.UpdateLastAck(seq)
		case msg.TYPE_NORMAL:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			ok, err := c.Ack(seq)
			if err != nil {
				return err
			}
			if !ok {
				c.CTXLogger.Debugf("Ack failed, %x", maxBuf)
				continue
			}
			c.In <- maxBuf[msg.MSG_HEADER_END:]
		default:
			c.CTXLogger.Debugf("not implemented msg type %d", t)
			return fmt.Errorf("not implemented msg type %d", t)
		}
	}
}

const (
	TICK_PERIOD = 60
)

func (c *ClientUDPConn) ping() error {
	return c.WriteBytes(msg.GenPingMsg())
}

func (c *ClientUDPConn) WriteLoop() (err error) {
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
			c.CTXLogger.Debug("Ping out")
			err := c.ping()
			if err != nil {
				return err
			}
		case m, ok := <-c.Out:
			if !ok {
				c.CTXLogger.Debug("udp conn closed")
				return nil
			}
			//c.CTXLogger.Debugf("msg out %x", m)
			err := c.Write(m)
			if err != nil {
				c.CTXLogger.Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *ClientUDPConn) WriteBytes(bytes []byte) error {
	_, err := c.UdpConn.Write(bytes)
	return err
}
