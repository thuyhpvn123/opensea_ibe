package state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
)

type AccountStatesManager struct {
	sync.RWMutex
	trie          *merkle_patricia_trie.Trie
	originStates  map[common.Address]types.AccountState
	dirtyStates   map[common.Address]types.AccountState
	pendingStates map[common.Address]types.AccountState
}

func NewAccountStatesManager(
	trie *merkle_patricia_trie.Trie,
) types.AccountStatesManager {
	return &AccountStatesManager{
		trie:          trie,
		originStates:  make(map[common.Address]types.AccountState),
		dirtyStates:   make(map[common.Address]types.AccountState),
		pendingStates: make(map[common.Address]types.AccountState),
	}
}

// general
func (am *AccountStatesManager) Commit() (common.Hash, error) {
	am.Lock()
	defer am.Unlock()

	hash, err := am.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	am.trie.Commit()

	am.originStates = make(map[common.Address]types.AccountState)
	am.pendingStates = make(map[common.Address]types.AccountState)
	return hash, err
}

func (am *AccountStatesManager) CommitDirty() {
	am.Lock()
	defer am.Unlock()
	for i, v := range am.dirtyStates {
		am.pendingStates[i] = v
	}
	am.dirtyStates = make(map[common.Address]types.AccountState)
}

func (am *AccountStatesManager) Revert() {
	am.Lock()
	defer am.Unlock()
	am.dirtyStates = make(map[common.Address]types.AccountState)
	am.pendingStates = make(map[common.Address]types.AccountState)
	am.trie.Revert()
}

func (am *AccountStatesManager) RevertDirty() {
	am.Lock()
	defer am.Unlock()

	am.dirtyStates = make(map[common.Address]types.AccountState)
}

// getter
func (am *AccountStatesManager) IntermediateRoot() (common.Hash, error) {
	// update account state changes to trie
	for address, v := range am.pendingStates {
		bAddress := address.Bytes()
		bData, err := v.Marshal()
		if err != nil {
			return common.Hash{}, err
		}
		am.trie.Set(bAddress, bData)
	}
	_, rootHash, err := am.trie.HashRoot()
	return rootHash, err
}

func (am *AccountStatesManager) AccountState(address common.Address) types.AccountState {
	am.RLock()
	defer am.RUnlock()
	return am.accountState(address)
}

func (am *AccountStatesManager) OriginAccountState(address common.Address) types.AccountState {
	am.Lock()
	defer am.Unlock()
	if originState, ok := am.originStates[address]; ok {
		return originState.Copy()
	}

	bData, err := am.trie.Get(address.Bytes())
	var accountState types.AccountState
	if err != nil {
		accountState = NewAccountState(address)
	} else {
		accountState = &AccountState{}
		accountState.Unmarshal(bData)
	}

	am.originStates[address] = accountState
	return am.originStates[address]
}

func (am *AccountStatesManager) Exist(address common.Address) bool {
	am.RLock()
	defer am.RUnlock()
	if _, ok := am.dirtyStates[address]; ok {
		return true
	}

	if _, ok := am.pendingStates[address]; ok {
		return true
	}

	if _, ok := am.originStates[address]; ok {
		return true
	}

	_, err := am.trie.Get(address.Bytes())
	return err == nil
}

func (am *AccountStatesManager) accountState(address common.Address) types.AccountState {
	if dirtyState, ok := am.dirtyStates[address]; ok {
		return dirtyState.Copy()
	}

	if pendingState, ok := am.pendingStates[address]; ok {
		return pendingState.Copy()
	}

	if originState, ok := am.originStates[address]; ok {
		return originState.Copy()
	}

	bData, _ := am.trie.Get(address.Bytes())
	var accountState types.AccountState
	if bData == nil {
		accountState = NewAccountState(address)
	} else {
		accountState = &AccountState{}
		err := accountState.Unmarshal(bData)
		if err != nil {
			logger.Error("error when unmarshal account state", err)
		}
	}

	am.originStates[address] = accountState
	return am.originStates[address].Copy()
}

func (am *AccountStatesManager) setState(newState types.AccountState) {
	am.dirtyStates[newState.Address()] = newState
}

func (am *AccountStatesManager) SetState(newState types.AccountState) {
	am.Lock()
	defer am.Unlock()
	am.setState(newState)
}

func (am *AccountStatesManager) PendingAccountStates() map[common.Address]types.AccountState {
	return am.pendingStates
}

func (am *AccountStatesManager) GetStorageSnapShot() storage.SnapShot {
	am.RLock()
	defer am.RUnlock()
	return am.trie.Storage().GetSnapShot()
}

func (am *AccountStatesManager) Copy() types.AccountStatesManager {
	am.RLock()
	defer am.RUnlock()
	cp := &AccountStatesManager{
		trie:          am.trie.Copy(),
		originStates:  make(map[common.Address]types.AccountState, len(am.originStates)),
		dirtyStates:   make(map[common.Address]types.AccountState, len(am.dirtyStates)),
		pendingStates: make(map[common.Address]types.AccountState, len(am.pendingStates)),
	}
	//
	for i, v := range am.originStates {
		cp.originStates[i] = v
	}

	for i, v := range am.dirtyStates {
		cp.dirtyStates[i] = v
	}

	for i, v := range am.pendingStates {
		cp.pendingStates[i] = v
	}
	return cp
}

// setter
func (am *AccountStatesManager) SetSmartContractState(address common.Address, smState types.SmartContractState) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetSmartContractState(smState)
	am.setState(as)
}

func (am *AccountStatesManager) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetNewDeviceKey(newDeviceKey)
	am.setState(as)
}

func (am *AccountStatesManager) SetLastHash(address common.Address, newLastHash common.Hash) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetLastHash(newLastHash)
	am.setState(as)
}

func (am *AccountStatesManager) AddPendingBalance(address common.Address, amount *uint256.Int) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.AddPendingBalance(amount)
	am.setState(as)
}

func (am *AccountStatesManager) SubPendingBalance(address common.Address, amount *uint256.Int) error {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	err := as.SubPendingBalance(amount)
	if err != nil {
		return err
	}
	am.setState(as)
	return nil
}

func (am *AccountStatesManager) SubBalance(address common.Address, amount *uint256.Int) error {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	err := as.SubBalance(amount)
	if err != nil {
		return err
	}
	am.setState(as)
	return nil
}

func (am *AccountStatesManager) AddBalance(address common.Address, amount *uint256.Int) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.AddBalance(amount)
	am.setState(as)
}

func (am *AccountStatesManager) SubTotalBalance(address common.Address, amount *uint256.Int) error {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	err := as.SubTotalBalance(amount)
	if err != nil {
		return err
	}
	am.setState(as)
	return nil
}

func (am *AccountStatesManager) SetCodeHash(address common.Address, hash common.Hash) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetCodeHash(hash)
	am.setState(as)
}

func (am *AccountStatesManager) SetStorageHost(address common.Address, storageHost string) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetStorageHost(storageHost)
	am.setState(as)
}

func (am *AccountStatesManager) SetStorageRoot(address common.Address, hash common.Hash) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetStorageRoot(hash)
	am.setState(as)
}

func (am *AccountStatesManager) SetLogsHash(address common.Address, hash common.Hash) {
	am.Lock()
	defer am.Unlock()
	as := am.accountState(address)
	as.SetLogsHash(hash)
	am.setState(as)
}
