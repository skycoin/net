package conn

import (
	"net"
	"github.com/skycoin/net/msg"
	"bufio"
	"io"
	"encoding/binary"
	"time"
	"log"
	"sync"
	"sync/atomic"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	TCP_READ_TIMEOUT = 90
)

type TCPConn struct {
	TcpConn *net.TCPConn
	In      chan []byte
	Out     chan []byte

	seq          uint32
	PendingMap

	closed      bool
	pubkey      cipher.PubKey
	fieldsMutex sync.RWMutex
}

func (c *TCPConn) ReadLoop() error {
	defer func() {
		c.Close()
	}()
	header := make([]byte, msg.MSG_HEADER_SIZE)
	reader := bufio.NewReader(c.TcpConn)

	for {
		t, err := reader.Peek(msg.MSG_TYPE_SIZE)
		if err != nil {
			return err
		}
		msg_t := t[msg.MSG_TYPE_BEGIN]
		switch msg_t {
		case msg.TYPE_ACK:
			_, err = io.ReadAtLeast(reader, header, msg.MSG_SEQ_END)
			if err != nil {
				return err
			}
			seq := binary.BigEndian.Uint32(header[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.DelMsgToPendingMap(seq)
			c.PendingMap.RLock()
			log.Printf("acked %d, Pending:%d, %v", seq, len(c.Pending), c.Pending)
			c.PendingMap.RUnlock()
		case msg.TYPE_PING:
			reader.Discard(msg.MSG_TYPE_SIZE)
			err = c.WriteBytes([]byte{msg.TYPE_PONG})
			if err != nil {
				return err
			}
			log.Println("recv Ping")
		case msg.TYPE_PONG:
			reader.Discard(msg.MSG_TYPE_SIZE)
			log.Println("recv pong")
		case msg.TYPE_NORMAL:
			_, err = io.ReadAtLeast(reader, header, msg.MSG_HEADER_SIZE)
			if err != nil {
				return err
			}

			m := msg.NewByHeader(header)
			_, err = io.ReadAtLeast(reader, m.Body, int(m.Len))
			if err != nil {
				return err
			}

			seq := binary.BigEndian.Uint32(header[msg.MSG_TYPE_END:msg.MSG_SEQ_END])
			c.Ack(seq)

			c.In <- m.Body
		}

		c.UpdateLastTime()
	}
	return nil
}

func (c *TCPConn) WriteLoop() error {
	for {
		select {
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

func (c *TCPConn) IsClosed() bool {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.closed
}

func getTCPReadDeadline() time.Time {
	return time.Now().Add(time.Second * TCP_READ_TIMEOUT)
}

func (c *TCPConn) Write(bytes []byte) error {
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_NORMAL, new, bytes)
	c.AddMsgToPendingMap(new, m)
	return c.WriteBytes(m.Bytes())
}

func (c *TCPConn) WriteSlice(bytes ...[]byte) error {
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_NORMAL, new, nil)
	for _, s := range bytes {
		m.Len += uint32(len(s))
	}
	m.BodySlice = bytes
	c.AddMsgToPendingMap(new, m)
	err := c.WriteBytes(m.HeaderBytes())
	if err != nil {
		return err
	}

	for _, m := range bytes {
		err := c.WriteBytes(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *TCPConn) SendReg(key cipher.PubKey) error {
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_REG, new, key[:])
	c.AddMsgToPendingMap(new, m)
	return c.WriteBytes(m.Bytes())
}

func (c *TCPConn) WriteBytes(bytes []byte) error {
	index := 0
	for n, err := c.TcpConn.Write(bytes[index:]); index != len(bytes); index += n {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TCPConn) Ack(seq uint32) error {
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	return c.WriteBytes(resp)
}

func (c *TCPConn) Ping() error {
	b := make([]byte, msg.MSG_TYPE_SIZE)
	b[msg.MSG_TYPE_BEGIN] = msg.TYPE_PING
	return c.WriteBytes(b)
}

func (c *TCPConn) GetChanOut() chan<- []byte {
	return c.Out
}

func (c *TCPConn) GetChanIn() <-chan []byte {
	return c.In
}

func (c *TCPConn) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("closing closed udpconn")
		}
	}()
	c.fieldsMutex.Lock()
	if c.closed {
		c.fieldsMutex.Unlock()
		return
	}
	c.closed = true
	c.fieldsMutex.Unlock()
	close(c.In)
	close(c.Out)
}

func (c *TCPConn) GetPublicKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.pubkey
}

func (c *TCPConn) SetPublicKey(key cipher.PubKey) {
	c.fieldsMutex.Lock()
	c.pubkey = key
	c.fieldsMutex.Unlock()
}

func (c *TCPConn) UpdateLastTime() {
	c.TcpConn.SetReadDeadline(getTCPReadDeadline())
}