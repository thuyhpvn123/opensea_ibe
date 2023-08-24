package state

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrorInvalidSubPendingAmount      = errors.New("invalid sub pending amount")
	ErrorInvalidSubStakingAmount      = errors.New("invalid sub staking amount")
	ErrorInvalidSubBalanceAmount      = errors.New("invalid sub balance amount")
	ErrorInvalidSubTotalBalanceAmount = errors.New("invalid sub total balance amount")

	ErrorStakeStateNotFound = errors.New("stake info not found")
)

type AccountState struct {
	proto *pb.AccountState
}

func AccountStateFromProto(proto *pb.AccountState) types.AccountState {
	return &AccountState{
		proto,
	}
}

func NewAccountState(address common.Address) types.AccountState {
	return &AccountState{
		proto: &pb.AccountState{
			Address:  address.Bytes(),
			LastHash: common.Hash{}.Bytes(),
		},
	}
}

// general

func (as *AccountState) Marshal() ([]byte, error) {
	return proto.Marshal(as.proto)
}

func (as *AccountState) Unmarshal(b []byte) error {
	asProto := &pb.AccountState{}
	err := proto.Unmarshal(b, asProto)
	if err != nil {
		return err
	}
	as.proto = asProto
	return nil
}

func (as *AccountState) Proto() protoreflect.ProtoMessage {
	return as.proto
}

func (as *AccountState) Copy() types.AccountState {
	copyAs := &AccountState{
		proto: proto.Clone(as.proto).(*pb.AccountState),
	}
	return copyAs
}

func (as *AccountState) String() string {
	str := fmt.Sprintf(
		"Address: %v \n"+
			"LastHash: %v \n"+
			"Balance: %v \n"+
			"PendingBalance: %v \n"+
			"StateChannelState: %v \n"+
			"SmartContractInfo: %v \n"+
			"DeviceKey: %v \n",

		hex.EncodeToString(as.proto.Address),
		hex.EncodeToString(as.proto.LastHash),
		uint256.NewInt(0).SetBytes(as.proto.Balance),
		uint256.NewInt(0).SetBytes(as.proto.PendingBalance),
		as.StateChannelState(),
		as.SmartContractState(),
		hex.EncodeToString(as.proto.DeviceKey),
	)
	return str
}

// getter
func (as *AccountState) Address() common.Address {
	return common.BytesToAddress(as.proto.Address)
}

func (as *AccountState) Balance() *uint256.Int {
	return uint256.NewInt(0).SetBytes(as.proto.Balance)
}

func (as *AccountState) PendingBalance() *uint256.Int {
	return uint256.NewInt(0).SetBytes(as.proto.PendingBalance)
}

func (as *AccountState) TotalBalance() *uint256.Int {
	return uint256.NewInt(0).Add(
		as.Balance(),
		as.PendingBalance(),
	)
}

func (as *AccountState) LastHash() common.Hash {
	return common.BytesToHash(as.proto.LastHash)
}

func (as *AccountState) SmartContractState() types.SmartContractState {
	if as.proto.SmartContractState == nil {
		return nil
	}
	return SmartContractStateFromProto(as.proto.SmartContractState)
}

func (as *AccountState) DeviceKey() common.Hash {
	return common.BytesToHash(as.proto.DeviceKey)
}

func (as *AccountState) StateChannelState() types.StateChannelState {
	if (as.proto.StateChannelState == nil || proto.Equal(as.proto.StateChannelState, &pb.StateChannelState{})) {
		return nil
	}
	stateChannelState := &StateChannelState{}
	stateChannelState.FromProto(as.proto.StateChannelState)
	return stateChannelState
}

// setter
func (as *AccountState) SetBalance(newBalance *uint256.Int) {
	as.proto.Balance = newBalance.Bytes()
}

func (as *AccountState) SetStateChannelState(scs types.StateChannelState) {
	as.proto.StateChannelState = scs.Proto().(*pb.StateChannelState)
}

func (as *AccountState) SetNewDeviceKey(newDeviceKey common.Hash) {
	as.proto.DeviceKey = newDeviceKey.Bytes()
}

func (as *AccountState) SetLastHash(newLastHash common.Hash) {
	as.proto.LastHash = newLastHash.Bytes()
}

func (as *AccountState) SetSmartContractState(smState types.SmartContractState) {
	as.proto.SmartContractState = smState.Proto().(*pb.SmartContractState)
}

func (as *AccountState) AddPendingBalance(amount *uint256.Int) {
	pendingBalance := uint256.NewInt(0).SetBytes(as.proto.PendingBalance)
	pendingBalance = pendingBalance.Add(pendingBalance, amount)
	as.proto.PendingBalance = pendingBalance.Bytes()
}

func (as *AccountState) SubPendingBalance(amount *uint256.Int) error {
	pendingBalance := as.PendingBalance()
	if amount.Gt(pendingBalance) {
		return ErrorInvalidSubPendingAmount
	}
	newPendingBalance := uint256.NewInt(0).Sub(pendingBalance, amount)
	as.proto.PendingBalance = newPendingBalance.Bytes()
	return nil
}

func (as *AccountState) SubBalance(amount *uint256.Int) error {
	balance := as.Balance()
	if amount.Gt(balance) {
		return ErrorInvalidSubBalanceAmount
	}
	newBalance := uint256.NewInt(0).Sub(balance, amount)
	as.proto.Balance = newBalance.Bytes()
	return nil
}

func (as *AccountState) SubTotalBalance(amount *uint256.Int) error {
	totalBalance := uint256.NewInt(0).Add(as.PendingBalance(), as.Balance())
	if amount.Gt(totalBalance) {
		return ErrorInvalidSubBalanceAmount
	}
	newTotalBalance := uint256.NewInt(0).Sub(totalBalance, amount)
	as.proto.PendingBalance = uint256.NewInt(0).Bytes()
	as.proto.Balance = newTotalBalance.Bytes()
	return nil
}

func (as *AccountState) AddBalance(amount *uint256.Int) {
	balance := as.Balance()
	newBalance := uint256.NewInt(0).Add(balance, amount)
	as.proto.Balance = newBalance.Bytes()
}

func (as *AccountState) SetCodeHash(hash common.Hash) {
	scState := as.SmartContractState()
	scState.SetCodeHash(hash)
}

func (as *AccountState) SetStorageHost(storageHost string) {
	scState := as.SmartContractState()
	scState.SetStorageHost(storageHost)
}

func (as *AccountState) SetStorageRoot(hash common.Hash) {
	scState := as.SmartContractState()
	scState.SetStorageRoot(hash)
}

func (as *AccountState) SetLogsHash(hash common.Hash) {
	scState := as.SmartContractState()
	scState.SetLogsHash(hash)
}

func (as *AccountState) SetPendingBalance(newBalance *uint256.Int) {
	as.proto.PendingBalance = newBalance.Bytes()
}
