package conn

import (
	"sync"
)

type PendingMap struct {
	Pending map[uint32]interface{}
	sync.RWMutex
}

func (m *PendingMap) AddMsgToPendingMap(k uint32, v interface{}) {
	m.Lock()
	m.Pending[k] = v
	m.Unlock()
}

func (m *PendingMap) DelMsgToPendingMap(k uint32) {
	m.Lock()
	delete(m.Pending, k)
	m.Unlock()
}
