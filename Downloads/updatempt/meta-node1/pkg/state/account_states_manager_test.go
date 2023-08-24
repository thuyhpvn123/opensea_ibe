package state

import (
	"reflect"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
)

func TestAccountStatesManager_SetSmartContractState(t *testing.T) {
	type fields struct {
		RWMutex       sync.RWMutex
		trie          *merkle_patricia_trie.Trie
		originStates  map[common.Address]types.AccountState
		dirtyStates   map[common.Address]types.AccountState
		pendingStates map[common.Address]types.AccountState
	}
	type args struct {
		address common.Address
		smState types.SmartContractState
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			"Test SetSmartContractState",
			fields{
				trie:          merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode(), storage.NewMemoryDb()),
				pendingStates: make(map[common.Address]types.AccountState),
				originStates:  make(map[common.Address]types.AccountState),
				dirtyStates:   make(map[common.Address]types.AccountState),
			},
			args{
				common.HexToAddress("0x0000000000000000000000000000000000000001"),
				NewSmartContractState(nil, "127.0.0.1:3051", nil, nil, nil, nil, nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountStatesManager{
				RWMutex:       tt.fields.RWMutex,
				trie:          tt.fields.trie,
				originStates:  tt.fields.originStates,
				dirtyStates:   tt.fields.dirtyStates,
				pendingStates: tt.fields.pendingStates,
			}
			am.SetSmartContractState(tt.args.address, tt.args.smState)
		})
	}
}

func TestAccountStatesManager_Commit(t *testing.T) {
	type fields struct {
		RWMutex       sync.RWMutex
		trie          *merkle_patricia_trie.Trie
		originStates  map[common.Address]types.AccountState
		dirtyStates   map[common.Address]types.AccountState
		pendingStates map[common.Address]types.AccountState
	}
	tests := []struct {
		name    string
		fields  fields
		want    common.Hash
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountStatesManager{
				RWMutex:       tt.fields.RWMutex,
				trie:          tt.fields.trie,
				originStates:  tt.fields.originStates,
				dirtyStates:   tt.fields.dirtyStates,
				pendingStates: tt.fields.pendingStates,
			}
			got, err := am.Commit()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountStatesManager.Commit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccountStatesManager.Commit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccountStatesManager_CommitDirty(t *testing.T) {
	type fields struct {
		RWMutex       sync.RWMutex
		trie          *merkle_patricia_trie.Trie
		originStates  map[common.Address]types.AccountState
		dirtyStates   map[common.Address]types.AccountState
		pendingStates map[common.Address]types.AccountState
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AccountStatesManager{
				RWMutex:       tt.fields.RWMutex,
				trie:          tt.fields.trie,
				originStates:  tt.fields.originStates,
				dirtyStates:   tt.fields.dirtyStates,
				pendingStates: tt.fields.pendingStates,
			}
			am.CommitDirty()
		})
	}
}
