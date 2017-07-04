package main

import (
	"log"
	"github.com/skycoin/net/factory"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tcpFactory := factory.NewTCPFactory()
	conn, err := tcpFactory.Connect(":8080")
	if err != nil {
		panic(err)
	}
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			log.Printf("received msg %s", m)
		}
	}
}

