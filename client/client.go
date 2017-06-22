package client

import (
	"net"
	"github.com/skycoin/net/conn"
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

var ErrInvalidConnectionType = errors.New("invalid connection type")

type Client struct {
	conn conn.Connection

	In   <-chan []byte
	Out  chan<- interface{}
}

func New() *Client {
	return &Client{}
}

func (client *Client) Connect(network, address string) error {
	c, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	switch c := c.(type) {
	case *net.TCPConn:
		cn := conn.NewClientTCPConn(c)
		client.conn = cn
		client.In = cn.In
		client.Out = cn.Out
	case *net.UDPConn:
		cn := conn.NewClientUDPConn(c)
		client.conn = cn
		client.In = cn.In
		client.Out = cn.Out
	default:
		return ErrInvalidConnectionType
	}
	return nil
}

func (client *Client) Reg(key cipher.PubKey) error {
	return client.conn.SendReg(key)
}

func (client *Client) Loop() error {
	go client.conn.ReadLoop()
	return client.conn.WriteLoop()
}

func (client *Client) IsClosed() bool {
	return client.conn.IsClosed()
}