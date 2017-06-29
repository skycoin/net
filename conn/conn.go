package conn

import "github.com/skycoin/skycoin/src/cipher"

type Connection interface {
	ReadLoop() error
	WriteLoop() error
	Write(bytes []byte) error
	WriteSlice(bytes ...[]byte) error
	GetChanIn() <-chan []byte
	Close()
	IsClosed() bool

	SendReg(key cipher.PubKey) error
	GetPublicKey() cipher.PubKey
}
