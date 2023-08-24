package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
)

type AccountStatesManager interface {
	// general
	Commit() (common.Hash, error)
	CommitDirty()
	Revert()
	RevertDirty()
	Copy() AccountStatesManager
	// getter
	IntermediateRoot() (common.Hash, error)
	AccountState(address common.Address) AccountState
	Exist(common.Address) bool
	PendingAccountStates() map[common.Address]AccountState
	OriginAccountState(address common.Address) AccountState
	GetStorageSnapShot() storage.SnapShot
	// setter
	SetState(AccountState)
	SetSmartContractState(address common.Address, smState SmartContractState)
	SetNewDeviceKey(address common.Address, newDeviceKey common.Hash)
	SetLastHash(address common.Address, newLastHash common.Hash)
	AddPendingBalance(address common.Address, amount *uint256.Int)
	SubPendingBalance(address common.Address, amount *uint256.Int) error
	SubBalance(address common.Address, amount *uint256.Int) error
	AddBalance(address common.Address, amount *uint256.Int)
	SubTotalBalance(address common.Address, amount *uint256.Int) error
	SetCodeHash(address common.Address, hash common.Hash)
	SetStorageHost(address common.Address, storageHost string)
	SetStorageRoot(address common.Address, hash common.Hash)
	SetLogsHash(address common.Address, hash common.Hash)
}

type StakeStatesManager interface {
	AddStakingBalance(common.Address, StakeState)
	SubStakingBalance(common.Address, StakeState) error
	StakeStates(common.Address) StakeStates
	SetState(common.Address, StakeStates)
	IntermediateRoot() (common.Hash, error)
	GetStorageSnapShot() storage.SnapShot
	CommitDirty()
	Revert()
	Commit() (common.Hash, error)
	Copy() StakeStatesManager
	PendingStates() map[common.Address]StakeStates
	OriginStakeStates(address common.Address) StakeStates
}
