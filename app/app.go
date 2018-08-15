package app

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
)

//App is a wrapper on MessengerFactory, service (ssh server or client), address & keys (set by client,
//
type App struct {
	net         *factory.MessengerFactory
	service     string
	serviceAddr string
	appType     Type
	allowNodes  NodeKeys
	Version     string

	AppConnectionInitCallback func(resp *factory.AppConnResp) *factory.AppFeedback
}

type NodeKeys []string

//func() String takes a pointer to a keys and returns a string of the key
func (keys *NodeKeys) String() string {
	return fmt.Sprintf("%v", []string(*keys))
}

//func() Set takes a pointer to keys and appends a new key string to that address. Returns error.
func (keys *NodeKeys) Set(key string) error {
	*keys = append(*keys, key)
	return nil
}

type Type int

const (
	Client Type = iota
	Public
	Private
)

//NewServer takes type (public/private), service (socks/sshs), port and /version. Returns a pointer
//to the server app wrapper
func NewServer(appType Type, service, addr, version string) *App {
	messengerFactory := factory.NewMessengerFactory()
	messengerFactory.SetLoggerLevel(factory.DebugLevel)
	return &App{
		net:         messengerFactory,
		service:     service,
		serviceAddr: addr,
		appType:     appType,
		Version:     version,
	}
}

//NewServer takes type (public/private), service (socks/sshs), port and version. Returns a pointer
//to the client app wrapper
func NewClient(appType Type, service, version string) *App {
	messengerFactory := factory.NewMessengerFactory()
	messengerFactory.SetLoggerLevel(factory.DebugLevel)
	return &App{
		net:     messengerFactory,
		service: service,
		appType: appType,
		Version: version,
	}
}

//Start is a member of the App struct. It takes the port and seedpath and returns error
//Calls app.net.ConnectWithConfig to connect with the configurations specified.
func (app *App) Start(addr, scPath string) error {
	err := app.net.ConnectWithConfig(addr, &factory.ConnConfig{
		SeedConfigPath: scPath,
		OnConnected: func(connection *factory.Connection) {
			switch app.appType {
			case Public:
				connection.OfferServiceWithAddress(app.serviceAddr, app.Version, app.service)
			case Client:
				fallthrough
			case Private:
				connection.OfferPrivateServiceWithAddress(app.serviceAddr, app.Version, app.allowNodes, app.service)
			}
		},
		OnDisconnected: func(connection *factory.Connection) {
			log.Debug("exit on disconnected")
			os.Exit(1)
		},
		FindServiceNodesByAttributesCallback: app.FindServiceByAttributesCallback,
		AppConnectionInitCallback:            app.AppConnectionInitCallback,
	})
	return err
}

//Member of App struct, takes resp, returns debug info
func (app *App) FindServiceByAttributesCallback(resp *factory.QueryByAttrsResp) {
	log.Debugf("findServiceByAttributesCallback resp %#v", resp)
}

//SetAllNodes sets apps nodes
func (app *App) SetAllowNodes(nodes NodeKeys) {
	app.allowNodes = nodes
}

//Used to connect to another node with public(node) key, private(App) key and discovery key
func (app *App) ConnectTo(nodeKeyHex, appKeyHex, discoveryKeyHex string) (err error) {
	nodeKey, err := cipher.PubKeyFromHex(nodeKeyHex)
	if err != nil {
		return
	}
	appKey, err := cipher.PubKeyFromHex(appKeyHex)
	if err != nil {
		return
	}

	discoveryKey := cipher.PubKey{}
	if len(discoveryKeyHex) != 0 {
		discoveryKey, err = cipher.PubKeyFromHex(discoveryKeyHex)
		if err != nil {
			return
		}
	}
	app.net.ForEachConn(func(connection *factory.Connection) {
		connection.BuildAppConnection(nodeKey, appKey, discoveryKey)
	})
	return
}

