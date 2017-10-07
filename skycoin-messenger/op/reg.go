package op

import (
	"sync"
	"time"

	"errors"

	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/websocket/data"
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
	keys := data.GetData()
	if len(keys) < 1 {
		return errors.New("no public key found")
	}
	sc, ok := keys[r.PublicKey]
	if !ok {
		return errors.New("public key not found")
	}
	f := factory.NewMessengerFactory()
	_, err := f.ConnectWithConfig(r.Address, &factory.ConnConfig{
		SeedConfig:    sc,
		Reconnect:     true,
		ReconnectWait: 2 * time.Second,
		OnConnected: func(connection *factory.Connection) {
			go c.PushLoop(connection)
		},
	})
	if err != nil {
		return err
	}
	c.SetFactory(f)
	return nil
}
