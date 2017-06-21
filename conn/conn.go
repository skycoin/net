package conn

import "github.com/skycoin/skycoin/src/cipher"

type Connection interface {
	ReadLoop() error
	WriteLoop() error
	Write(bytes []byte) error
	WriteSlice(bytes [][]byte) error
	GetChanOut() chan<- interface{}
	IsClosed() bool

	SendReg(key cipher.PubKey) error
}
