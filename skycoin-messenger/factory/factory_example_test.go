package factory_test

import (
	"github.com/skycoin/net/skycoin-messenger/factory"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

const recvMsgOffset = 67

var receiverPk cipher.PubKey

// This example shows how to send ping pong messages through two clients
func Example(){
	go startServer()
	go startReceiver()

	// Wait for receiver to be ready
	time.Sleep(time.Second*2)
	startSender()
	// Wait for message to arrive to receiver
	time.Sleep(time.Second*2)
}

func startServer() {
	f := factory.NewMessengerFactory()
	fmt.Println("serving")
	err := f.Listen(":8000")
	if err != nil {
		fmt.Println("error starting server...", err)
	}
}

func startReceiver() {
	f := factory.NewMessengerFactory()

	fmt.Println("receiver connecting...")
	err :=	f.Connect("localhost:8000")
	if err != nil {
		panic(err)
	}

	fmt.Println("receiver connected")
	f.ForEachConn(processConn)
}

func processConn(conn *factory.Connection) {
	receiverPk = conn.GetKey()
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				panic("something is not ok...")
				return
			}
			fmt.Printf("got message: %s\n",string(m[recvMsgOffset:]))
		}
	}
}

func startSender() {
	f := factory.NewMessengerFactory()

	fmt.Println("sender connecting...")
	err :=	f.Connect("localhost:8000")
	if err != nil {
		panic(err)
	}
	fmt.Println("sender connected")


	f.ForEachConn(processConn2)
}

func processConn2(conn *factory.Connection){
	msg :=	factory.GenSendMsg(conn.GetKey(),receiverPk,[]byte("hey!") )
	conn.Write(msg)
}
