package factory

import "sync"

type Factory interface {
	Connect(address string) (conn *Connection, err error)
	GetConns() (result []*Connection)
	ForEachConn(fn func(connection *Connection))
	Close() error
}

type FactoryCommonFields struct {
	AcceptedCallback func(connection *Connection)

	connections      map[*Connection]struct{}
	connectionsMutex sync.RWMutex

	fieldsMutex sync.RWMutex
}

func NewFactoryCommonFields() FactoryCommonFields {
	return FactoryCommonFields{connections: make(map[*Connection]struct{})}
}

func (f *FactoryCommonFields) AddConn(conn *Connection) {
	f.connectionsMutex.Lock()
	f.connections[conn] = struct{}{}
	f.connectionsMutex.Unlock()
	go func() {
		conn.WriteLoop()
		f.RemoveConn(conn)
	}()
	go conn.ReadLoop()
}

func (f *FactoryCommonFields) GetConns() (result []*Connection) {
	f.connectionsMutex.RLock()
	defer f.connectionsMutex.RUnlock()
	if len(f.connections) < 1 {
		return
	}
	result = make([]*Connection, 0, len(f.connections))
	for k := range f.connections {
		result = append(result, k)
	}
	return
}

func (f *FactoryCommonFields) ForEachConn(fn func(connection *Connection)) {
	f.connectionsMutex.RLock()
	defer f.connectionsMutex.RUnlock()
	if len(f.connections) < 1 {
		return
	}
	for k := range f.connections {
		fn(k)
	}
}

func (f *FactoryCommonFields) RemoveConn(conn *Connection) {
	f.connectionsMutex.Lock()
	delete(f.connections, conn)
	f.connectionsMutex.Unlock()
}

func (f *FactoryCommonFields) Close() (err error) {
	f.connectionsMutex.RLock()
	defer f.connectionsMutex.RUnlock()
	if len(f.connections) < 1 {
		return
	}
	for k := range f.connections {
		k.Close()
	}
	return
}