package util

import "sync"

var (
	FixedMtuPool = NewFixedSizePool(1500)
)

type FixedSizePool struct {
	pool sync.Pool
	size int
}

func NewFixedSizePool(size int) (fp *FixedSizePool) {
	fp = &FixedSizePool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
		size: size,
	}
	return fp
}

func (fp *FixedSizePool) Get() []byte {
	v := fp.pool.Get()
	return v.([]byte)
}

func (fp *FixedSizePool) Put(c []byte) {
	if len(c) != fp.size {
		if cap(c) != fp.size {
			return
		}
		c = c[:fp.size]
	}
	XorBytes(c, c, c)
	fp.pool.Put(c)
}
