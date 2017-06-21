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
	factory *ConnectionFactory
	tcpConn *net.TCPConn
	In      chan []byte
	Out     chan interface{}

	seq     uint32
	pending map[uint32]*msg.Message

	closed bool
	pubkey cipher.PubKey
	fieldsMutex *sync.RWMutex
}

func NewTCPConn(c *net.TCPConn, factory *ConnectionFactory) *TCPConn {
	return &TCPConn{tcpConn: c, factory: factory ,In: make(chan []byte), Out: make(chan interface{}), pending: make(map[uint32]*msg.Message), fieldsMutex:new(sync.RWMutex)}
}

func (c *TCPConn) ReadLoop() error {
	defer func() {
		c.close()
	}()
	header := make([]byte, msg.MSG_HEADER_SIZE)
	reader := bufio.NewReader(c.tcpConn)

	for {
		t, err := reader.Peek(msg.MSG_TYPE_SIZE)
		if err != nil {
			return err
		}
		msg_t := t[msg.MSG_TYPE_BEGIN]
		switch msg_t {
		case msg.TYPE_ACK:
			reader.Discard(msg.MSG_TYPE_SIZE)
			_, err = io.ReadAtLeast(reader, header, msg.MSG_SEQ_END)
			if err != nil {
				return err
			}
			seq := binary.BigEndian.Uint32(header[msg.MSG_TYPE_END:msg.MSG_SEQ_END])
			delete(c.pending, seq)
		case msg.TYPE_PING:
			reader.Discard(msg.MSG_TYPE_SIZE)
			err = c.writeBytes([]byte{msg.TYPE_PONG})
			if err != nil {
				return err
			}
		case msg.TYPE_PONG:
			reader.Discard(msg.MSG_TYPE_SIZE)
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
			c.ack(seq)

			c.In <- m.Body
		case msg.TYPE_REG:
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
			c.ack(seq)

			if m.Len != 33 {
				continue
			}
			key := cipher.NewPubKey(m.Body)
			c.fieldsMutex.Lock()
			c.pubkey = key
			c.fieldsMutex.Unlock()
			c.factory.Register(key.Hex(), c)
		}

		c.tcpConn.SetReadDeadline(getTCPReadDeadline())
	}
	return nil
}

func (c *TCPConn) WriteLoop() error {
	ticker := time.NewTicker(time.Second * TICK_PERIOD)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			err := c.ping()
			if err != nil {
				return err
			}
		case ob, ok := <-c.Out:
			if !ok {
				log.Println("conn closed")
				return nil
			}
			switch m := ob.(type) {
			case []byte:
				log.Printf("msg Out %x", m)
				err := c.Write(m)
				if err != nil {
					log.Printf("write msg is failed %v", err)
					return err
				}
			case [][]byte:
				log.Printf("msg Out %x", m)
				err := c.WriteSlice(m)
				if err != nil {
					log.Printf("write msg is failed %v", err)
					return err
				}
			default:
				log.Printf("WriteLoop writting %#v failed unsupported type", ob)
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
	c.pending[new] = m
	return c.writeBytes(m.Bytes())
}

func (c *TCPConn) WriteSlice(bytes [][]byte) error {
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_NORMAL, new, nil)
	for _, s := range bytes {
		m.Len += uint32(len(s))
	}
	m.BodySlice = bytes
	c.pending[new] = m
	err := c.writeBytes(m.HeaderBytes())
	if err != nil {
		return err
	}

	for _, m := range bytes {
		err := c.writeBytes(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *TCPConn) SendReg(key cipher.PubKey) error {
	new := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_REG, new, key[:])
	c.pending[new] = m
	return c.writeBytes(m.Bytes())
}

func (c *TCPConn) writeBytes(bytes []byte) error {
	index := 0
	for n, err := c.tcpConn.Write(bytes[index:]); index != len(bytes); index += n {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TCPConn) ack(seq uint32) error {
	resp := make([]byte, msg.MSG_SEQ_END)
	resp[msg.MSG_TYPE_BEGIN] = msg.TYPE_ACK
	binary.BigEndian.PutUint32(resp[msg.MSG_SEQ_BEGIN:], seq)
	return c.writeBytes(resp)
}

func (c *TCPConn) ping() error {
	b := make([]byte, msg.MSG_TYPE_SIZE)
	b[msg.MSG_TYPE_BEGIN] = msg.TYPE_PING
	return c.writeBytes(b)
}

func (c *TCPConn) GetChanOut() chan<- interface{} {
	return c.Out
}

func (c *TCPConn) close() {
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
	c.factory.UnRegister(c.pubkey.Hex(), c)
}

func (c *TCPConn) GetPublicKey() cipher.PubKey {
	c.fieldsMutex.RLock()
	defer c.fieldsMutex.RUnlock()
	return c.pubkey
}