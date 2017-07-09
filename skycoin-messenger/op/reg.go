package op

import (
	"github.com/skycoin/net/skycoin-messenger/msg"
	"sync"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

type Reg struct {
	Address   string
	PublicKey string
}

func init() {
	msg.OP_POOL[msg.OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(Reg)
		},
	}
}

func (r *Reg) Execute(c msg.OPer) error {
	key, err := cipher.PubKeyFromHex(r.PublicKey)
	if err != nil {
		return err
	}
	f := factory.NewMessengerFactory()
	conn, err := f.Connect(r.Address)
	if err != nil {
		return err
	}
	c.SetConnection(conn)
	err = conn.Reg(key)
	if err != nil {
		return err
	}
	go c.PushLoop(conn)
	return nil
}
