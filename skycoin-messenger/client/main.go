package main

import (
	"log"
	"net/http"
	"github.com/skycoin/net/skycoin-messenger/websocket"
	"github.com/skycoin/net/skycoin-messenger/rpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})
	go func() {
		log.Println("listening rpc")
		rpc.ServeRPC(":8083")
	}()
	log.Println("listening websocket")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
