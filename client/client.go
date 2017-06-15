package client

import (
	"net"
	"github.com/skycoin/net/conn"
	"errors"
	"log"
)

var ErrInvalidConnectionType = errors.New("invalid connection type")

type Client struct {
	conn conn.Connection
	In   chan []byte
	Out  chan []byte
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
		cn := conn.NewTCPConn(c)
		client.conn = cn
		client.In = cn.In
		client.Out = cn.Out
	case *net.UDPConn:
		cn := conn.NewUDPClientConn(c)
		client.conn = cn
		client.In = cn.In
		client.Out = cn.Out
	default:
		return ErrInvalidConnectionType
	}
	return nil
}

func (client *Client) Loop() (err error) {
	go client.conn.ReadLoop()
	for {
		select {
		case m, ok := <-client.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg In %x", m)
		case m, ok := <-client.Out:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg Out %x", m)
			err = client.conn.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return
			}
		}
	}
}
