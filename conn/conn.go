package conn

import (
	"sync"
	"log"
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

	LowestSequenceNumber       uint32 // lowest sequence number in map (that we still have stored)
	HighestACKedSequenceNumber uint32 // highest packet that has been ACKed

	PendingMessageCount uint32 // number of messages pending transmission
	PendingMessageBytes uint64 // number of bytes pending transmission
	LastAck             uint64 // last time an ACK of receipt was received (better to store id of highest packet id with an ACK?)

	Status int //<-- status may be a struct and include log
	// status of connection
	// connecting (setting up)
	// connected (is connected)
	// disconnected by your side, disconnected by other side, timeout, error etc
	Err error

	closed      bool
	fieldsMutex sync.RWMutex
	writeMutex  sync.Mutex
}

func NewConnCommonFileds() *ConnCommonFields {
	return &ConnCommonFields{PendingMap: NewPendingMap()}
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
