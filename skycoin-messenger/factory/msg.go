package factory

import (
	"github.com/skycoin/skycoin/src/cipher"
)

func GenRegMsg() []byte {
	result := make([]byte, MSG_HEADER_END)
	result[MSG_OP_BEGIN] = OP_REG
	return result
}

func GenOfferServiceMsg(sub []byte) []byte {
	result := make([]byte, MSG_HEADER_END+len(sub))
	result[MSG_OP_BEGIN] = OP_OFFER_SERVICE
	copy(result[MSG_HEADER_END:], sub)
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

func GenGetServiceNodesMsg(data []byte) []byte {
	result := make([]byte, MSG_HEADER_END+len(data))
	result[MSG_OP_BEGIN] = OP_GET_SERVICE_NODES
	copy(result[MSG_HEADER_END:], data)
	return result
}
