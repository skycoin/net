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
	if conn.CreatedByTransport == nil {
		err = errors.New("CreatedByTransport == nil")
		return
	}

	ok, err := conn.CreatedByTransport.submitTicket(wt)
	conn.GetContextLogger().Debugf("pow ticket %#v valid %t", wt, err == nil)
	if !ok || err != nil {
		return
	}

	return
}
