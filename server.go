package main

import (
	"github.com/skycoin/net/server"
	"log"
)

func main() {
	s := server.New()
	log.Println("listening")
	if err := s.ListenTCP(); err != nil {
		panic(err)
	}
}
