package conn

import (
	"github.com/skycoin/net/msg"
	"sync"
)

type PendingMap struct {
	pending map[uint32]*msg.Message
	sync.RWMutex
}

func (m *PendingMap) addMsgToPendingMap(k uint32, v *msg.Message) {
	m.Lock()
	m.pending[k] = v
	m.Unlock()
}

func (m *PendingMap) delMsgToPendingMap(k uint32) {
	m.Lock()
	delete(m.pending, k)
	m.Unlock()
}
