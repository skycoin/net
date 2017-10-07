package op

import (
	"sync"

	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/websocket/data"
)

type Account struct {
}

func init() {
	msg.OP_POOL[msg.OP_ACCOUNT] = &sync.Pool{
		New: func() interface{} {
			return new(Account)
		},
	}
}

func (r *Account) Execute(c msg.OPer) error {
	sc := data.GetData()
	keys := make([]string, 0, len(sc))
	for k := range sc {
		keys = append(keys, k)
	}
	c.Push(msg.OP_ACCOUNT, keys)
	return nil
}
