package conn

import "github.com/skycoin/net/msg"

type Connection interface {
	ReadLoop() error
	Write(msg *msg.Message) error
}
