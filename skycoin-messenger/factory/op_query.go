package factory

import (
	"sync"

	"sync/atomic"

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
			return new(QueryResp)
		},
	}

	ops[OP_QUERY_BY_ATTRS] = &sync.Pool{
		New: func() interface{} {
			return new(queryByAttrs)
		},
	}
	resps[OP_QUERY_BY_ATTRS] = &sync.Pool{
		New: func() interface{} {
			return new(QueryByAttrsResp)
		},
	}
}

var (
	querySeq uint32
)

type query struct {
	Keys []cipher.PubKey
	Seq  uint32
}

func newQuery(keys []cipher.PubKey) *query {
	q := &query{Keys: keys, Seq: atomic.AddUint32(&querySeq, 1)}
	return q
}

func (query *query) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if f.ServiceDiscoveryParent == nil {
		r = &QueryResp{
			Seq:    query.Seq,
			Result: f.findServiceAddresses(query.Keys, conn.GetKey()),
		}
		return
	}
	f.ServiceDiscoveryParent.ForEachConn(func(connection *Connection) {
		connection.setProxyConnection(query.Seq, conn)
		connection.writeOP(OP_QUERY_SERVICE_NODES, query)
	})

	return
}

type QueryResp struct {
	Seq    uint32
	Result map[string][]string
}

func (resp *QueryResp) Execute(conn *Connection) (err error) {
	if connection, ok := conn.removeProxyConnection(resp.Seq); ok {
		return connection.writeOP(OP_QUERY_SERVICE_NODES|RESP_PREFIX, resp)
	}
	if conn.findServiceNodesByKeysCallback != nil {
		conn.findServiceNodesByKeysCallback(resp)
	}
	return
}

// query nodes by attributes
type queryByAttrs struct {
	Attrs []string
	Seq   uint32
}

func newQueryByAttrs(attrs []string) *queryByAttrs {
	q := &queryByAttrs{Attrs: attrs, Seq: atomic.AddUint32(&querySeq, 1)}
	return q
}

func (query *queryByAttrs) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if f.ServiceDiscoveryParent == nil {
		r = &QueryByAttrsResp{Seq: query.Seq, Result: f.findByAttributes(query.Attrs...)}
		return
	}
	f.ServiceDiscoveryParent.ForEachConn(func(connection *Connection) {
		connection.setProxyConnection(query.Seq, conn)
		connection.writeOP(OP_QUERY_BY_ATTRS, query)
	})

	return
}

type QueryByAttrsResp struct {
	Result map[string][]cipher.PubKey
	Seq    uint32
}

func (resp *QueryByAttrsResp) Execute(conn *Connection) (err error) {
	if connection, ok := conn.removeProxyConnection(resp.Seq); ok {
		return connection.writeOP(OP_QUERY_BY_ATTRS|RESP_PREFIX, resp)
	}
	if conn.findServiceNodesByAttributesCallback != nil {
		conn.findServiceNodesByAttributesCallback(resp)
	}
	return
}
