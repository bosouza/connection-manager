package connmngr

import (
	"fmt"
	"sync"
)

func CreateConnManager() ConnManager {
	return ConnManager{
		mu:      &sync.RWMutex{},
		clients: make(map[string]ConnectionFactory),
	}
}

type ConnManager struct {
	mu      *sync.RWMutex
	clients map[string]ConnectionFactory
}

func (cm *ConnManager) AddClient(client string, cnnFactory ConnectionFactory) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//don't allow re-regisring of clients
	if _, ok := cm.clients[client]; ok {
		return fmt.Errorf("client %q is already registered", client)
	}

	//register client
	cm.clients[client] = cnnFactory
	return nil
}

func (cm *ConnManager) ConnectTo(client string) (Connection, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cnnFactory, ok := cm.clients[client]
	if !ok {
		return nil, fmt.Errorf("client %q is not registered", client)
	}

	return cnnFactory.CreateConnection()
}
