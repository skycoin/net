package client

import (
	"net"
	"time"

	"github.com/skycoin/net/conn"
)

// ClientTCPConn is a wrapper over skycoin TCPConn
type ClientTCPConn struct {
	conn.TCPConn
}

// NewClientTCPConn gets a net.Conn and returns own wrapper struct over TCPConn
func NewClientTCPConn(c net.Conn) *ClientTCPConn {
	return &ClientTCPConn{
		TCPConn: conn.TCPConn{
			TcpConn:          c,
			ConnCommonFields: conn.NewConnCommonFileds(),
		},
	}
}

// WriteLoop writes over the connection every TCP_PING_TICK_PERIOD what is send to
// the Out channel. Also keeps a ping on every period. WriteLoop exits when the ping fails,
// when write fails or when the Out channel is closed.
func (c *ClientTCPConn) WriteLoop() (err error) {
	ticker := time.NewTicker(time.Second * conn.TCP_PING_TICK_PERIOD)
	defer func() {
		ticker.Stop()
		if err != nil {
			c.SetStatusToError(err)
		}
	}()
	for {
		select {
		case <-ticker.C:
			err := c.Ping()
			if err != nil {
				return err
			}
		case m, ok := <-c.Out:
			if !ok {
				c.GetContextLogger().Debug("conn closed")
				return nil
			}
			//c.GetContextLogger().Debugf("msg Out %x", m)
			err := c.Write(m)
			if err != nil {
				c.GetContextLogger().Debugf("write msg is failed %v", err)
				return err
			}
		}
	}
}
