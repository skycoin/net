package factory

import "github.com/skycoin/net/conn"

type Connection struct {
	conn.Connection
	factory Factory
}

type Channel struct {
	Channel chan []byte //message queue, (could be a map?)

	SequenceId uint32 // id of last message, increment every new message

	LowestSequenceNumber       uint32 // lowest sequence number in map (that we still have stored)
	HighestACKedSequenceNumber uint32 // highest packet that has been ACKed

	PendingMessageCount uint32 // number of messages pending transmission
	PendingMessageBytes uint64 // number of bytes pending transmission
	LastAck             uint64 // last time an ACK of receipt was received (better to store id of highest packet id with an ACK?)

	StatusMap map[uint32]int32
	//each message, to status
	// 0 = pending transmission (written, but not transmitted yet)
	// 1 = transmitted (waiting for ACK)
	// 2 = ACKed (ACK has been received)
}

type MessageStatus struct {
	Status int
	// 0 = pending transmission (written, but not transmitted yet)
	// 1 = transmitted (waiting for ACK)
	// 2 = ACKed (ACK has been received)

	MessageSize   uint32
	InsertedAt    uint64 //unix time in ms, when it was queued for transmission
	TransmittedAt uint64 //unix time in ms, when it was transmitted, (0 if not)
	AckedAt       uint64 //unix time in ms, when ACK was received, (0 if not)
}
