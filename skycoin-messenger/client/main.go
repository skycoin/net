package main

import (
	"log"
	"net/http"
	"github.com/skycoin/net/skycoin-messenger/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
