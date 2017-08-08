package factory

import (
	"sync"
)

func init() {
	ops[OP_OFFER_SERVICE] = &sync.Pool{
		New: func() interface{} {
			return new(Offer)
		},
	}
}

type Offer struct {
	abstractJsonOP
	Services []*Service
}

func (offer *Offer) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	f.serviceDiscovery.register(conn, offer.Services)
	return
}
