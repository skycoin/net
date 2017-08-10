package factory

import (
	"time"
)

type ConnConfig struct {
	Reconnect     bool
	ReconnectWait time.Duration

	// callbacks

	// call after received resp for FindServiceNodesByKeys
	FindServiceNodesCallback func(result map[string][]string)

	// call after connected to server
	OnConnected func(connection *Connection)
}
