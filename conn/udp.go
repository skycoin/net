package conn

import (
	"encoding/binary"
	"net"
	"sync/atomic"
	"time"

	"github.com/skycoin/net/msg"
	"sync"
	"hash/crc32"
)

const (
	MAX_UDP_PACKAGE_SIZE = 1200
)

type UDPConn struct {
	ConnCommonFields
	*UDPPendingMap
	StreamQueue
	UdpConn *net.UDPConn
	addr    *net.UDPAddr

	lastTime int64

	// write loop with ping
	SendPing bool
	wmx      sync.Mutex
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
	if c.SendPing {
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
			err := c.Ping()
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
	c.wmx.Lock()
	defer c.wmx.Unlock()
	s := c.GetNextSeq()
	m := msg.New(msg.TYPE_NORMAL, s, bytes)
	c.AddMsg(s, m)
	c.GetContextLogger().Debugf("Write msg seq %d", s)
	err = c.WriteBytes(m.PkgBytes())
	return
}

func (c *UDPConn) WriteBytes(bytes []byte) error {
	//c.CTXLogger.Debugf("write %x", bytes)
	c.WriteMutex.Lock()
	defer c.WriteMutex.Unlock()
	_, err := c.UdpConn.WriteToUDP(bytes, c.addr)
	return err
}

func (c *UDPConn) Ack(seq uint32) error {
	p := make([]byte, msg.MSG_SEQ_END+msg.PKG_HEADER_SIZE)
	m := p[msg.PKG_HEADER_SIZE:]
	m[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(m[msg.MSG_SEQ_BEGIN:], seq)
	checksum := crc32.ChecksumIEEE(m)
	binary.BigEndian.PutUint32(p[msg.PKG_CRC32_BEGIN:], checksum)
	return c.WriteBytes(p)
}

func (c *UDPConn) Ping() error {
	p := make([]byte, msg.PING_MSG_HEADER_SIZE+msg.PKG_HEADER_SIZE)
	m := p[msg.PKG_HEADER_SIZE:]
	m[msg.PING_MSG_TYPE_BEGIN] = msg.TYPE_PING
	binary.BigEndian.PutUint64(m[msg.PING_MSG_TIME_BEGIN:], msg.UnixMillisecond())
	checksum := crc32.ChecksumIEEE(m)
	binary.BigEndian.PutUint32(p[msg.PKG_CRC32_BEGIN:], checksum)
	return c.WriteBytes(p)
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
