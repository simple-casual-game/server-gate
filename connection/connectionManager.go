package connection

import "sync"

type ConnectionManager struct {
	connections map[string]*Connection
	lock        sync.RWMutex
}

var connectionManager *ConnectionManager
var lock sync.RWMutex

func GetConnectionManager() *ConnectionManager {
	lock.RLock()
	if connectionManager == nil {

		lock.RUnlock()
		lock.Lock()
		connectionManager = &ConnectionManager{
			connections: make(map[string]*Connection),
		}
		lock.Unlock()
	} else {
		defer lock.RUnlock()
	}
	return connectionManager
}

func (m *ConnectionManager) AddConnection(name string, c *Connection) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.connections[name] = c
}

func (m *ConnectionManager) GetConnection(name string) *Connection {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.connections[name]
}
