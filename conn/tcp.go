package conn

import (
	"net"
	"github.com/skycoin/net/msg"
	"bufio"
	"io"
	"encoding/binary"
)

type TCPConn struct {
	tcpConn *net.TCPConn
	In      chan *msg.Message
	Out     chan *msg.Message

	seq     uint32
	pending map[uint32]*msg.Message
}

func NewTCPConn(c *net.TCPConn) *TCPConn {
	return &TCPConn{tcpConn: c, In: make(chan *msg.Message), Out: make(chan *msg.Message)}
}

func (c *TCPConn) ReadLoop() error {
	defer func() {
		close(c.In)
	}()
	header := make([]byte, msg.MSG_HEADER_SIZE)
	reader := bufio.NewReader(c.tcpConn)

	for {
		t, err := reader.Peek(msg.MSG_TYPE_SIZE)
		if err != nil {
			return err
		}
		switch t[0] {
		case msg.TYPE_ACK:
			_, err = io.ReadAtLeast(reader, header, msg.MSG_SEQ_END)
			if err != nil {
				return err
			}
			seq := binary.BigEndian.Uint32(header[msg.MSG_TYPE_END:msg.MSG_SEQ_END])
			delete(c.pending, seq)
			continue
		case msg.TYPE_PING:
			err = c.WriteBytes([]byte{msg.TYPE_PING})
			if err != nil {
				return err
			}
			continue
		}

		_, err = io.ReadAtLeast(reader, header, msg.MSG_HEADER_SIZE)
		if err != nil {
			return err
		}

		m := msg.New(header)
		_, err = io.ReadAtLeast(reader, m.Body, int(m.Len))
		if err != nil {
			return err
		}

		c.In <- m
	}
	return nil
}

func (c *TCPConn) Write(msg *msg.Message) error {
	c.seq++
	c.pending[c.seq] = msg
	return c.WriteBytes(msg.Bytes())
}

func (c *TCPConn) WriteBytes(bytes []byte) error {
	index := 0
	for n, err := c.tcpConn.Write(bytes[index:]); index != len(bytes); index += n {
		if err != nil {
			return err
		}
	}
	return nil
}
