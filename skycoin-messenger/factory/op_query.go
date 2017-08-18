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
	result, ok := f.findServiceAddresses(query.Keys, conn.GetKey())
	if ok {
		r = &queryResp{Result: result}
	}
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
	result, ok := f.findByAttributes(query.Attrs...)
	if ok {
		r = &queryByAttrsResp{Result: result}
	}
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
