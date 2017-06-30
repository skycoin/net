package server

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"github.com/skycoin/net/conn"
	"github.com/skycoin/net/msg"
	"github.com/skycoin/skycoin/src/cipher"
	"net"
)

type ServerTCPConn struct {
	conn.TCPConn
	factory *ConnectionFactory
}

func NewServerTCPConn(c *net.TCPConn, factory *ConnectionFactory) *ServerTCPConn {
	return &ServerTCPConn{TCPConn: conn.TCPConn{TcpConn: c, In: make(chan []byte), Out: make(chan []byte), PendingMap: conn.PendingMap{Pending: make(map[uint32]*msg.Message)}}, factory: factory}
}

func (c *ServerTCPConn) ReadLoop() error {
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

			seq := binary.BigEndian.Uint32(header[msg.MSG_TYPE_END:msg.MSG_SEQ_END])
			c.Ack(seq)

			func() {
				defer func() {
					if err := recover(); err != nil {
						log.Printf("c.In <- m.Body panic %v, %x", err, m.Body)
					}
				}()
				c.In <- m.Body
			}()
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
			c.Ack(seq)

			if m.Len != 33 {
				continue
			}
			key := cipher.NewPubKey(m.Body)
			c.SetPublicKey(key)
			c.factory.Register(key.Hex(), c)
		}

		c.UpdateLastTime()
	}
	return nil
}
