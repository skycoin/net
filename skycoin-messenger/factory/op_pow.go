package factory

import (
	"github.com/pkg/errors"
	"sync"
)

func init() {
	ops[OP_POW] = &sync.Pool{
		New: func() interface{} {
			return new(workTicket)
		},
	}
}

type workTicket struct {
	Seq  uint32
	Code []byte
}

func (wt *workTicket) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if f.Proxy {
		return
	}
	pair := conn.GetTransportPair()
	if pair == nil {
		err = errors.New("GetTransportPair == nil")
		return
	}

	ok, err := pair.submitTicket(wt)
	conn.GetContextLogger().Debugf("pow ticket %#v valid %t", wt, err == nil)
	if !ok || err != nil {
		return
	}

	return
}
