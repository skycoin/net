package server

import (
	"net"
	"log"
	"github.com/skycoin/net/conn"
)

var (
	DefaultConnectionFactory = conn.NewFactory()
)

type Server struct {
	TCPAddress string
	UDPAddress string
}

func New() *Server {
	return &Server{TCPAddress: ":8080", UDPAddress:":8081"}
}

func (server *Server) Init() {
	DefaultConnectionFactory.TCPClientHandler = handleTCP
	DefaultConnectionFactory.UDPClientHandler = handleUDP
}

func (server *Server) ListenTCP() error {
	addr, err := net.ResolveTCPAddr("tcp", server.TCPAddress)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	for {
		c, err := ln.AcceptTCP()
		if err != nil {
			return err
		}
		connection := DefaultConnectionFactory.CreateTCPConn(c)
		go connection.ReadLoop()
	}
}

func handleTCP(connection *conn.TCPConn) {
	for {
		select {
		case m, ok := <-connection.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg in %x", m)
		case m, ok := <-connection.Out:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg out %x", m)
			err := connection.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return
			}
		}
	}
}

func (server *Server) ListenUDP() error {
	addr, err := net.ResolveUDPAddr("udp", server.UDPAddress)
	if err != nil {
		return err
	}
	udp, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	udpc := conn.NewUDPServerConn(udp, DefaultConnectionFactory)
	return udpc.ReadLoop()
}

func handleUDP(connection *conn.UDPConn) {
	for {
		select {
		case m, ok := <-connection.In:
			if !ok {
				log.Println("udp conn closed")
				return
			}
			log.Printf("msg in %x", m)
		case m, ok := <-connection.Out:
			if !ok {
				log.Println("udp conn closed")
				return
			}
			log.Printf("msg out %x", m)
			err := connection.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return
			}
		}
	}
}
