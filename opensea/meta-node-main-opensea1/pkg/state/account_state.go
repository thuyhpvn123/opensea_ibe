package state

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
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

type IAccountState interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetProto() protoreflect.ProtoMessage
	Copy() IAccountState
	String() string

	// getter
	GetAddress() common.Address
	GetLastHash() common.Hash
	GetBalance() *uint256.Int
	GetPendingBalance() *uint256.Int
	GetTotalBalance() *uint256.Int
	GetStakeStates(minStakeAmount *uint256.Int) map[common.Address]IStakeState
	GetStakeStatesByType(_type pb.STAKE_TYPE, minStakeAmount *uint256.Int) map[common.Address]IStakeState
	GetStakeState(address common.Address) (IStakeState, int)
	GetSmartContractState() ISmartContractState
	GetDeviceKey() common.Hash

	// setter
	SetSmartContractState(smState ISmartContractState)
	SetBalance(newBalance *uint256.Int)
	SetPendingBalance(newBalance *uint256.Int)
	SetNewDeviceKey(newDeviceKey common.Hash)
	SetLastHash(newLastHash common.Hash)
	AddPendingBalance(amount *uint256.Int)
	SubPendingBalance(amount *uint256.Int) error
	AddStakingBalance(
		address common.Address,
		stakeState IStakeState,
	)
	SubStakingBalance(
		address common.Address,
		amount *uint256.Int,
	) error
	SubBalance(amount *uint256.Int) error
	SubTotalBalance(amount *uint256.Int) error
	AddBalance(amount *uint256.Int)
	SetCodeHash(hash common.Hash)
	SetStorageHost(storageHost string)
	SetStorageRoot(hash common.Hash)
	SetLogsHash(hash common.Hash)
	AddStakeState(address common.Address, stakeState IStakeState)
}

type AccountState struct {
	proto *pb.AccountState
}

func AccountStateFromProto(proto *pb.AccountState) IAccountState {
	return &AccountState{
		proto,
	}
}

func NewAccountState(address common.Address) IAccountState {
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

func (as *AccountState) GetProto() protoreflect.ProtoMessage {
	return as.proto
}

func (as *AccountState) Copy() IAccountState {
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
			"StakeStates: %v \n"+
			"SmartContractInfo: %v \n"+
			"DeviceKey: %v \n",

		hex.EncodeToString(as.proto.Address),
		hex.EncodeToString(as.proto.LastHash),
		uint256.NewInt(0).SetBytes(as.proto.Balance),
		uint256.NewInt(0).SetBytes(as.proto.PendingBalance),
		as.GetStakeStates(uint256.NewInt(0)),
		as.GetSmartContractState(),
		hex.EncodeToString(as.proto.DeviceKey),
	)
	return str
}

// getter
func (as *AccountState) GetAddress() common.Address {
	return common.BytesToAddress(as.proto.Address)
}

func (as *AccountState) GetBalance() *uint256.Int {
	return uint256.NewInt(0).SetBytes(as.proto.Balance)
}

func (as *AccountState) GetPendingBalance() *uint256.Int {
	return uint256.NewInt(0).SetBytes(as.proto.PendingBalance)
}

func (as *AccountState) GetTotalBalance() *uint256.Int {
	return uint256.NewInt(0).Add(
		as.GetBalance(),
		as.GetPendingBalance(),
	)
}

func (as *AccountState) GetLastHash() common.Hash {
	return common.BytesToHash(as.proto.LastHash)
}

func (as *AccountState) GetStakeStates(
	minStakeAmount *uint256.Int,
) map[common.Address]IStakeState {
	stakeState := MapStakeStateFromProto(as.proto.StakeStates)
	rs := make(map[common.Address]IStakeState)
	for i, v := range stakeState {
		if v.GetAmount().Gt(minStakeAmount) || v.GetAmount().Eq(minStakeAmount) {
			rs[i] = v

		}
	}
	return rs
}

func (as *AccountState) GetStakeStatesByType(
	_type pb.STAKE_TYPE,
	minStakeAmount *uint256.Int,
) map[common.Address]IStakeState {
	allStakeState := as.GetStakeStates(minStakeAmount)
	rs := make(map[common.Address]IStakeState)
	for i, v := range allStakeState {
		if v.GetType() == _type {
			rs[i] = v
		}
	}
	return rs
}

func (as *AccountState) GetStakeState(address common.Address) (IStakeState, int) {
	if as.proto.StakeStates == nil {
		return nil, -1
	}
	for i, v := range as.proto.StakeStates {
		if common.BytesToAddress(v.Address) == address {
			return StakeStateFromProto(v), i
		}
	}
	return nil, -1
}

func (as *AccountState) GetSmartContractState() ISmartContractState {
	if as.proto.SmartContractState == nil {
		return nil
	}
	return SmartContractStateFromProto(as.proto.SmartContractState)
}

func (as *AccountState) GetDeviceKey() common.Hash {
	return common.BytesToHash(as.proto.DeviceKey)
}

// setter
func (as *AccountState) SetBalance(newBalance *uint256.Int) {
	as.proto.Balance = newBalance.Bytes()
}

func (as *AccountState) SetNewDeviceKey(newDeviceKey common.Hash) {
	as.proto.DeviceKey = newDeviceKey.Bytes()
}

func (as *AccountState) SetLastHash(newLastHash common.Hash) {
	as.proto.LastHash = newLastHash.Bytes()
}

func (as *AccountState) SetSmartContractState(smState ISmartContractState) {
	as.proto.SmartContractState = smState.GetProto().(*pb.SmartContractState)
}

func (as *AccountState) AddPendingBalance(amount *uint256.Int) {
	pendingBalance := uint256.NewInt(0).SetBytes(as.proto.PendingBalance)
	pendingBalance = pendingBalance.Add(pendingBalance, amount)
	as.proto.PendingBalance = pendingBalance.Bytes()
}

func (as *AccountState) SubPendingBalance(amount *uint256.Int) error {
	pendingBalance := as.GetPendingBalance()
	if amount.Gt(pendingBalance) {
		return ErrorInvalidSubPendingAmount
	}
	newPendingBalance := uint256.NewInt(0).Sub(pendingBalance, amount)
	as.proto.PendingBalance = newPendingBalance.Bytes()
	return nil
}

func (as *AccountState) AddStakingBalance(
	address common.Address,
	stakeState IStakeState,
) {
	if as.proto.StakeStates == nil {
		as.proto.StakeStates = make([]*pb.StakeState, 0)
	}
	currentStakeState, _ := as.GetStakeState(address)
	if currentStakeState == nil {
		as.AddStakeState(address, stakeState)
		return
	}
	currentStakeState.AddAmount(stakeState.GetAmount())
	currentStakeState.SetType(stakeState.GetType())
	currentStakeState.SetPublicConnectionAddress(stakeState.GetPublicConnectionAddress())
}

func (as *AccountState) SubStakingBalance(
	address common.Address,
	amount *uint256.Int) error {
	if as.proto.StakeStates == nil {
		return ErrorStakeStateNotFound
	}
	// stakeState, i := as.GetStakeState(address)
	stakeState, _ := as.GetStakeState(address)
	if stakeState == nil {
		return ErrorStakeStateNotFound
	}

	stakingAmount := stakeState.GetAmount()
	if amount.Gt(stakingAmount) {
		return ErrorInvalidSubStakingAmount
	}
	stakeState.SubAmount(amount)
	logger.Debug("Sub Stake state", stakeState)
	// if stakeState.GetAmount().IsZero() {
	// 	// delete stake state
	// 	as.proto.StakeStates = append(as.proto.StakeStates[:i], as.proto.StakeStates[i+1:]...)
	// }
	return nil
}

func (as *AccountState) SubBalance(amount *uint256.Int) error {
	balance := as.GetBalance()
	if amount.Gt(balance) {
		return ErrorInvalidSubBalanceAmount
	}
	newBalance := uint256.NewInt(0).Sub(balance, amount)
	as.proto.Balance = newBalance.Bytes()
	return nil
}

func (as *AccountState) SubTotalBalance(amount *uint256.Int) error {
	totalBalance := uint256.NewInt(0).Add(as.GetPendingBalance(), as.GetBalance())
	if amount.Gt(totalBalance) {
		return ErrorInvalidSubBalanceAmount
	}
	newTotalBalance := uint256.NewInt(0).Sub(totalBalance, amount)
	as.proto.PendingBalance = uint256.NewInt(0).Bytes()
	as.proto.Balance = newTotalBalance.Bytes()
	return nil
}

func (as *AccountState) AddBalance(amount *uint256.Int) {
	balance := as.GetBalance()
	newBalance := uint256.NewInt(0).Add(balance, amount)
	as.proto.Balance = newBalance.Bytes()
}

func (as *AccountState) SetCodeHash(hash common.Hash) {
	scState := as.GetSmartContractState()
	scState.SetCodeHash(hash)
}

func (as *AccountState) SetStorageHost(storageHost string) {
	scState := as.GetSmartContractState()
	scState.SetStorageHost(storageHost)
}

func (as *AccountState) SetStorageRoot(hash common.Hash) {
	scState := as.GetSmartContractState()
	scState.SetStorageRoot(hash)
}

func (as *AccountState) SetLogsHash(hash common.Hash) {
	scState := as.GetSmartContractState()
	scState.SetLogsHash(hash)
}

func (as *AccountState) SetPendingBalance(newBalance *uint256.Int) {
	as.proto.PendingBalance = newBalance.Bytes()
}

func (as *AccountState) AddStakeState(address common.Address, stakeState IStakeState) {
	if as.proto.StakeStates == nil {
		as.proto.StakeStates = make([]*pb.StakeState, 1)
	}
	as.proto.StakeStates = append(as.proto.StakeStates, stakeState.GetProto().(*pb.StakeState))
}

func SortAdressStakeStatesByAmount(stakeStates map[common.Address]IStakeState) []common.Address {
	keys := make([]common.Address, 0, len(stakeStates))
	for key := range stakeStates {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return stakeStates[keys[i]].GetAmount().Gt(
			stakeStates[keys[j]].GetAmount(),
		)
	})
	return keys
}
