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
	Keys []cipher.PubKey
}

func (query *query) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	r = &queryResp{Result: f.findServiceAddresses(query.Keys, conn.GetKey())}
	return
}

type queryResp struct {
	Result map[string][]string
}

func (resp *queryResp) Execute(conn *Connection) (err error) {
	if conn.findServiceNodesCallback != nil {
		conn.findServiceNodesCallback(resp.Result)
	}
	return
}
