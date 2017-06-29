package main

import (
	"github.com/skycoin/net/server"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := server.New(":8080", ":8081")
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
