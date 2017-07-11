package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
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
	wsId           uint32
)

func GetFactory() *Factory {
	once.Do(func() {
		defaultFactory = NewFactory()
		go defaultFactory.logStatus()
	})
	return defaultFactory
}

func (factory *Factory) NewClient(c *websocket.Conn) *Client {
	logger := log.WithField("wsId", atomic.AddUint32(&wsId, 1))
	client := &Client{conn: c, push: make(chan interface{}), PendingMap: PendingMap{Pending: make(map[uint32]interface{})}, logger: logger}
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
			log.Debugf("websocket connection clients count:%d", len(factory.clients))
			factory.clientsMutex.RUnlock()
		}
	}
}
