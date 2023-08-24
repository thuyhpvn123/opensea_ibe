package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type LeaderSchedule interface {
	// general
	Proto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	LeaderAtSlot(slot *uint256.Int) common.Address
	ToSlot() *uint256.Int

	// setter
	SetSlots(slots map[uint256.Int]common.Address)
	SetToSlot(toSlot *uint256.Int)
}
