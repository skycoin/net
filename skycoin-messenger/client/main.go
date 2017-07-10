package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/skycoin/net/skycoin-messenger/websocket"
	"github.com/skycoin/net/skycoin-messenger/rpc"
)

func main() {
	go func() {
		log.Println("listening rpc")
		err := rpc.ServeRPC("localhost:8083")
		if err != nil {
			log.Fatal("rpc.ServeRPC: ", err)
		}
	}()

	log.Println("listening websocket")
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})
	err := http.ListenAndServe("localhost:8082", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}
