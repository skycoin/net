package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	f := factory.NewMessengerFactory()
	conn, err := f.Connect(":8080")
	if err != nil {
		panic(err)
	}

	err = conn.Reg()
	if err != nil {
		panic(err)
	}

	go client2()

	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			log.Printf("received msg %x", m)
		}
	}
}

func client2() {
	f := factory.NewMessengerFactory()
	conn, err := f.Connect(":8080")
	if err != nil {
		panic(err)
	}

	key := cipher.PubKey([33]byte{0xf2})
	conn.GetChanOut() <- factory.GenRegMsg()

	f1 := cipher.PubKey([33]byte{0xf1})
	conn.GetChanOut() <- factory.GenSendMsg(key, f1, []byte("Hello 0xf1 1"))
	conn.GetChanOut() <- factory.GenSendMsg(key, f1, []byte("Hello 0xf1 2"))
	conn.Write(factory.GenSendMsg(key, f1, []byte("Hello 0xf1 3")))
	conn.Write(factory.GenSendMsg(key, f1, []byte("Hello 0xf1 4")))
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			log.Printf("received msg %x", m)
		}
	}
}
