package conn

import (
	"encoding/binary"
	"github.com/skycoin/net/msg"
	"net"
	"log"
	"time"
	"sync"
)

const (
	MAX_UDP_PACKAGE_SIZE = 1024
)

type UDPConn struct {
	udpConn *net.UDPConn
	addr    *net.UDPAddr
	In      chan []byte
	Out     chan []byte

	seq     uint32
	pending map[uint32]*msg.Message

	lastTime    int64
	closed      bool
	fieldsMutex *sync.RWMutex
}

type UDPServerConn struct {
	UDPConn
	factory *ConnectionFactory
}

func NewUDPConn(c *net.UDPConn, addr *net.UDPAddr) *UDPConn {
	return &UDPConn{udpConn: c, addr: addr, lastTime: time.Now().Unix(), fieldsMutex: new(sync.RWMutex), In: make(chan []byte), Out: make(chan []byte), pending: make(map[uint32]*msg.Message)}
}

func NewUDPServerConn(c *net.UDPConn, factory *ConnectionFactory) *UDPServerConn {
	sc := &UDPServerConn{}
	sc.factory = factory
	sc.udpConn = c
	return sc
}

func (c *UDPServerConn) ReadLoop() error {
	for {
		maxBuf := make([]byte, MAX_UDP_PACKAGE_SIZE)
		n, addr, err := c.udpConn.ReadFromUDP(maxBuf)
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					cc := c.factory.GetOrCreateUDPConn(c.udpConn, addr)
					log.Println("close in")
					close(cc.In)
					continue
				}
			}
			return err
		}
		maxBuf = maxBuf[:n]
		cc := c.factory.GetOrCreateUDPConn(c.udpConn, addr)

		switch maxBuf[msg.MSG_TYPE_BEGIN] {
		case msg.TYPE_ACK:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			delete(c.pending, seq)
		case msg.TYPE_PING:
			err = c.writeBytes([]byte{msg.TYPE_PONG})
			if err != nil {
				return err
			}
		case msg.TYPE_NORMAL:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			err = cc.ack(seq)
			if err != nil {
				return err
			}
			cc.In <- maxBuf[msg.MSG_HEADER_END:]
		}

		cc.fieldsMutex.Lock()
		cc.lastTime = time.Now().Unix()
		cc.fieldsMutex.Unlock()
	}
	return nil
}

func (c *UDPConn) WriteLoop() error {
	for {
		select {
		case m, ok := <-c.Out:
			if !ok {
				log.Println("udp conn closed")
				return nil
			}
			log.Printf("msg out %x", m)
			err := c.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPConn) Write(bytes []byte) error {
	c.seq++
	m := msg.New(msg.TYPE_NORMAL, c.seq, bytes)
	c.pending[c.seq] = m
	return c.writeBytes(m.Bytes())
}

func (c *UDPConn) writeBytes(bytes []byte) error {
	_, err := c.udpConn.WriteToUDP(bytes, c.addr)
	return err
}

func (c *UDPConn) ack(seq uint32) error {
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	return c.writeBytes(resp)
}

func (c *UDPConn) GetChanOut() chan<- []byte {
	return c.Out
}

func (c *UDPConn) IsClosed() bool {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.closed
}

func (c *UDPConn) GetLastTime() int64 {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.lastTime
}

func (c *UDPConn) close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("closing closed udpconn")
		}
	}()
	c.fieldsMutex.Lock()
	c.closed = true
	c.fieldsMutex.Unlock()
	close(c.In)
	close(c.Out)
}

func (c *UDPConn) Ping() error {
	b := make([]byte, msg.MSG_TYPE_SIZE)
	b[msg.MSG_TYPE_BEGIN] = msg.TYPE_PING
	return c.writeBytes(b)
}

type UDPClientConn struct {
	UDPConn
}

func NewUDPClientConn(c *net.UDPConn) *UDPClientConn {
	cc := &UDPClientConn{}
	cc.udpConn = c
	cc.In = make(chan []byte)
	cc.Out = make(chan []byte)
	cc.pending = make(map[uint32]*msg.Message)
	return cc
}

func (c *UDPClientConn) ReadLoop() error {
	for {
		maxBuf := make([]byte, MAX_UDP_PACKAGE_SIZE)
		n, err := c.udpConn.Read(maxBuf)
		if err != nil {
			return err
		}
		maxBuf = maxBuf[:n]

		switch maxBuf[msg.MSG_TYPE_BEGIN] {
		case msg.TYPE_ACK:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			delete(c.pending, seq)
		case msg.TYPE_NORMAL:
			seq := binary.BigEndian.Uint32(maxBuf[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			err = c.ack(seq)
			if err != nil {
				return err
			}
			c.In <- maxBuf[msg.MSG_HEADER_END:]
		}
	}
	return nil
}

func (c *UDPClientConn) WriteLoop() error {
	for {
		select {
		case m, ok := <-c.Out:
			if !ok {
				log.Println("udp conn closed")
				return nil
			}
			log.Printf("msg out %x", m)
			err := c.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return err
			}
		}
	}
}

func (c *UDPClientConn) Write(bytes []byte) error {
	c.seq++
	m := msg.New(msg.TYPE_NORMAL, c.seq, bytes)
	c.pending[c.seq] = m
	return c.writeBytes(m.Bytes())
}

func (c *UDPClientConn) writeBytes(bytes []byte) error {
	_, err := c.udpConn.Write(bytes)
	return err
}

func (c *UDPClientConn) ack(seq uint32) error {
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	return c.writeBytes(resp)
}
