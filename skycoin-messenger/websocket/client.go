package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"log"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"sync"
	"encoding/json"
	"io"
	"encoding/binary"
	"sync/atomic"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/conn"
)

type Client struct {
	connection *factory.Connection
	sync.RWMutex

	seq uint32
	conn.PendingMap

	conn *websocket.Conn
	push chan interface{}
}

func (c *Client) GetConnection() *factory.Connection {
	c.RLock()
	defer c.RUnlock()
	return c.connection
}

func (c *Client) SetConnection(connection *factory.Connection) {
	c.Lock()
	if c.connection != nil {
		c.connection.Close()
	}
	c.connection = connection
	c.Unlock()
}

func (c *Client) PushLoop(conn *factory.Connection, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{PublicKey: conn.Key.Hex(), Msg: string(data)}
	c.push <- push
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			push.Msg = string(m)
			c.push <- push
		}
	}
}

func (c *Client) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("readLoop recovered err %v", err)
		}
		c.SetConnection(nil)
		c.conn.Close()
		close(c.push)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		opn := int(m[msg.MSG_OP_BEGIN])
		if opn == msg.OP_ACK {
			c.DelMsg(binary.BigEndian.Uint32(m[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END]))
			continue
		}
		op := msg.GetOP(opn)
		if op == nil {
			log.Printf("op not found, %d", opn)
			continue
		}

		c.ack(m[msg.MSG_OP_BEGIN:msg.MSG_SEQ_END])

		err = json.Unmarshal(m[msg.MSG_HEADER_END:], op)
		if err == nil {
			err = op.Execute(c)
			if err != nil {
				log.Printf("websocket readLoop executed err: %v", err)
			}
		} else {
			log.Printf("websocket readLoop json Unmarshal err: %v", err)
		}
		msg.PutOP(opn, op)
	}
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("writeLoop recovered err %v", err)
		}
		ticker.Stop()
		c.SetConnection(nil)
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.push:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			switch m := message.(type) {
			case msg.PushMsg:
				c.write(w, &m)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) write(w io.WriteCloser, m *msg.PushMsg) (err error) {
	_, err = w.Write([]byte{msg.PUSH_MSG})
	if err != nil {
		return
	}
	ss := make([]byte, 4)
	nseq := atomic.AddUint32(&c.seq, 1)
	c.AddMsg(nseq, m)
	binary.BigEndian.PutUint32(ss, nseq)
	_, err = w.Write(ss)
	if err != nil {
		return
	}
	jbs, err := json.Marshal(m)
	if err != nil {
		return
	}
	_, err = w.Write(jbs)
	if err != nil {
		return
	}

	return nil
}

func (c *Client) ack(data []byte) error {
	data[msg.MSG_OP_BEGIN] = msg.PUSH_ACK
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
