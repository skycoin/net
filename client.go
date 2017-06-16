package main

import (
	"github.com/skycoin/net/client"
	"log"
)

func main() {
	c := client.New()
	//err := c.Connect("tcp", ":8080")
	err := c.Connect("udp", ":8081")
	if err != nil {
		panic(err)
	}
	go c.Loop()

	c.Out <- []byte("hello1world")
	c.Out <- []byte("hello2world")
	c.Out <- []byte("hello3world")

	for {
		select {
		case m, ok := <-c.In:
			if !ok {
				log.Println("conn closed")
			}
			log.Printf("msg In %x", m)
		}
	}
}
