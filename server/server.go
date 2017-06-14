package server

import (
	"net"
	"github.com/skycoin/net/conn"
	"log"
)

type Server struct {
	TCPAddress string
}

func New() *Server {
	return &Server{TCPAddress: ":8080"}
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
		go server.handleTCP(c)
	}
}

func (server *Server) handleTCP(c *net.TCPConn) {
	connection := conn.NewTCPConn(c)
	go connection.ReadLoop()
	for {
		select {
		case m, ok := <-connection.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg in %v", m)
		case m, ok := <-connection.Out:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg out %v", m)
			err := connection.Write(m)
			if err != nil {
				log.Printf("write msg is failed %v", err)
				return
			}
		}
	}
}
