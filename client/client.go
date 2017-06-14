package client

import (
	"net"
	"github.com/skycoin/net/conn"
	"errors"
	"github.com/skycoin/net/msg"
	"log"
)

var ErrInvalidConnectionType = errors.New("invalid connection type")

type Client struct {
	conn conn.Connection
	in   chan *msg.Message
	out  chan *msg.Message
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
		client.in = cn.In
		client.out = cn.Out
	case *net.UDPConn:
	default:
		return ErrInvalidConnectionType
	}
	return nil
}

func (client *Client) Loop() (err error) {
	go client.conn.ReadLoop()
	for {
		select {
		case m, ok := <-client.in:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg in %v", m)
		case m, ok := <-client.out:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg out %v", m)
			err = client.conn.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return
			}
		}
	}
}
