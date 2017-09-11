package factory

import (
	"time"
)

type ConnConfig struct {
	Reconnect     bool
	ReconnectWait time.Duration
	Creator       *MessengerFactory

	// callbacks

	FindServiceNodesByKeysCallback func(resp *QueryResp)

	FindServiceNodesByAttributesCallback func(resp *QueryByAttrsResp)

	AppConnectionInitCallback func(resp *AppConnResp)

	// call after connected to server
	OnConnected func(connection *Connection)
}
