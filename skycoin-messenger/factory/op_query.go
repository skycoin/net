package factory

import (
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
)

func init() {
	ops[OP_GET_SERVICE_NODES] = &sync.Pool{
		New: func() interface{} {
			return new(query)
		},
	}
	resps[OP_GET_SERVICE_NODES] = &sync.Pool{
		New: func() interface{} {
			return new(queryResp)
		},
	}
}

type query struct {
	abstractJsonOP
	Service *Service
}

func (query *query) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if len(query.Service.Attributes) > 0 {
		r = &queryResp{PubKeys: f.findByAttributes(query.Service.Attributes)}
	} else {
		r = &queryResp{PubKeys: f.find(query.Service.Key)}
	}
	return
}

type queryResp struct {
	PubKeys []cipher.PubKey
}

func (resp *queryResp) Execute(conn *Connection) (err error) {
	conn.getServicesChan <- resp.PubKeys
	return
}
