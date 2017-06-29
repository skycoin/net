package op

import (
	"github.com/skycoin/net/skycoin-messenger/msg"
	"sync"
	"github.com/skycoin/net/client"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type Reg struct {
	Network   string
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
	f := client.NewClientConnectionFactory()
	f.SetIncomingCallback(func(conn *client.ClientConnection, data []byte) bool {
		log.Printf("msg from %s In %s", conn.Key.Hex(), data)
		go c.PushLoop(conn, data)
		return true
	})
	c.SetFactory(f)
	return f.Connect(r.Network, r.Address, key)
}
