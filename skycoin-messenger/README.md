# Skycoin Messenger

## Client Websocket Protocol
                       +--+--------+----------------------------------------------+
                       |  |        |                                              |
                       +-++-------++-------------------------+--------------------+
                         |        |                          |
                         v        v                          v
                      op type    seq                     json body
                       1 byte   4 byte

                       +----------------------------------------------------------+
               reg     |00|{"Network":"", "Address":"", "PublicKey":""}           |
                       +----------------------------------------------------------+
           ^
           |           +----------------------------------------------------------+
      req  |   send    |01|{"PublicKey":"", "Msg":""}                             |
           |           +----------------------------------------------------------+
           |
    +--------------------------------------------------------------------------------------+
           |
           |           +-----------+
      resp |   ack     |00|  seq   |
           |           +-----------+
           v
                       +----------------------------------------------------------+
               push    |01|{"PublicKey":"", "Msg":""}                             |
                       +----------------------------------------------------------+


## RPC Client Example

Look inside rpc/rpc_test.go

```
	client, err := rpc.DialHTTP("tcp", ":8083")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var code int
	key := cipher.PubKey([33]byte{0xf3})
	err = client.Call("Gateway.Reg", &op.Reg{PublicKey:key.Hex(), Address:":8080", Network:"tcp"}, &code)
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
```

