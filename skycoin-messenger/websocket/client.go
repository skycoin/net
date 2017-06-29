package websocket

import (
	"github.com/gorilla/websocket"
	"time"
	"log"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"sync"
	"github.com/skycoin/net/client"
	"encoding/json"
	"github.com/skycoin/net/util"
)

type Client struct {
	factory *client.ClientConnectionFactory
	sync.RWMutex

	conn *websocket.Conn
	push chan interface{}
}

func (c *Client) GetFactory() *client.ClientConnectionFactory {
	c.RLock()
	defer c.RUnlock()
	return c.factory
}

func (c *Client) SetFactory(factory *client.ClientConnectionFactory) {
	c.Lock()
	if c.factory != nil {
		c.factory.Close()
	}
	c.factory = factory
	c.Unlock()
}

func (c *Client) PushLoop(conn *client.ClientConnection, data []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{PublicKey: conn.Key.Hex(), Msg: util.ByteSlice2String(data)}
	c.push <- push
	for {
		select {
		case m, ok := <-conn.In:
			if !ok {
				return
			}
			push.Msg = util.ByteSlice2String(m)
			c.push <- push
		}
	}
}

func (c *Client) readLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("readLoop recovered err %v", err)
		}
		c.SetFactory(nil)
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
		op := msg.GetOP(opn)
		if op == nil {
			log.Printf("op not found, %d", opn)
			continue
		}
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
		c.SetFactory(nil)
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			switch m := message.(type) {
			case msg.PushMsg:
				_, err := w.Write([]byte{msg.PUSH_MSG})
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
