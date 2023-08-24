package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Block interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Hash() e_common.Hash
	Number() *uint256.Int
	LastEntryHash() e_common.Hash
	LeaderAddress() e_common.Address
	AccountStatesRoot() e_common.Hash
	ReceiptRoot() e_common.Hash
	Prevrandao() uint64
	BaseFee() uint64
	GasLimit() uint64
	Type() pb.BLOCK_TYPE
	TimeStamp() uint64
	StakeStatesRoot() e_common.Hash
	CalculateHash() (e_common.Hash, error)

	// setter
	SetTimeStamp(uint64)
}

type ConfirmBlock interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Hash() e_common.Hash
	Number() *uint256.Int
	NextLeaderAddress() e_common.Address
	ValidatorSigns() map[common.PublicKey]common.Sign
	TimeStamp() uint64
	AccountStatesRoot() e_common.Hash

	// setter
	AddValidatorSign(common.PublicKey, common.Sign)
}

type FullBlock interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Block() Block
	ValidatorSigns() map[common.PublicKey]common.Sign
	AccountStateChanges() map[e_common.Address]AccountState
	StakeStateChanges() map[e_common.Address]StakeStates
	Receipts() Receipts
	GasUsed() uint64
	Fee() *uint256.Int
	NextLeaderAddress() e_common.Address

	// setter
	SetBlock(Block)
	AddValidatorSign(common.PublicKey, common.Sign)
	SetValidatorSigns(map[common.PublicKey]common.Sign)
	SetTimeStamp(uint64)
	SetNextLeaderAddress(e_common.Address)
}
