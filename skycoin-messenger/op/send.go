package op

import (
	"sync"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/net/util"
)

type Send struct {
	PublicKey string
	Msg       string
}

func init() {
	msg.OP_POOL[msg.OP_SEND] = &sync.Pool{
		New: func() interface{} {
			return new(Send)
		},
	}
}

func (s *Send) Execute(c msg.OPer) error {
	key, err := cipher.PubKeyFromHex(s.PublicKey)
	if err != nil {
		return err
	}
	connection := c.GetFactory().Dial(key)
	connection.Out <- util.String2ByteSlice(s.Msg)
	return nil
}
