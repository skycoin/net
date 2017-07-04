package rpc

import (
	"testing"
	"net/rpc"
	"log"
	"github.com/skycoin/net/skycoin-messenger/op"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/net/skycoin-messenger/msg"
)

func TestServeRPC(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", ":8083")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var code int
	key := cipher.PubKey([33]byte{0xf3})
	err = client.Call("Gateway.Reg", &op.Reg{PublicKey:key.Hex(), Address:":8080"}, &code)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Log("code", code)

	_ = msg.PUSH_MSG
	target := cipher.PubKey([33]byte{0xf1})
	err = client.Call("Gateway.Send", &op.Send{PublicKey:target.Hex(), Msg:"What a beautiful day!"}, &code)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Log("code", code)

	var msgs []*msg.PushMsg
	err = client.Call("Gateway.Receive", 0, &msgs)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Logf("%v", msgs)
}
