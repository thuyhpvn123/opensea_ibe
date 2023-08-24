package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/types"
)

type ConnectionsManager interface {
	// getter
	ConnectionsByType(cType int) map[common.Address]Connection
	ConnectionByTypeAndAddress(cType int, address common.Address) Connection
	ConnectionsByTypeAndAddresses(cType int, addresses []common.Address) map[common.Address]Connection
	FilterAddressAvailable(cType int, addresses map[common.Address]types.StakeState) map[common.Address]types.StakeState
	ParentConnection() Connection

	// setter
	AddConnection(Connection, bool)
	RemoveConnection(Connection)
	AddParentConnection(Connection)
}
