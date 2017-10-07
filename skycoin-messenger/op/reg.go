package op

import (
	"sync"
	"time"

	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/skycoin-messenger/msg"
)

type Reg struct {
	Address string
}

func init() {
	msg.OP_POOL[msg.OP_REG] = &sync.Pool{
		New: func() interface{} {
			return new(Reg)
		},
	}
}

func (r *Reg) Execute(c msg.OPer) error {
	f := factory.NewMessengerFactory()
	_, err := f.ConnectWithConfig(r.Address, &factory.ConnConfig{
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
