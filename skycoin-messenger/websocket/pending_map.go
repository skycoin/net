package websocket

import (
	"sync"
	log "github.com/sirupsen/logrus"
)

type PendingMap struct {
	Pending map[uint32]interface{}
	sync.RWMutex
}

func (m *PendingMap) AddMsg(k uint32, v interface{}) {
	m.Lock()
	m.Pending[k] = v
	m.Unlock()
}

func (m *PendingMap) DelMsg(k uint32) {
	m.Lock()
	delete(m.Pending, k)
	log.Printf("acked %d, Pending:%d, %v", k, len(m.Pending), m.Pending)
	m.Unlock()
}
