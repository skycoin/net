package factory

import "sync"

type Factory interface {
	Close() error
}

type FactoryCommonFields struct {
	AcceptedCallback  func(connection *Connection)

	connections map[*Connection]bool
	connectionsMutex sync.RWMutex
}

func NewFactoryCommonFields() FactoryCommonFields {
	return FactoryCommonFields{connections:make(map[*Connection]bool)}
}

func (f *FactoryCommonFields) AddConn(conn *Connection) {
	f.connectionsMutex.Lock()
	f.connections[conn] = true
	f.connectionsMutex.Unlock()
}

func (f *FactoryCommonFields) RemoveConn(conn *Connection) {
	f.connectionsMutex.Lock()
	delete(f.connections, conn)
	f.connectionsMutex.Unlock()
}
