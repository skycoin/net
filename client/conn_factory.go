package client

import "github.com/skycoin/skycoin/src/cipher"

type ClientConnection struct {
	key    cipher.PubKey
	client *Client
}

func (c *ClientConnection) Write(msg []byte) {
	c.client.Out <- [][]byte{c.key[:], msg}
}

type ClientConnectionFactory struct {
	client *Client
}

func NewClientConnectionFactory(client *Client) *ClientConnectionFactory {
	return &ClientConnectionFactory{client: client}
}

func (factory *ClientConnectionFactory) GetConn(key cipher.PubKey) *ClientConnection {
	return &ClientConnection{key: key, client: factory.client}
}

type ClientDirectConnectionFactory struct {
	ClientConnectionFactory
}

func NewClientDirectConnectionFactory(client *Client) *ClientDirectConnectionFactory {
	return &ClientDirectConnectionFactory{ClientConnectionFactory{client: client}}
}

func (factory *ClientDirectConnectionFactory) GetConn(key cipher.PubKey) *ClientConnection {
	return &ClientConnection{key: key, client: factory.client}
}