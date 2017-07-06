package conn

import (
	"net"
	"github.com/skycoin/net/msg"
	"bufio"
	"io"
	"encoding/binary"
	"time"
	"log"
	"sync/atomic"
	"fmt"
)

const (
	TCP_READ_TIMEOUT = 90
)

type TCPConn struct {
	TcpConn net.Conn
	In      chan []byte
	Out     chan []byte

	*ConnCommonFields
}

func (c *TCPConn) ReadLoop() (err error) {
	defer func() {
		if err != nil {
			c.SetStatusToError(err)
		}
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
			_, err = io.ReadAtLeast(reader, header[:msg.MSG_SEQ_END], msg.MSG_SEQ_END)
			if err != nil {
				return err
			}
			seq := binary.BigEndian.Uint32(header[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.DelMsg(seq)
			c.UpdateLastAck(seq)
		case msg.TYPE_PING:
			reader.Discard(msg.MSG_TYPE_SIZE)
			err = c.WriteBytes([]byte{msg.TYPE_PONG})
			if err != nil {
				return err
			}
			log.Println("recv ping")
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

			seq := binary.BigEndian.Uint32(header[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END])
			c.Ack(seq)

			c.In <- m.Body
		default:
			return fmt.Errorf("not implemented msg type %d", msg_t)
		}
		c.UpdateLastTime()
	}
	return nil
}

func (c *TCPConn) WriteLoop() (err error) {
	defer func() {
		if err != nil {
			c.SetStatusToError(err)
		}
	}()
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
	s := atomic.AddUint32(&c.seq, 1)
	m := msg.New(msg.TYPE_NORMAL, s, bytes)
	c.AddMsg(s, m)
	return c.WriteBytes(m.Bytes())
}

func (c *TCPConn) WriteBytes(bytes []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	//log.Printf("write %x", bytes)
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

func (c *TCPConn) UpdateLastTime() {
	c.TcpConn.SetReadDeadline(getTCPReadDeadline())
}