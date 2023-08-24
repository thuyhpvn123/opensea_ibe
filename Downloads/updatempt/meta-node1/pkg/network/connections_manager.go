package network

import (
	"sync"

	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/types"
	"github.com/meta-node-blockchain/meta-node/types/network"

	"github.com/ethereum/go-ethereum/common"
)

type ConnectionsManager struct {
	mu                          sync.RWMutex
	parentConnection            network.Connection
	typeToMapAddressConnections []map[common.Address]network.Connection
}

func NewConnectionsManager(
	connectionTypes []string,
) network.ConnectionsManager {
	cm := &ConnectionsManager{}
	cm.typeToMapAddressConnections = make([]map[common.Address]network.Connection, 10)
	for i := range cm.typeToMapAddressConnections {
		cm.typeToMapAddressConnections[i] = make(map[common.Address]network.Connection)
	}
	return cm
}

// getter
func (cm *ConnectionsManager) ConnectionsByType(cType int) map[common.Address]network.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.typeToMapAddressConnections[cType]
}

func (cm *ConnectionsManager) ConnectionByTypeAndAddress(cType int, address common.Address) network.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.typeToMapAddressConnections[cType][address]
}

func (cm *ConnectionsManager) ConnectionsByTypeAndAddresses(cType int, addresses []common.Address) map[common.Address]network.Connection {
	rs := make(map[common.Address]network.Connection, len(addresses))
	for _, v := range addresses {
		rs[v] = cm.ConnectionByTypeAndAddress(cType, v)
	}

	return rs
}

func (cm *ConnectionsManager) FilterAddressAvailable(cType int, addresses map[common.Address]types.StakeState) map[common.Address]types.StakeState {
	availableAddresses := make(map[common.Address]types.StakeState)
	for address := range addresses {
		if cm.ConnectionByTypeAndAddress(cType, address) != nil {
			availableAddresses[address] = addresses[address]
		}
	}
	return availableAddresses
}

func (cm *ConnectionsManager) ParentConnection() network.Connection {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.parentConnection
}

// setter
func (cm *ConnectionsManager) AddParentConnection(conn network.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.parentConnection = conn
}

func (cm *ConnectionsManager) RemoveConnection(conn network.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cType := p_common.MapTypeToIdex(conn.Type())
	if cm.typeToMapAddressConnections[cType][conn.Address()] == conn {
		delete(cm.typeToMapAddressConnections[cType], conn.Address())
	}
}

func (cm *ConnectionsManager) AddConnection(conn network.Connection, replace bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	address := conn.Address()
	cType := p_common.MapTypeToIdex(conn.Type())

	if (address != common.Address{} &&
		cm.typeToMapAddressConnections[cType][address] == nil) ||
		replace {
		cm.typeToMapAddressConnections[cType][address] = conn
	}
}

func MapAddressConnectionToInterface(data map[common.Address]network.Connection) map[common.Address]interface{} {
	rs := make(map[common.Address]interface{})
	for i, v := range data {
		rs[i] = v
	}
	return rs
}
