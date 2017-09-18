package conn

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/skycoin/net/msg"
)

type PendingMap struct {
	Pending              map[uint32]msg.Interface
	sync.RWMutex
	ackedMessages        map[uint32]msg.Interface
	ackedMessagesMutex   sync.RWMutex
	lastMinuteAcked      map[uint32]msg.Interface
	lastMinuteAckedMutex sync.RWMutex

	statistics string
}

func NewPendingMap() *PendingMap {
	pendingMap := &PendingMap{Pending: make(map[uint32]msg.Interface), ackedMessages: make(map[uint32]msg.Interface)}
	go pendingMap.analyse()
	return pendingMap
}

func (m *PendingMap) AddMsg(k uint32, v *msg.Message) {
	m.Lock()
	m.Pending[k] = v
	m.Unlock()
	v.Transmitted()
}

func (m *PendingMap) DelMsg(k uint32) (ok bool) {
	m.RLock()
	v, ok := m.Pending[k]
	m.RUnlock()

	if !ok {
		return
	}

	v.Acked()

	m.ackedMessagesMutex.Lock()
	m.ackedMessages[k] = v
	m.ackedMessagesMutex.Unlock()

	m.Lock()
	delete(m.Pending, k)
	m.Unlock()
	return
}

func (m *PendingMap) analyse() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			m.ackedMessagesMutex.Lock()
			m.lastMinuteAckedMutex.Lock()
			m.lastMinuteAcked = m.ackedMessages
			m.lastMinuteAckedMutex.Unlock()
			m.ackedMessages = make(map[uint32]msg.Interface)
			m.ackedMessagesMutex.Unlock()

			m.lastMinuteAckedMutex.RLock()
			if len(m.lastMinuteAcked) < 1 {
				m.lastMinuteAckedMutex.RUnlock()
				continue
			}
			var max, min int64
			sum := new(big.Int)
			bytesSent := 0
			for _, v := range m.lastMinuteAcked {
				latency := v.GetRTT().Nanoseconds()
				if max < latency {
					max = latency
				}
				if min == 0 || min > latency {
					min = latency
				}
				y := new(big.Int)
				y.SetInt64(latency)
				sum.Add(sum, y)

				bytesSent += v.TotalSize()
			}
			n := new(big.Int)
			n.SetInt64(int64(len(m.lastMinuteAcked)))
			avg := new(big.Int)
			avg.Div(sum, n)
			m.lastMinuteAckedMutex.RUnlock()

			m.statistics = fmt.Sprintf("sent: %d bytes, latency: max %d ns, min %d ns, avg %s ns, count %s", bytesSent, max, min, avg, n)
		}
	}
}

type UDPPendingMap struct {
	*PendingMap
	waitBits byte
	waitCond *sync.Cond
	ringMask uint32
}

func NewUDPPendingMap() *UDPPendingMap {
	m := &UDPPendingMap{PendingMap: NewPendingMap()}
	m.waitCond = sync.NewCond(&m.RWMutex)
	m.ringMask = 16
	return m
}

func (m *UDPPendingMap) AddMsg(k uint32, v msg.Interface) {
	m.Lock()
	i := k % m.ringMask
	for m.waitBits&(1<<i) > 0 {
		m.waitCond.Wait()
	}
	m.Pending[k] = v
	m.waitBits |= 1 << i
	m.Unlock()
}

func (m *UDPPendingMap) DelMsgAndGetLossMsgs(k uint32) (ok bool, loss []msg.Interface) {
	m.Lock()
	v, ok := m.Pending[k]
	if !ok {
		m.Unlock()
		return
	}
	v.Acked()
	delete(m.Pending, k)
	i := k % m.ringMask
	m.waitBits &^= 1 << i
	var prev byte
	prev = ^(1 << i) & ^(1 << ((k - 1) % m.ringMask ))
	// loss
	if m.waitBits&prev > 0 {
		for n := m.ringMask - 1; n > 1; n-- {
			pk := k - uint32(n)
			ii := 1 << (pk % m.ringMask)
			if m.waitBits&byte(ii) > 0 {
				l, ok := m.Pending[pk]
				if ok {
					loss = append(loss, l)
				}
			}
		}
	}
	m.Unlock()
	m.waitCond.Broadcast()

	m.ackedMessagesMutex.Lock()
	m.ackedMessages[k] = v
	m.ackedMessagesMutex.Unlock()

	return
}

type StreamQueue struct {
	ackedSeq uint32
	msgs     [][]byte
}

func (q *StreamQueue) Push(k uint32, m []byte) (ok bool, msgs [][]byte) {
	if k <= q.ackedSeq {
		return
	}
	if k == q.ackedSeq+1 {
		ok = true
		if len(q.msgs) < 1 {
			msgs = [][]byte{m}
			q.ackedSeq = k
			return
		}
		q.push(k, m)
		msgs = q.pop()
		return
	}
	q.push(k, m)
	return
}

func (q *StreamQueue) pop() (msgs [][]byte) {
	index := len(q.msgs)
	for i, mm := range q.msgs {
		if mm == nil {
			index = i
			break
		}
	}
	msgs = q.msgs[:index]
	q.ackedSeq += uint32(index)
	if len(q.msgs) > index {
		for _, mm := range q.msgs[index:] {
			if mm != nil {
				q.msgs = q.msgs[index:]
				return
			}
		}
	}
	q.msgs = nil
	return
}

func (q *StreamQueue) push(k uint32, m []byte) {
	if q.msgs == nil {
		q.msgs = make([][]byte, 8)
	}
	index := k - q.ackedSeq - 1
	if len(q.msgs) <= int(index) {
		n := make([][]byte, index+8)
		copy(n, q.msgs)
		q.msgs = n
	}
	q.msgs[index] = m
}
