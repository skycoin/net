package conn

import (
	"sync"
	"log"
)

type PendingMap struct {
	Pending map[uint32]interface{}
	sync.RWMutex
}

func NewPendingMap() *PendingMap {
	return &PendingMap{Pending: make(map[uint32]interface{})}
}

func (m *PendingMap) AddMsg(k uint32, v interface{}) {
	m.Lock()
	m.Pending[k] = v
	//log.Printf("add %d, Pending:%d, %v", k, len(m.Pending), m.Pending)
	m.Unlock()
}

func (m *PendingMap) DelMsg(k uint32) {
	m.Lock()
	delete(m.Pending, k)
	log.Printf("acked %d, Pending:%d, %v", k, len(m.Pending), m.Pending)
	m.Unlock()
}
