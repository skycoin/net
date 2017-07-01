package conn

import (
	"encoding/binary"
	"github.com/skycoin/net/msg"
	"net"
	"log"
	"time"
	"sync"
	"sync/atomic"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	MAX_UDP_PACKAGE_SIZE = 1024
)

type UDPConn struct {
	UdpConn *net.UDPConn
	addr    *net.UDPAddr
	In      chan []byte
	Out     chan []byte

	seq          uint32
	PendingMap


	lastTime    int64
	closed      bool
	pubkey      cipher.PubKey
	fieldsMutex sync.RWMutex
}

func NewUDPConn(c *net.UDPConn, addr *net.UDPAddr) *UDPConn {
	return &UDPConn{UdpConn: c, addr: addr, lastTime: time.Now().Unix(), In: make(chan []byte), Out: make(chan []byte), PendingMap: PendingMap{Pending: make(map[uint32]interface{})}}
}

func (c *UDPConn) ReadLoop() error {
	panic("UDPConn unimplemented ReadLoop")
}

func (c *UDPConn) WriteSlice(bytes ...[]byte) error {
	panic("UDPConn unimplemented WriteSlice")
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
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_NORMAL, new, bytes)
	c.AddMsgToPendingMap(new, m)
	return c.WriteBytes(m.Bytes())
}

func (c *UDPConn) GetPublicKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.pubkey
}

func (c *UDPConn) SetPublicKey(key cipher.PubKey) {
	c.fieldsMutex.Lock()
	c.pubkey = key
	c.fieldsMutex.Unlock()
}

func (c *UDPConn) WriteBytes(bytes []byte) error {
	_, err := c.UdpConn.WriteToUDP(bytes, c.addr)
	return err
}

func (c *UDPConn) Ack(seq uint32) error {
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	return c.WriteBytes(resp)
}

func (c *UDPConn) GetChanOut() chan<- []byte {
	return c.Out
}

func (c *UDPConn) GetChanIn() <-chan []byte {
	return c.In
}

func (c *UDPConn) IsClosed() bool {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.closed
}

func (c *UDPConn) SendReg(key cipher.PubKey) error {
	panic("UDPConn unimplemented SendReg")
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

func (c *UDPConn) Close() {
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

func (c *UDPConn) GetNextSeq() uint32 {
	return atomic.AddUint32(&c.seq, 1)
}