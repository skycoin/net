package main

import (
	"github.com/skycoin/net/client"
	"time"
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
	time.Sleep(time.Second * 2)
}
