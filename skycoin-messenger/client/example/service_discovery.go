package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

func main() {
	f := factory.NewMessengerFactory()
	conn, err := f.Connect(":8080")
	if err != nil {
		panic(err)
	}
	conn.OfferService("tracker")
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			log.Printf("in %x", m)
			if m[factory.MSG_OP_BEGIN] == factory.OP_OFFER_SERVICE {
				log.Printf("offer services %v", conn.GetServices())
			}
		}
	}
}
