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
)

const (
	TCP_READ_TIMEOUT = 90
)

type TCPConn struct {
	tcpConn *net.TCPConn
	In      chan []byte
	Out     chan []byte

	seq     uint32
	pending map[uint32]*msg.Message

	closed bool
	pubkey string
	fieldsMutex *sync.RWMutex
}

func NewTCPConn(c *net.TCPConn) *TCPConn {
	return &TCPConn{tcpConn: c, In: make(chan []byte), Out: make(chan []byte), pending: make(map[uint32]*msg.Message), fieldsMutex:new(sync.RWMutex)}
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
		switch t[msg.MSG_TYPE_BEGIN] {
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
		}

		c.tcpConn.SetReadDeadline(getTCPReadDeadline())
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
	c.seq++
	m := msg.New(msg.TYPE_NORMAL, c.seq, bytes)
	c.pending[c.seq] = m
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

func (c *TCPConn) Ping() error {
	b := make([]byte, msg.MSG_TYPE_SIZE)
	b[msg.MSG_TYPE_BEGIN] = msg.TYPE_PING
	return c.writeBytes(b)
}

func (c *TCPConn) GetChanOut() chan<- []byte {
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
}