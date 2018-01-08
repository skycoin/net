package conn

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"io"
)

type Crypto struct {
	key    cipher.PubKey
	secKey cipher.SecKey
	target cipher.PubKey
	block  cipher2.Block
	es     cipher2.Stream
	ds     cipher2.Stream
}

func NewCrypto(key cipher.PubKey, secKey cipher.SecKey) *Crypto {
	return &Crypto{
		key:    key,
		secKey: secKey,
	}
}

func (c *Crypto) SetTargetKey(target cipher.PubKey) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("SetTargetKey recovered err %v", e)
		}
	}()
	c.target = target
	ecdh := cipher.ECDH(target, c.secKey)
	c.block, err = aes.NewCipher(ecdh)
	return
}

func (c *Crypto) Encrypt(data []byte) (result []byte, err error) {
	if c.block == nil {
		err = errors.New("call SetTargetKey first")
		return
	}

	if c.es == nil {
		result = make([]byte, aes.BlockSize+len(data))
		if _, err = io.ReadFull(rand.Reader, result[:aes.BlockSize]); err != nil {
			return
		}
		c.es = cipher2.NewCFBEncrypter(c.block, result[:aes.BlockSize])
		c.es.XORKeyStream(result[aes.BlockSize:], data)
		return
	}

	c.es.XORKeyStream(data, data)
	result = data
	return
}

func (c *Crypto) Decrypt(data []byte) (result []byte, err error) {
	if c.block == nil {
		err = errors.New("call SetTargetKey first")
		return
	}

	if c.ds == nil {
		c.ds = cipher2.NewCFBDecrypter(c.block, data[:aes.BlockSize])
		data = data[aes.BlockSize:]
	}

	c.ds.XORKeyStream(data, data)
	result = data
	return
}
