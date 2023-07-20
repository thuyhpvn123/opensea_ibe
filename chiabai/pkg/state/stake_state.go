package state

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IStakeState interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	GetProto() protoreflect.ProtoMessage
	Copy() IStakeState
	String() string

	// getter
	GetAddress() common.Address
	GetAmount() *uint256.Int
	GetType() pb.STAKE_TYPE
	GetPublicConnectionAddress() string

	// setter
	AddAmount(*uint256.Int)
	SubAmount(*uint256.Int) error
	SetAmount(*uint256.Int)
	SetType(pb.STAKE_TYPE)
	SetPublicConnectionAddress(string)
}

type StakeState struct {
	proto *pb.StakeState
}

func StakeStateFromProto(proto *pb.StakeState) IStakeState {
	if proto == nil {
		return nil
	}
	return &StakeState{
		proto,
	}
}

func MapStakeStateFromProto(pbStakeState []*pb.StakeState) map[common.Address]IStakeState {
	rs := make(map[common.Address]IStakeState, len(pbStakeState))
	for _, v := range pbStakeState {
		ss := StakeStateFromProto(v)
		rs[ss.GetAddress()] = ss
	}
	return rs
}

func NewStakeState(amount *uint256.Int, _type pb.STAKE_TYPE, connectionAddress string) IStakeState {
	return &StakeState{
		proto: &pb.StakeState{
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

func (ss *StakeState) GetAddress() common.Address {
	return common.BytesToAddress(ss.proto.Address)
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

func (ss *StakeState) GetProto() protoreflect.ProtoMessage {
	return ss.proto
}

func (ss *StakeState) Copy() IStakeState {
	return StakeStateFromProto(proto.Clone(ss.proto).(*pb.StakeState))
}

func (ss *StakeState) String() string {
	str := fmt.Sprintf(
		"Amount: %v \n"+
			"Type: %v \n"+
			"Public Connection Address: %v \n",
		uint256.NewInt(0).SetBytes(ss.proto.Amount),
		ss.proto.Type,
		ss.proto.PublicConnectionAddress,
	)
	return str
}

// getter
func (ss *StakeState) GetAmount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(ss.proto.Amount)
}

func (ss *StakeState) GetType() pb.STAKE_TYPE {
	return ss.proto.Type
}

func (ss *StakeState) GetPublicConnectionAddress() string {
	return ss.proto.PublicConnectionAddress
}

// setter
func (ss *StakeState) AddAmount(amount *uint256.Int) {
	stakingAmount := ss.GetAmount()
	ss.proto.Amount = uint256.NewInt(0).Add(
		stakingAmount,
		amount,
	).Bytes()
}

func (ss *StakeState) SetAmount(amount *uint256.Int) {
	ss.proto.Amount = amount.Bytes()
}

func (ss *StakeState) SubAmount(amount *uint256.Int) error {
	stakingAmount := ss.GetAmount()
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
