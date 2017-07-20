package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/factory"
)

func main() {
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
