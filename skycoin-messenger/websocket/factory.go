package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Factory struct {
	clients map[*Client]bool
	clientsMutex sync.RWMutex
}

func NewFactory() *Factory {
	return &Factory{clients:make(map[*Client]bool)}
}

var (
	once = &sync.Once{}
	defaultFactory *Factory
)

func GetFactory() *Factory {
	once.Do(func() {
		defaultFactory = NewFactory()
	})
	return defaultFactory
}

func (factory *Factory) NewClient(conn *websocket.Conn) *Client {
	client := &Client{conn: conn, send: make(chan []byte)}
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
