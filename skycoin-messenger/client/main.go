package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/rpc"
	"github.com/skycoin/net/skycoin-messenger/websocket"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
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
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}
