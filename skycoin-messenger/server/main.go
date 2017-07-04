package main

import (
	"github.com/skycoin/net/skycoin-messenger/factory"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("listening tcp")
	f := factory.NewMessengerFactory()
	err := f.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
