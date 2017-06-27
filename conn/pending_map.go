package conn

import (
	"github.com/skycoin/net/msg"
	"sync"
)

type PendingMap struct {
	Pending map[uint32]*msg.Message
	sync.RWMutex
}

func (m *PendingMap) AddMsgToPendingMap(k uint32, v *msg.Message) {
	m.Lock()
	m.Pending[k] = v
	m.Unlock()
}

func (m *PendingMap) DelMsgToPendingMap(k uint32) {
	m.Lock()
	delete(m.Pending, k)
	m.Unlock()
}
