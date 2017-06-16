package main

import (
	"github.com/skycoin/net/server"
	"log"
	"github.com/skycoin/net/conn"
)

func main() {
	s := server.New()
	server.DefaultConnectionFactory.TCPClientHandler = handleTCP
	server.DefaultConnectionFactory.UDPClientHandler = handleUDP
	go func() {
		log.Println("listening udp")
		if err := s.ListenUDP(); err != nil {
			panic(err)
		}
	}()
	log.Println("listening tcp")
	if err := s.ListenTCP(); err != nil {
		panic(err)
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
		}
	}
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
		}
	}
}
