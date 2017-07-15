package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/factory"
)

func main() {
	log.SetLevel(log.DebugLevel)
	go func() {
		log.Println("listening udp")
		udpFactory := factory.NewUDPFactory()
		udpFactory.AcceptedCallback = func(connection *factory.Connection) {
			connection.GetChanOut() <- []byte("hello")
		}
		udpFactory.Listen(":8081")
	}()
	log.Println("listening tcp")
	tcpFactory := factory.NewTCPFactory()
	tcpFactory.AcceptedCallback = func(connection *factory.Connection) {
		connection.GetChanOut() <- []byte("hello")
	}
	if err := tcpFactory.Listen(":8080"); err != nil {
		panic(err)
	}
}
