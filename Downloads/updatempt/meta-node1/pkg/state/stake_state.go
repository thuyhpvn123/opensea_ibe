package state

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type StakeState struct {
	proto *pb.StakeState
}

func StakeStateFromProto(proto *pb.StakeState) types.StakeState {
	if proto == nil {
		return nil
	}
	return &StakeState{
		proto,
	}
}

func MapStakeStateFromProto(pbStakeState []*pb.StakeState) map[common.Address]types.StakeState {
	rs := make(map[common.Address]types.StakeState, len(pbStakeState))
	for _, v := range pbStakeState {
		ss := StakeStateFromProto(v)
		rs[ss.Address()] = ss
	}
	return rs
}

func NewStakeState(address e_common.Address, amount *uint256.Int, _type pb.STAKE_TYPE, connectionAddress string) types.StakeState {
	return &StakeState{
		proto: &pb.StakeState{
			Address:                 address.Bytes(),
			Amount:                  amount.Bytes(),
			Type:                    _type,
			PublicConnectionAddress: connectionAddress,
		},
	}
}

// general
func (ss *StakeState) Marshal() ([]byte, error) {
	return proto.Marshal(ss.proto)
}

func (ss *StakeState) Unmarshal(b []byte) error {
	ssPb := &pb.StakeState{}
	err := proto.Unmarshal(b, ssPb)
	if err != nil {
		return err
	}
	ss.proto = ssPb
	return nil
}

func (ss *StakeState) Proto() protoreflect.ProtoMessage {
	return ss.proto
}

func (ss *StakeState) Copy() types.StakeState {
	return StakeStateFromProto(proto.Clone(ss.proto).(*pb.StakeState))
}

func (ss *StakeState) String() string {
	str := fmt.Sprintf(
		"Address: %v \n"+
			"Amount: %v \n"+
			"Type: %v \n"+
			"Public Connection Address: %v \n",
		hex.EncodeToString(ss.proto.Address),
		uint256.NewInt(0).SetBytes(ss.proto.Amount),
		ss.proto.Type,
		ss.proto.PublicConnectionAddress,
	)
	return str
}

// getter

func (ss *StakeState) Address() common.Address {
	return common.BytesToAddress(ss.proto.Address)
}

func (ss *StakeState) Amount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(ss.proto.Amount)
}

func (ss *StakeState) Type() pb.STAKE_TYPE {
	return ss.proto.Type
}

func (ss *StakeState) PublicConnectionAddress() string {
	return ss.proto.PublicConnectionAddress
}

// setter
func (ss *StakeState) SetAddress(address common.Address) {
	ss.proto.Address = address.Bytes()
}

func (ss *StakeState) AddAmount(amount *uint256.Int) {
	stakingAmount := ss.Amount()
	ss.proto.Amount = uint256.NewInt(0).Add(
		stakingAmount,
		amount,
	).Bytes()
}

func (ss *StakeState) SetAmount(amount *uint256.Int) {
	ss.proto.Amount = amount.Bytes()
}

func (ss *StakeState) SubAmount(amount *uint256.Int) error {
	stakingAmount := ss.Amount()
	if amount.Gt(stakingAmount) {
		return ErrorInvalidSubStakingAmount
	}
	newAmount := uint256.NewInt(0).Sub(stakingAmount, amount)
	ss.proto.Amount = newAmount.Bytes()
	return nil
}

func (ss *StakeState) SetType(_type pb.STAKE_TYPE) {
	ss.proto.Type = _type
}

func (ss *StakeState) SetPublicConnectionAddress(str string) {
	ss.proto.PublicConnectionAddress = str
}

type StakeStates struct {
	stakeStates []types.StakeState
}

func NewStakeStates(stakeStates []types.StakeState) types.StakeStates {
	return &StakeStates{
		stakeStates: stakeStates,
	}
}

// general
func (ss *StakeStates) Marshal() ([]byte, error) {
	return proto.Marshal(ss.Proto())
}

func (ss *StakeStates) Unmarshal(b []byte) error {
	ssPb := &pb.StakeStates{}
	err := proto.Unmarshal(b, ssPb)
	if err != nil {
		return err
	}
	ss.FromProto(ssPb)
	return nil
}

func (ss *StakeStates) Proto() protoreflect.ProtoMessage {
	ssPb := &pb.StakeStates{}
	stakeStatesPb := make([]*pb.StakeState, len(ss.stakeStates))
	for i, v := range ss.stakeStates {
		stakeStatesPb[i] = v.Proto().(*pb.StakeState)
	}
	ssPb.StakeStates = stakeStatesPb
	return ssPb
}

func (ss *StakeStates) FromProto(pbMessage protoreflect.ProtoMessage) {
	ssPb := pbMessage.(*pb.StakeStates)
	ss.stakeStates = make([]types.StakeState, len(ssPb.StakeStates))
	for i, v := range ssPb.StakeStates {
		ss.stakeStates[i] = StakeStateFromProto(v)
	}
}

func (ss *StakeStates) Copy() types.StakeStates {
	copyStakeStates := make([]types.StakeState, len(ss.stakeStates))
	for i, v := range ss.stakeStates {
		copyStakeStates[i] = v.Copy()
	}
	return NewStakeStates(ss.stakeStates)
}

// getter
func (ss *StakeStates) MapStakeState(
	_type pb.STAKE_TYPE,
	maxStaker int,
	minStakeAmount *uint256.Int,
) (map[e_common.Address]types.StakeState, error) {
	rs := make(map[common.Address]types.StakeState)
	for _, v := range ss.stakeStates {
		if v.Type() == _type && !v.Amount().Lt(minStakeAmount) {
			rs[v.Address()] = v
		}
	}
	sortedAddressByAmount := SortAdressStakeStatesByAmount(rs)
	totalStaker := len(sortedAddressByAmount)
	if totalStaker > maxStaker {
		totalStaker = maxStaker
	}
	addressWithStakeStates := make(map[common.Address]types.StakeState, totalStaker)
	for i := 0; i < totalStaker; i++ {
		addressWithStakeStates[sortedAddressByAmount[i]] = rs[sortedAddressByAmount[i]]
	}
	return addressWithStakeStates, nil
}

func (ss *StakeStates) AddStakingBalance(newState types.StakeState) {
	for _, v := range ss.stakeStates {
		if v.Address() == newState.Address() && v.Type() == newState.Type() {
			v.AddAmount(newState.Amount())
			return
		}
	}
	// staking not exist, add new state to stake states
	ss.stakeStates = append(ss.stakeStates, newState)
}

func (ss *StakeStates) SubStakingBalance(newState types.StakeState) error {
	for _, v := range ss.stakeStates {
		if v.Address() == newState.Address() && v.Type() == newState.Type() {
			return v.SubAmount(newState.Amount())
		}
	}
	return fmt.Errorf("stake state not found")
}

func (ss *StakeStates) ChangePublicConnectionAddress(newState types.StakeState) {
	for _, v := range ss.stakeStates {
		if v.Address() == newState.Address() && v.Type() == newState.Type() {
			v.SetPublicConnectionAddress(newState.PublicConnectionAddress())
			return
		}
	}
	// staking not exist, add new state to stake states
	ss.stakeStates = append(ss.stakeStates, newState)
}

func (ss *StakeStates) StakeState(address common.Address, _type pb.STAKE_TYPE) types.StakeState {
	for _, v := range ss.stakeStates {
		if v.Address() == address && v.Type() == _type {
			return v
		}
	}
	return NewStakeState(address, uint256.NewInt(0), _type, "")
}

func (ss *StakeStates) String() string {
	str := fmt.Sprintf("Stake states(%v): \n", len(ss.stakeStates))
	for _, v := range ss.stakeStates {
		str += v.String()
	}
	return str
}
func SortAdressStakeStatesByAmount(stakeStates map[common.Address]types.StakeState) []common.Address {
	keys := make([]common.Address, 0, len(stakeStates))
	for key := range stakeStates {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return stakeStates[keys[i]].Amount().Gt(
			stakeStates[keys[j]].Amount(),
		)
	})
	return keys
}
