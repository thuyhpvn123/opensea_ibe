package state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
)

type StakeStatesManager struct {
	sync.RWMutex
	trie          *merkle_patricia_trie.Trie
	originStates  map[common.Address]types.StakeStates
	dirtyStates   map[common.Address]types.StakeStates
	pendingStates map[common.Address]types.StakeStates
}

func NewStakeStatesManager(
	trie *merkle_patricia_trie.Trie,
) types.StakeStatesManager {
	return &StakeStatesManager{
		trie:          trie,
		originStates:  make(map[common.Address]types.StakeStates),
		dirtyStates:   make(map[common.Address]types.StakeStates),
		pendingStates: make(map[common.Address]types.StakeStates),
	}
}

func (ss *StakeStatesManager) setState(address common.Address, newState types.StakeStates) {
	ss.dirtyStates[address] = newState
}

func (ss *StakeStatesManager) SetState(address common.Address, newState types.StakeStates) {
	ss.Lock()
	defer ss.Unlock()
	ss.setState(address, newState)
}

func (ss *StakeStatesManager) AddStakingBalance(address common.Address, newState types.StakeState) {
	ss.Lock()
	defer ss.Unlock()
	state := ss.stakeStates(address)
	state.AddStakingBalance(newState)
	ss.setState(address, state)
}

func (ss *StakeStatesManager) SubStakingBalance(address common.Address, newState types.StakeState) error {
	ss.Lock()
	defer ss.Unlock()
	state := ss.stakeStates(address)
	err := state.SubStakingBalance(newState)
	if err != nil {
		return err
	}
	ss.setState(address, state)
	return nil
}

func (ss *StakeStatesManager) StakeStates(address common.Address) types.StakeStates {
	ss.Lock()
	defer ss.Unlock()
	return ss.stakeStates(address)
}

func (ss *StakeStatesManager) stakeStates(address common.Address) types.StakeStates {
	if dirtyState, ok := ss.dirtyStates[address]; ok {
		return dirtyState.Copy()
	}

	if pendingState, ok := ss.pendingStates[address]; ok {
		return pendingState.Copy()
	}

	if originState, ok := ss.originStates[address]; ok {
		return originState.Copy()
	}

	bData, _ := ss.trie.Get(address.Bytes())
	stakeStates := &StakeStates{}
	stakeStates.Unmarshal(bData)
	ss.originStates[address] = stakeStates
	return stakeStates
}

func (ss *StakeStatesManager) OriginStakeStates(address common.Address) types.StakeStates {
	ss.Lock()
	defer ss.Unlock()
	if originState, ok := ss.originStates[address]; ok {
		return originState.Copy()
	}

	bData, _ := ss.trie.Get(address.Bytes())
	stakeStates := &StakeStates{}
	stakeStates.Unmarshal(bData)
	ss.originStates[address] = stakeStates
	return stakeStates
}

func (ss *StakeStatesManager) ChangePublicConnectionAddress(address common.Address, newState types.StakeState) {
	ss.Lock()
	defer ss.Unlock()
	state := ss.stakeStates(address)
	state.ChangePublicConnectionAddress(newState)
	ss.setState(address, state)
}

func (ss *StakeStatesManager) IntermediateRoot() (common.Hash, error) {
	// update account state changes to trie
	for address, state := range ss.pendingStates {
		bAddress := address.Bytes()
		bData, err := state.Marshal()
		if err != nil {
			return common.Hash{}, err
		}
		ss.trie.Set(bAddress, bData)
	}
	_, rootHash, err := ss.trie.HashRoot()
	return rootHash, err
}

func (ss *StakeStatesManager) CommitDirty() {
	ss.Lock()
	defer ss.Unlock()
	for i, v := range ss.dirtyStates {
		ss.pendingStates[i] = v
	}
	ss.dirtyStates = make(map[common.Address]types.StakeStates)
}

func (ss *StakeStatesManager) Commit() (common.Hash, error) {
	ss.Lock()
	defer ss.Unlock()

	hash, err := ss.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	ss.trie.Commit()

	ss.originStates = make(map[common.Address]types.StakeStates)
	ss.pendingStates = make(map[common.Address]types.StakeStates)
	return hash, err
}

func (ss *StakeStatesManager) Revert() {
	ss.Lock()
	defer ss.Unlock()
	ss.dirtyStates = make(map[common.Address]types.StakeStates)
	ss.pendingStates = make(map[common.Address]types.StakeStates)
	ss.trie.Revert()
}

func (ss *StakeStatesManager) Copy() types.StakeStatesManager {
	ss.Lock()
	defer ss.Unlock()
	cp := &StakeStatesManager{
		trie:          ss.trie.Copy(),
		originStates:  make(map[common.Address]types.StakeStates, len(ss.originStates)),
		dirtyStates:   make(map[common.Address]types.StakeStates, len(ss.dirtyStates)),
		pendingStates: make(map[common.Address]types.StakeStates, len(ss.pendingStates)),
	}
	//
	for i, v := range ss.originStates {
		cp.originStates[i] = v
	}

	for i, v := range ss.dirtyStates {
		cp.dirtyStates[i] = v
	}

	for i, v := range ss.pendingStates {
		cp.pendingStates[i] = v
	}
	return cp
}

func (ss *StakeStatesManager) PendingStates() map[common.Address]types.StakeStates {
	return ss.pendingStates
}

func (ss *StakeStatesManager) GetStorageSnapShot() storage.SnapShot {
	ss.RLock()
	defer ss.RUnlock()
	return ss.trie.Storage().GetSnapShot()
}
