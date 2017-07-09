package factory

import (
	"github.com/skycoin/net/server"
	"net"
	"github.com/skycoin/net/client"
	"github.com/skycoin/net/conn"
)

type TCPFactory struct {
	listener *net.TCPListener

	FactoryCommonFields
}

func NewTCPFactory() *TCPFactory {
	return &TCPFactory{FactoryCommonFields:NewFactoryCommonFields()}
}

func (factory *TCPFactory) Listen(address string) error {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	factory.listener = ln
	for {
		c, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		factory.createConn(c)
	}
}

func (factory *TCPFactory) Close() error {
	return factory.listener.Close()
}

func (factory *TCPFactory) createConn(c *net.TCPConn) *Connection {
	tcpConn := server.NewServerTCPConn(c)
	tcpConn.SetStatusToConnected()
	conn := &Connection{Connection: tcpConn, factory: factory}
	factory.AddConn(conn)
	factory.AcceptedCallback(conn)
	return conn
}

func (factory *TCPFactory) Connect(address string) (conn *Connection, err error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	cn := client.NewClientTCPConn(c)
	cn.SetStatusToConnected()
	conn = &Connection{Connection: cn, factory: factory}
	factory.AddConn(conn)
	return
}

