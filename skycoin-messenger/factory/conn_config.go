package factory

import (
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

type ConnConfig struct {
	Reconnect     bool
	ReconnectWait time.Duration

	// callbacks

	FindServiceNodesByKeysCallback func(result map[string][]string)

	FindServiceNodesByAttributesCallback func([]cipher.PubKey)

	// call after connected to server
	OnConnected func(connection *Connection)
}
