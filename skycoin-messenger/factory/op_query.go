package factory

import (
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
)

func init() {
	ops[OP_QUERY_SERVICE_NODES] = &sync.Pool{
		New: func() interface{} {
			return new(query)
		},
	}
	resps[OP_QUERY_SERVICE_NODES] = &sync.Pool{
		New: func() interface{} {
			return new(queryResp)
		},
	}

	ops[OP_QUERY_BY_ATTRS] = &sync.Pool{
		New: func() interface{} {
			return new(queryByAttrs)
		},
	}
	resps[OP_QUERY_BY_ATTRS] = &sync.Pool{
		New: func() interface{} {
			return new(queryByAttrsResp)
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
	if conn.findServiceNodesByKeysCallback != nil {
		conn.findServiceNodesByKeysCallback(resp.Result)
	}
	return
}

// query nodes by attributes
type queryByAttrs struct {
	abstractJsonOP
	Attrs []string
}

func (query *queryByAttrs) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	r = &queryByAttrsResp{Result: f.findByAttributes(query.Attrs...)}
	return
}

type queryByAttrsResp struct {
	Result []cipher.PubKey
}

func (resp *queryByAttrsResp) Execute(conn *Connection) (err error) {
	if conn.findServiceNodesByKeysCallback != nil {
		conn.findServiceNodesByAttributesCallback(resp.Result)
	}
	return
}
