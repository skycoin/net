package factory

import (
	"github.com/skycoin/skycoin/src/cipher"
)

func GenRegMsg() []byte {
	result := make([]byte, MSG_HEADER_END)
	result[MSG_OP_BEGIN] = OP_REG
	return result
}

func GenOfferServiceMsg(service string) []byte {
	result := make([]byte, MSG_HEADER_END+len(service))
	result[MSG_OP_BEGIN] = OP_OFFER_SERVICE
	copy(result[MSG_HEADER_END:], service)
	return result
}

func GenSendMsg(from, to cipher.PubKey, msg []byte) []byte {
	result := make([]byte, SEND_MSG_TO_PUBLIC_KEY_END+len(msg))
	result[MSG_OP_BEGIN] = OP_SEND
	copy(result[SEND_MSG_PUBLIC_KEY_BEGIN:], from[:])
	copy(result[SEND_MSG_TO_PUBLIC_KEY_BEGIN:], to[:])
	copy(result[SEND_MSG_TO_PUBLIC_KEY_END:], msg)
	return result
}

func GenCustomMsg(msg []byte) []byte {
	result := make([]byte, MSG_HEADER_END+len(msg))
	result[MSG_OP_BEGIN] = OP_CUSTOM
	copy(result[MSG_HEADER_END:], msg)
	return result
}

// Server side

func GenRegRespMsg(key cipher.PubKey) []byte {
	result := make([]byte, MSG_HEADER_END+MSG_PUBLIC_KEY_SIZE)
	result[MSG_OP_BEGIN] = OP_REG
	copy(result[MSG_HEADER_END:], key[:])
	return result
}

func GenOfferServiceRespMsg(keys []cipher.PubKey) []byte {
	result := make([]byte, MSG_HEADER_END+len(keys)*MSG_PUBLIC_KEY_SIZE)
	result[MSG_OP_BEGIN] = OP_OFFER_SERVICE
	for i, k := range keys {
		copy(result[MSG_HEADER_END+i*MSG_PUBLIC_KEY_SIZE:], k[:])
	}
	return result
}
