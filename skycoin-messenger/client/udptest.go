package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/factory"
)

func main() {
	log.SetLevel(log.DebugLevel)
	f := factory.NewUDPFactory()
	conn, err := f.Connect("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	conn.GetChanOut() <- []byte("123")
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

