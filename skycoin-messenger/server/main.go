package main

import (
	"github.com/skycoin/net/skycoin-messenger/factory"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("listening tcp")
	f := factory.NewMessengerFactory()
	err := f.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
