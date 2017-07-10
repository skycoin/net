package conn

import (
	"sync"
	log "github.com/sirupsen/logrus"
	"time"
	"sync/atomic"
)

var (
	ctxId uint32
)

type Connection interface {
	ReadLoop() error
	WriteLoop() error
	Write(bytes []byte) error
	GetChanIn() <-chan []byte
	GetChanOut() chan<- []byte
	Close()
	IsClosed() bool
}

type ConnCommonFields struct {
	*PendingMap
	seq uint32 // id of last message, increment every new message

	HighestACKedSequenceNumber uint32 // highest packet that has been ACKed
	LastAck                    int64  // last time an ACK of receipt was received (better to store id of highest packet id with an ACK?)

	Status int // STATUS_CONNECTING, STATUS_CONNECTED, STATUS_ERROR
	Err    error

	closed      bool
	fieldsMutex sync.RWMutex
	writeMutex  sync.Mutex

	CTXLogger *log.Entry
}

func NewConnCommonFileds() ConnCommonFields {
	return ConnCommonFields{
		CTXLogger:  log.WithField("ctxId", atomic.AddUint32(&ctxId, 1)),
		PendingMap: NewPendingMap()}
}

func (c *ConnCommonFields) SetStatusToConnected() {
	c.fieldsMutex.Lock()
	c.Status = STATUS_CONNECTED
	c.fieldsMutex.Unlock()
}

func (c *ConnCommonFields) SetStatusToError(err error) {
	c.fieldsMutex.Lock()
	c.Status = STATUS_ERROR
	c.Err = err
	c.fieldsMutex.Unlock()
	log.Printf("SetStatusToError %v", err)
}

func (c *ConnCommonFields) UpdateLastAck(s uint32) {
	c.fieldsMutex.Lock()
	c.LastAck = time.Now().Unix()
	if s > c.HighestACKedSequenceNumber {
		c.HighestACKedSequenceNumber = s
	}
	c.fieldsMutex.Unlock()
}
