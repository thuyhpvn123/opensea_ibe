package state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/merkle_patricia_trie"
	"gitlab.com/meta-node/meta-node/pkg/storage"
)

type IAccountStates interface {
	SetAccountState(accountState IAccountState) error
	IntermediateRoot() (common.Hash, error)
	Copy() IAccountStates
	Commit() (common.Hash, error)
	SetStorage(storage.IStorage)
	OpenStorage() error
	CloseStorage() error

	// getter
	GetStorageIterator() storage.IIterator
	GetTrie() *merkle_patricia_trie.Trie
	GetAccountState(address common.Address) (IAccountState, error)
	GetDirtyAccountStates() map[common.Address]IAccountState

	// update account state func
	SetSmartContractState(address common.Address, smState ISmartContractState)
	SetNewDeviceKey(address common.Address, newDeviceKey common.Hash)
	SetLastHash(address common.Address, newLastHash common.Hash)
	AddPendingBalance(address common.Address, amount *uint256.Int)
	SubPendingBalance(address common.Address, amount *uint256.Int) error
	AddStakingBalance(
		sendAddress common.Address,
		receiveAddress common.Address,
		stakeState IStakeState,
	)
	SubStakingBalance(
		sendAddress common.Address,
		receiveAddress common.Address,
		amount *uint256.Int,
	) error
	SubBalance(address common.Address, amount *uint256.Int) error
	AddBalance(address common.Address, amount *uint256.Int)
	SubTotalBalance(address common.Address, amount *uint256.Int) error
	SetCodeHash(address common.Address, hash common.Hash)
	SetStorageHost(address common.Address, storageHost string)
	SetStorageRoot(address common.Address, hash common.Hash)
	SetLogsHash(address common.Address, hash common.Hash)
}

type AccountStates struct {
	storage storage.IStorage
	trie    *merkle_patricia_trie.Trie

	loadedAccountStates map[common.Address]IAccountState
	dirtyAccountStates  map[common.Address]struct{}

	sync.RWMutex
}

func NewAccountStates(storage storage.IStorage, trie *merkle_patricia_trie.Trie) IAccountStates {
	return &AccountStates{
		storage:             storage,
		trie:                trie,
		loadedAccountStates: make(map[common.Address]IAccountState),
		dirtyAccountStates:  make(map[common.Address]struct{}),
	}
}

func (ass *AccountStates) SetAccountState(accountState IAccountState) error {
	ass.Lock()
	defer ass.Unlock()
	address := accountState.GetAddress()
	ass.loadedAccountStates[address] = accountState
	ass.dirtyAccountStates[address] = struct{}{}
	return nil
}

func (ass *AccountStates) GetAccountState(address common.Address) (IAccountState, error) {
	ass.Lock()
	defer ass.Unlock()

	if ass.loadedAccountStates[address] != nil {
		return ass.loadedAccountStates[address], nil
	}
	bAddress := address.Bytes()
	b, err := ass.trie.Get(bAddress)
	if err != nil {
		return nil, err
	}
	var accountState IAccountState
	if len(b) == 0 { // account doesn't exist => create new one
		accountState = NewAccountState(address)
	} else {
		accountState = &AccountState{}
		err = accountState.Unmarshal(b)
		if err != nil {
			return nil, err
		}
	}

	ass.loadedAccountStates[address] = accountState
	return accountState, nil
}

func (ass *AccountStates) IntermediateRoot() (common.Hash, error) {
	// update account state changes to trie
	for address, _ := range ass.dirtyAccountStates {
		bAddress := address.Bytes()
		accountState := ass.loadedAccountStates[address]
		bData, err := accountState.Marshal()
		if err != nil {
			return common.Hash{}, err
		}
		ass.trie.Set(bAddress, bData)
	}
	ass.dirtyAccountStates = make(map[common.Address]struct{})
	_, rootHash, err := ass.trie.HashRoot()
	return rootHash, err
}

func (ass *AccountStates) Copy() IAccountStates {
	ass.Lock()
	defer ass.Unlock()
	as := &AccountStates{
		trie:                ass.trie.Copy(),
		storage:             ass.storage,
		loadedAccountStates: make(map[common.Address]IAccountState, len(ass.loadedAccountStates)), // TODO
		dirtyAccountStates:  make(map[common.Address]struct{}, len(ass.dirtyAccountStates)),       // TODO
	}

	for i, v := range ass.loadedAccountStates {
		as.loadedAccountStates[i] = v
	}

	for i, v := range ass.dirtyAccountStates {
		as.dirtyAccountStates[i] = v
	}

	return as
}

func (ass *AccountStates) Commit() (common.Hash, error) {

	hash, err := ass.IntermediateRoot()
	if err != nil {
		return common.Hash{}, err
	}
	ass.trie.Commit(ass.storage)
	ass.storage.Close()

	ass.loadedAccountStates = make(map[common.Address]IAccountState)
	return hash, err
}

func (ass *AccountStates) SetStorage(s storage.IStorage) {
	ass.storage = s
}

func (ass *AccountStates) OpenStorage() error {
	return ass.storage.Open()
}

func (ass *AccountStates) CloseStorage() error {
	return ass.storage.Close()
}

// getter

func (ass *AccountStates) GetTrie() *merkle_patricia_trie.Trie {
	return ass.trie
}

func (ass *AccountStates) GetStorageIterator() storage.IIterator {
	return ass.storage.GetIterator()
}

func (ass *AccountStates) GetDirtyAccountStates() map[common.Address]IAccountState {
	rs := make(map[common.Address]IAccountState, len(ass.dirtyAccountStates))
	for i := range ass.dirtyAccountStates {
		rs[i], _ = ass.GetAccountState(i)
	}
	return rs
}

// update account state func
func (ass *AccountStates) SetSmartContractState(address common.Address, smState ISmartContractState) {
	as, _ := ass.GetAccountState(address)
	as.SetSmartContractState(smState)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetNewDeviceKey(address common.Address, newDeviceKey common.Hash) {
	as, _ := ass.GetAccountState(address)
	as.SetNewDeviceKey(newDeviceKey)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetLastHash(address common.Address, newLastHash common.Hash) {
	as, _ := ass.GetAccountState(address)
	as.SetLastHash(newLastHash)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) AddPendingBalance(address common.Address, amount *uint256.Int) {
	as, _ := ass.GetAccountState(address)
	as.AddPendingBalance(amount)
	ass.dirtyAccountStates[address] = struct{}{}

}

func (ass *AccountStates) SubPendingBalance(address common.Address, amount *uint256.Int) error {
	as, _ := ass.GetAccountState(address)
	err := as.SubPendingBalance(amount)
	if err != nil {
		return err
	}
	ass.dirtyAccountStates[address] = struct{}{}
	return nil
}

func (ass *AccountStates) AddStakingBalance(
	sendAddress common.Address,
	receiveAddress common.Address,
	stakeState IStakeState,
) {
	as, _ := ass.GetAccountState(receiveAddress)
	as.AddStakingBalance(sendAddress, stakeState)
	ass.dirtyAccountStates[receiveAddress] = struct{}{}
}

func (ass *AccountStates) SubStakingBalance(
	sendAddress common.Address,
	receiveAddress common.Address,
	amount *uint256.Int,
) error {
	as, _ := ass.GetAccountState(receiveAddress)
	err := as.SubStakingBalance(sendAddress, amount)
	if err != nil {
		logger.Warn("Error when sub staking balance", err)
		logger.Debug("SendAddress:", sendAddress)
		logger.Debug("account state:", as)
		return err
	}
	ass.dirtyAccountStates[receiveAddress] = struct{}{}
	return nil
}

func (ass *AccountStates) SubBalance(address common.Address, amount *uint256.Int) error {
	as, _ := ass.GetAccountState(address)
	err := as.SubBalance(amount)
	if err != nil {
		return err
	}
	ass.dirtyAccountStates[address] = struct{}{}
	return nil
}

func (ass *AccountStates) SubTotalBalance(address common.Address, amount *uint256.Int) error {
	as, _ := ass.GetAccountState(address)
	err := as.SubTotalBalance(amount)
	if err != nil {
		return err
	}
	ass.dirtyAccountStates[address] = struct{}{}
	return nil
}

func (ass *AccountStates) AddBalance(address common.Address, amount *uint256.Int) {
	as, _ := ass.GetAccountState(address)
	as.AddBalance(amount)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetCodeHash(address common.Address, hash common.Hash) {
	as, _ := ass.GetAccountState(address)
	as.SetCodeHash(hash)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetStorageHost(address common.Address, storageHost string) {
	as, _ := ass.GetAccountState(address)
	as.SetStorageHost(storageHost)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetStorageRoot(address common.Address, hash common.Hash) {
	as, _ := ass.GetAccountState(address)
	as.SetStorageRoot(hash)
	ass.dirtyAccountStates[address] = struct{}{}
}

func (ass *AccountStates) SetLogsHash(address common.Address, hash common.Hash) {
	as, _ := ass.GetAccountState(address)
	as.SetLogsHash(hash)
	ass.dirtyAccountStates[address] = struct{}{}
}
