package main

import (
	"github.com/skycoin/net/server"
	"log"
)

func main() {
	s := server.New()
	s.Init()
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
