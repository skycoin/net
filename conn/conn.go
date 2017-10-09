package conn

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
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

	GetContextLogger() *log.Entry
	SetContextLogger(*log.Entry)

	GetRemoteAddr() net.Addr
	IsTCP() bool
	IsUDP() bool
}

type ConnCommonFields struct {
	seq uint32 // id of last message, increment every new message

	HighestACKedSequenceNumber uint32 // highest packet that has been ACKed
	LastAck                    int64  // last time an ACK of receipt was received (better to store id of highest packet id with an ACK?)

	lastReadTime int64

	sentBytes uint64
	receivedBytes uint64

	Status int // STATUS_CONNECTING, STATUS_CONNECTED, STATUS_ERROR
	Err    error

	In          chan []byte
	Out         chan []byte
	closed      bool
	FieldsMutex sync.RWMutex
	WriteMutex  sync.Mutex

	CTXLogger *log.Entry
}

func NewConnCommonFileds() ConnCommonFields {
	entry := log.WithField("ctxId", atomic.AddUint32(&ctxId, 1))
	return ConnCommonFields{
		lastReadTime: time.Now().Unix(),
		CTXLogger:    entry,
		In:           make(chan []byte, 1),
		Out:          make(chan []byte, 1),
	}
}

func (c *ConnCommonFields) SetStatusToConnected() {
	c.FieldsMutex.Lock()
	c.Status = STATUS_CONNECTED
	c.FieldsMutex.Unlock()
}

func (c *ConnCommonFields) SetStatusToError(err error) {
	c.FieldsMutex.Lock()
	c.Status = STATUS_ERROR
	c.Err = err
	c.FieldsMutex.Unlock()
	c.CTXLogger.Debugf("SetStatusToError %v", err)
}

func (c *ConnCommonFields) UpdateLastAck(s uint32) {
	c.FieldsMutex.Lock()
	c.LastAck = time.Now().Unix()
	if s > c.HighestACKedSequenceNumber {
		c.HighestACKedSequenceNumber = s
	}
	c.FieldsMutex.Unlock()
}

func (c *ConnCommonFields) GetContextLogger() *log.Entry {
	c.FieldsMutex.RLock()
	defer c.FieldsMutex.RUnlock()
	return c.CTXLogger
}

func (c *ConnCommonFields) SetContextLogger(l *log.Entry) {
	c.FieldsMutex.Lock()
	c.CTXLogger = l
	c.FieldsMutex.Unlock()
}

func (c *ConnCommonFields) Close() {
	c.FieldsMutex.Lock()
	defer c.FieldsMutex.Unlock()

	if c.closed {
		return
	}
	c.closed = true

	close(c.In)
	c.In = nil
	close(c.Out)
	c.Out = nil
}

func (c *ConnCommonFields) IsClosed() bool {
	c.FieldsMutex.RLock()
	defer c.FieldsMutex.RUnlock()
	return c.closed
}

func (c *ConnCommonFields) GetLastTime() int64 {
	return atomic.LoadInt64(&c.lastReadTime)
}

func (c *ConnCommonFields) UpdateLastTime() {
	atomic.StoreInt64(&c.lastReadTime, time.Now().Unix())
}

func (c *ConnCommonFields) GetSentBytes() uint64 {
	return atomic.LoadUint64(&c.sentBytes)
}

func (c *ConnCommonFields) AddSentBytes(n int) {
	atomic.AddUint64(&c.sentBytes, uint64(n))
}

func (c *ConnCommonFields) GetReceivedBytes() uint64 {
	return atomic.LoadUint64(&c.receivedBytes)
}

func (c *ConnCommonFields) AddReceivedBytes(n int) {
	atomic.AddUint64(&c.receivedBytes, uint64(n))
}