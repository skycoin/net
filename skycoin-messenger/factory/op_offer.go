package factory

import (
	"strings"
	"sync"
	"encoding/json"
)

func init() {
	ops[OP_OFFER_SERVICE] = &sync.Pool{
		New: func() interface{} {
			return new(offer)
		},
	}
}

type offer struct {
	abstractJsonOP
	Services *NodeServices
}

func (offer *offer) UnmarshalJSON(data []byte) (err error) {
	ss := &NodeServices{}
	err = json.Unmarshal(data, ss)
	if err != nil {
		return
	}
	offer.Services = ss
	return
}

func (offer *offer) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	remote := conn.GetRemoteAddr().String()
	addr := remote[:strings.LastIndex(remote, ":")]
	lastIndex := strings.LastIndex(offer.Services.ServiceAddress, ":")
	if lastIndex < 0 {
		return
	}
	addr += offer.Services.ServiceAddress[lastIndex:]
	offer.Services.ServiceAddress = addr
	f.serviceDiscovery.register(conn, offer.Services)
	return
}
