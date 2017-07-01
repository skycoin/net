package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"log"
	"github.com/skycoin/net/conn"
)

type Factory struct {
	clients      map[*Client]bool
	clientsMutex sync.RWMutex
}

func NewFactory() *Factory {
	return &Factory{clients: make(map[*Client]bool)}
}

var (
	once           = &sync.Once{}
	defaultFactory *Factory
)

func GetFactory() *Factory {
	once.Do(func() {
		defaultFactory = NewFactory()
		go defaultFactory.logStatus()
	})
	return defaultFactory
}

func (factory *Factory) NewClient(conn *websocket.Conn) *Client {
	client := &Client{conn: conn, push: make(chan interface{}), PendingMap: conn.PendingMap{Pending: make(map[uint32]interface{})}}
	factory.clientsMutex.Lock()
	factory.clients[client] = true
	factory.clientsMutex.Unlock()
	go func() {
		client.writeLoop()
		factory.clientsMutex.Lock()
		delete(factory.clients, client)
		factory.clientsMutex.Unlock()
	}()
	return client
}

func (factory *Factory) logStatus() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			factory.clientsMutex.RLock()
			log.Printf("websocket factory clients count:%d", len(factory.clients))
			factory.clientsMutex.RUnlock()
		}
	}
}
