package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"sync"
	"encoding/json"
	"io"
	"encoding/binary"
	"sync/atomic"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

type Client struct {
	connection *factory.Connection
	sync.RWMutex

	seq uint32
	PendingMap

	conn *websocket.Conn
	push chan interface{}

	logger *log.Entry
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

func (c *Client) PushLoop(conn *factory.Connection) {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Errorf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{}
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			key := cipher.NewPubKey(m[factory.MSG_PUBLIC_KEY_BEGIN:factory.MSG_PUBLIC_KEY_END])
			push.From = key.Hex()
			push.Msg = string(m[factory.MSG_META_END:])
			c.push <- push
		}
	}
}

func (c *Client) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			c.logger.Errorf("readLoop recovered err %v", err)
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
				c.logger.Errorf("error: %v", err)
			}
			c.logger.Errorf("error: %v", err)
			break
		}
		c.logger.Debugf("recv %x", m)
		opn := int(m[msg.MSG_OP_BEGIN])
		if opn == msg.OP_ACK {
			c.DelMsg(binary.BigEndian.Uint32(m[msg.MSG_SEQ_BEGIN:msg.MSG_SEQ_END]))
			continue
		}
		op := msg.GetOP(opn)
		if op == nil {
			c.logger.Errorf("op not found, %d", opn)
			continue
		}

		c.ack(m[msg.MSG_OP_BEGIN:msg.MSG_SEQ_END])

		err = json.Unmarshal(m[msg.MSG_HEADER_END:], op)
		if err == nil {
			err = op.Execute(c)
			if err != nil {
				c.logger.Errorf("websocket readLoop executed err: %v", err)
			}
		} else {
			c.logger.Errorf("websocket readLoop json Unmarshal err: %v", err)
		}
		msg.PutOP(opn, op)
	}
}

func (c *Client) writeLoop() (err error) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if err := recover(); err != nil {
			c.logger.Errorf("writeLoop recovered err %v", err)
		}
		ticker.Stop()
		c.SetConnection(nil)
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.push:
			c.logger.Debug("push", message)
			if !ok {
				c.logger.Debug("closed c.push")
				err = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					c.logger.Error(err)
					return
				}
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				c.logger.Error(err)
				return err
			}
			switch m := message.(type) {
			case *msg.PushMsg:
				err = c.write(w, m)
				if err != nil {
					c.logger.Error(err)
					return err
				}
			default:
				c.logger.Errorf("not implemented msg %v", m)
			}
			if err = w.Close(); err != nil {
				c.logger.Error(err)
				return err
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.logger.Error(err)
				return err
			}
		}
	}
}

func (c *Client) write(w io.WriteCloser, m *msg.PushMsg) (err error) {
	_, err = w.Write([]byte{msg.PUSH_MSG})
	c.logger.Debugf("op %d", msg.PUSH_MSG)
	if err != nil {
		return
	}
	ss := make([]byte, 4)
	nseq := atomic.AddUint32(&c.seq, 1)
	c.AddMsg(nseq, m)
	binary.BigEndian.PutUint32(ss, nseq)
	_, err = w.Write(ss)
	c.logger.Debugf("seq %x", ss)
	if err != nil {
		return
	}
	jbs, err := json.Marshal(m)
	if err != nil {
		return
	}
	_, err = w.Write(jbs)
	c.logger.Debugf("json %x", jbs)
	if err != nil {
		return
	}

	return nil
}

func (c *Client) ack(data []byte) error {
	data[msg.MSG_OP_BEGIN] = msg.PUSH_ACK
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
