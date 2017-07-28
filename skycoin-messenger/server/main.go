package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

func main() {
	log.Println("listening tcp")
	f := factory.NewMessengerFactory()
	err := f.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
