package main

import (
	"github.com/skycoin/net/client"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	go client2()
	c := client.New()
	err := c.Connect("tcp", ":8080")
	//err := c.Connect("udp", ":8081")
	if err != nil {
		panic(err)
	}
	go c.Loop()

	c.Reg(cipher.PubKey([33]byte{0xf1}))
	factory := client.NewClientConnectionFactory(c)
	conn := factory.GetConn(cipher.PubKey([33]byte{0xf2}))
	conn.Write([]byte("Hello 0xf2"))
	for {
		select {
		case m, ok := <-c.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg In %x", m)
		}
	}
}

func client2() {
	c := client.New()
	err := c.Connect("tcp", ":8080")
	//err := c.Connect("udp", ":8081")
	if err != nil {
		panic(err)
	}
	go c.Loop()

	c.Reg(cipher.PubKey([33]byte{0xf2}))
	factory := client.NewClientConnectionFactory(c)
	conn := factory.GetConn(cipher.PubKey([33]byte{0xf2}))
	conn.Write([]byte("Hello 0xf2"))
	for {
		select {
		case m, ok := <-c.In:
			if !ok {
				log.Println("conn closed")
				return
			}
			log.Printf("msg In %x", m)
		}
	}
}
