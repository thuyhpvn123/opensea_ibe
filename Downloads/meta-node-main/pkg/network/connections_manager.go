package network

import (
	"sync"

	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/state"

	"github.com/ethereum/go-ethereum/common"
)

type IConnectionsManager interface {
	AddConnection(IConnection, bool)
	RemoveConnection(IConnection)
	GetConnectionByTypeAndAddress(cType int, address common.Address) IConnection
	GetConnectionsByTypeAndAddresses(cType int, addresses []common.Address) map[common.Address]IConnection
	FilterAddressAvailable(cType int, addresses map[common.Address]state.IStakeState) map[common.Address]state.IStakeState
	AddParentConnection(IConnection)
	GetParentConnection() IConnection
}

type ConnectionsManager struct {
	mu                          sync.RWMutex
	parentConnection            IConnection
	typeToMapAddressConnections []map[common.Address]IConnection
}

func NewConnectionsManager(
	connectionTypes []string,
) IConnectionsManager {
	cm := &ConnectionsManager{}
	cm.typeToMapAddressConnections = make([]map[common.Address]IConnection, 9)
	for i := range cm.typeToMapAddressConnections {
		cm.typeToMapAddressConnections[i] = make(map[common.Address]IConnection)
	}
	return cm
}

func (cm *ConnectionsManager) FilterAddressAvailable(cType int, addresses map[common.Address]state.IStakeState) map[common.Address]state.IStakeState {
	availableAddresses := make(map[common.Address]state.IStakeState)
	for address := range addresses {
		if cm.GetConnectionByTypeAndAddress(cType, address) != nil {
			availableAddresses[address] = addresses[address]
		}
	}
	return availableAddresses
}

func (cm *ConnectionsManager) AddParentConnection(conn IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.parentConnection = conn
}

func (cm *ConnectionsManager) GetParentConnection() IConnection {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.parentConnection
}

func (cm *ConnectionsManager) AddConnection(conn IConnection, replace bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	address := conn.GetAddress()
	cType := p_common.MapTypeToIdex(conn.GetType())

	if (address != common.Address{} &&
		cm.typeToMapAddressConnections[cType][address] == nil) ||
		replace {
		cm.typeToMapAddressConnections[cType][address] = conn
	}
}

func (cm *ConnectionsManager) GetConnectionByTypeAndAddress(cType int, address common.Address) IConnection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.typeToMapAddressConnections[cType][address]
}

func (cm *ConnectionsManager) GetConnectionsByTypeAndAddresses(cType int, addresses []common.Address) map[common.Address]IConnection {
	rs := make(map[common.Address]IConnection, len(addresses))
	for _, v := range addresses {
		rs[v] = cm.GetConnectionByTypeAndAddress(cType, v)
	}

	return rs
}

func (cm *ConnectionsManager) RemoveConnection(conn IConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cType := p_common.MapTypeToIdex(conn.GetType())
	if cm.typeToMapAddressConnections[cType][conn.GetAddress()] == conn {
		delete(cm.typeToMapAddressConnections[cType], conn.GetAddress())
	}
}

func MapAddressConnectionToInterface(data map[common.Address]IConnection) map[common.Address]interface{} {
	rs := make(map[common.Address]interface{})
	for i, v := range data {
		rs[i] = v
	}
	return rs
}
