package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
	"github.com/stretchr/testify/assert"
)

var manager types.StakeStatesManager
var testAddress common.Address

func init() {
	logger.Info("Initted")
	trie := merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode(), storage.NewMemoryDb())
	manager = NewStakeStatesManager(trie)
	testAddress = common.HexToAddress("0x0000000000000000000000000000000000000001")
}

func TestMapp(t *testing.T) {
	maps := make(map[common.Address]string)
	maps[common.Address{1}] = "2"
	maps[common.Address{1}] = "4"
	logger.Info(maps[common.Address{1}])
}
func TestStakeStates(t *testing.T) {
	logger.Info("Test stake states")
	stakeStates := manager.StakeStates(testAddress)
	logger.Info(stakeStates)
}

func TestAddStakeStates(t *testing.T) {
	manager.AddStakingBalance(testAddress, NewStakeState(testAddress, uint256.NewInt(100), pb.STAKE_TYPE_VALIDATOR, ""))
}

func TestStakeStatesManager_ChangePublicConnectionAddress(t *testing.T) {
	type fields struct {
		trie          *merkle_patricia_trie.Trie
		storage       storage.Storage
		originStates  map[common.Address]types.StakeStates
		dirtyStates   map[common.Address]types.StakeStates
		pendingStates map[common.Address]types.StakeStates
	}
	type args struct {
		address  common.Address
		newState types.StakeState
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			"Test case1",
			fields{
				trie:          merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode(), storage.NewMemoryDb()),
				pendingStates: make(map[common.Address]types.StakeStates),
				originStates:  make(map[common.Address]types.StakeStates),
				dirtyStates:   make(map[common.Address]types.StakeStates),
			},
			args{
				testAddress,
				NewStakeState(testAddress, uint256.NewInt(100), pb.STAKE_TYPE_VALIDATOR, "127.0.0.1:3030"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &StakeStatesManager{
				trie:          tt.fields.trie,
				pendingStates: tt.fields.pendingStates,
				originStates:  tt.fields.originStates,
				dirtyStates:   tt.fields.dirtyStates,
			}
			ss.ChangePublicConnectionAddress(tt.args.address, tt.args.newState)
			logger.Info(ss.StakeStates(testAddress).StakeState(testAddress, pb.STAKE_TYPE_VALIDATOR).PublicConnectionAddress())
			assert.Equal(t, ss.StakeStates(testAddress).StakeState(testAddress, pb.STAKE_TYPE_VALIDATOR).PublicConnectionAddress(), "127.0.0.1:3030")
		})
	}
}

func TestStakeStatesManager_Commit(t *testing.T) {
	type fields struct {
		trie          *merkle_patricia_trie.Trie
		originStates  map[common.Address]types.StakeStates
		dirtyStates   map[common.Address]types.StakeStates
		pendingStates map[common.Address]types.StakeStates
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		struct {
			name   string
			fields fields
		}{
			"Test case1",
			fields{
				trie:          merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode(), storage.NewMemoryDb()),
				pendingStates: make(map[common.Address]types.StakeStates),
				originStates:  make(map[common.Address]types.StakeStates),
				dirtyStates:   make(map[common.Address]types.StakeStates),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &StakeStatesManager{
				trie:          tt.fields.trie,
				pendingStates: tt.fields.pendingStates,
				originStates:  tt.fields.originStates,
				dirtyStates:   tt.fields.dirtyStates,
			}
			hash, err := ss.Commit()
			assert.Nil(t, err)
			logger.Info("Commit hash", hash)
		})
	}
}
