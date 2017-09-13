package conn

import (
	"encoding/binary"
	"net"
	"sync/atomic"
	"time"

	"github.com/skycoin/net/msg"
)

const (
	MAX_UDP_PACKAGE_SIZE = 1200
)

type UDPConn struct {
	ConnCommonFields
	*UDPPendingMap
	UdpConn *net.UDPConn
	addr    *net.UDPAddr

	lastTime int64
	ackSeq   uint32

	// write loop with ping
	Ping bool
}

// used for server spawn udp conn
func NewUDPConn(c *net.UDPConn, addr *net.UDPAddr) *UDPConn {
	return &UDPConn{
		UdpConn:          c,
		addr:             addr,
		lastTime:         time.Now().Unix(),
		ConnCommonFields: NewConnCommonFileds(),
		UDPPendingMap:    NewUDPPendingMap(),
	}
}

func (c *UDPConn) ReadLoop() error {
	return nil
}

func (c *UDPConn) WriteLoop() (err error) {
	if c.Ping {
		return c.writeLoopWithPing()
	} else {
		return c.writeLoop()
	}
}

func (c *UDPConn) writeLoop() (err error) {
	defer func() {
		if err != nil {
			c.SetStatusToError(err)
		}
	}()
	for {
		select {
		case m, ok := <-c.Out:
			if !ok {
				c.CTXLogger.Debug("udp conn closed")
				return nil
			}
			err := c.Write(m)
			if err != nil {
				c.CTXLogger.Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPConn) writeLoopWithPing() (err error) {
	ticker := time.NewTicker(time.Second * UDP_PING_TICK_PERIOD)
	defer func() {
		ticker.Stop()
		if err != nil {
			c.SetStatusToError(err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			err := c.WriteBytes(msg.GenPingMsg())
			if err != nil {
				return err
			}
		case m, ok := <-c.Out:
			if !ok {
				c.CTXLogger.Debug("udp conn closed")
				return nil
			}
			err := c.Write(m)
			if err != nil {
				c.CTXLogger.Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPConn) Write(bytes []byte) (err error) {
	s := c.GetNextSeq()
	m := msg.New(msg.TYPE_NORMAL, s, bytes)
	c.AddMsg(s, m)
	err = c.WriteBytes(m.Bytes())
	return
}

func (c *UDPConn) WriteBytes(bytes []byte) error {
	//c.CTXLogger.Debugf("write %x", bytes)
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	_, err := c.UdpConn.WriteToUDP(bytes, c.addr)
	return err
}

func (c *UDPConn) Ack(seq uint32) (ok bool, err error) {
	ok = atomic.CompareAndSwapUint32(&c.ackSeq, seq-1, seq)
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	err = c.WriteBytes(resp)
	if !ok {
		c.CTXLogger.Debugf("Ack now is %d try to ack %d %v %v", atomic.LoadUint32(&c.ackSeq), seq, ok, err)
	}
	return
}

func (c *UDPConn) GetChanOut() chan<- []byte {
	return c.Out
}

func (c *UDPConn) GetChanIn() <-chan []byte {
	return c.In
}

func (c *UDPConn) GetLastTime() int64 {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.lastTime
}

func (c *UDPConn) UpdateLastTime() {
	c.fieldsMutex.Lock()
	c.lastTime = time.Now().Unix()
	c.fieldsMutex.Unlock()
}

func (c *UDPConn) GetNextSeq() uint32 {
	return atomic.AddUint32(&c.seq, 1)
}

func (c *UDPConn) Close() {
	c.fieldsMutex.Lock()
	if c.UdpConn != nil {
		c.UdpConn.Close()
	}
	c.fieldsMutex.Unlock()
	c.ConnCommonFields.Close()
}

func (c *UDPConn) GetRemoteAddr() net.Addr {
	return c.addr
}
