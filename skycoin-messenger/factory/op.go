package factory

import (
	"sync"
)

type simpleOP interface {
	Execute(f *MessengerFactory, conn *Connection) (r resp, err error)
}

type rawOP interface {
	RawExecute(f *MessengerFactory, conn *Connection, m []byte) (rb []byte, err error)
}

type resp interface {
	Execute(conn *Connection) (err error)
}

var (
	ops   = make([]*sync.Pool, OP_SIZE)
	resps = make([]*sync.Pool, OP_SIZE)
)

func getOP(n int) rawOP {
	if n < 0 || n > OP_SIZE {
		return nil
	}
	pool := ops[n]
	if pool == nil {
		return nil
	}
	return pool.Get().(rawOP)
}

func putOP(n int, op rawOP) {
	if n < 0 || n > OP_SIZE {
		return
	}
	pool := ops[n]
	if pool == nil {
		return
	}
	pool.Put(op)
}

func getResp(n int) resp {
	if n < 0 || n > OP_SIZE {
		return nil
	}
	pool := resps[n]
	if pool == nil {
		return nil
	}
	return pool.Get().(resp)
}

func putResp(n int, r resp) {
	if n < 0 || n > OP_SIZE {
		return
	}
	pool := resps[n]
	if pool == nil {
		return
	}
	pool.Put(r)
}
