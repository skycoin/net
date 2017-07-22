package rpc

import (
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/net/skycoin-messenger/op"
	"github.com/skycoin/skycoin/src/cipher"
	"net/rpc"
	"testing"
)

func TestServeRPC(t *testing.T) {
	client, err := rpc.DialHTTP("tcp", ":8083")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var code int
	err = client.Call("Gateway.Reg", &op.Reg{Address: ":8080"}, &code)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Log("code", code)

	target := cipher.PubKey([33]byte{0xf1})
	err = client.Call("Gateway.Send", &op.Send{PublicKey: target.Hex(), Msg: "What a beautiful day!"}, &code)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Log("code", code)

	var msgs []interface{}
	err = client.Call("Gateway.Receive", 0, &msgs)
	if err != nil {
		log.Fatal("calling:", err)
	}
	t.Logf("%v", msgs)
}
