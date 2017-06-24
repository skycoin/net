package main

import (
	"github.com/skycoin/net/client"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)

	go client2()

	factory := client.NewClientConnectionFactory()
	factory.Connect("tcp", ":8080", cipher.PubKey([33]byte{0xf1}))
	conn := factory.Dial(cipher.PubKey([33]byte{0xf2}))
	conn.Out <- []byte("Hello 0xf2")

	for {
		select {
		case m, ok := <-conn.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg In %s", m)
		}
	}
}

func client2() {
	factory := client.NewClientConnectionFactory()
	factory.SetIncomingCallback(func(conn *client.ClientConnection, data []byte) bool {
		log.Printf("msg from %s In %s", conn.Key.Hex(), data)
		return true
	})
	factory.Connect("udp", ":8081", cipher.PubKey([33]byte{0xf2}))
}