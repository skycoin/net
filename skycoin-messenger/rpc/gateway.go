package rpc

import (
	"github.com/skycoin/net/skycoin-messenger/op"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"fmt"
)

type Gateway struct {
}

func (g *Gateway) Reg(op *op.Reg, result *int) error {
	return op.Execute(DefaultClient)
}

func (g *Gateway) Send(op *op.Send, result *int) error {
	return op.Execute(DefaultClient)
}

func (g *Gateway) Receive(option int, msgs *[]*msg.PushMsg) error {
	for {
		select {
		case m, ok := <-DefaultClient.push:
			if !ok {
				return nil
			}
			switch m := m.(type) {
			case *msg.PushMsg:
				*msgs = append(*msgs, m)
			default:
				return fmt.Errorf("recv not implemented msg %#v", m)
			}
		default:
			return nil
		}
	}
	return nil
}
