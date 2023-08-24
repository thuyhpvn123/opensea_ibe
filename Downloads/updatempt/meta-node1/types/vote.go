package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Vote interface {
	Value() interface{}
	Hash() e_common.Hash
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
}

type BlockVote interface {
	Proto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
	Unmarshal(
		b []byte,
		pubkey common.PublicKey,
		sign common.Sign,
	) error
	FromProto(protoreflect.ProtoMessage)

	BlockNumber() *uint256.Int
	Hash() e_common.Hash
	Value() interface{}
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
}

type VerifyPackSignResultVote interface {
	Hash() e_common.Hash
	Value() interface{}
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
	PackHash() e_common.Hash
	Valid() bool
}

type VerifyPacksSignResultVote interface {
	Hash() e_common.Hash
	Value() interface{}
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
	RequestHash() e_common.Hash
	Valid() bool
}

type VerifyTransactionSignVote interface {
	Hash() e_common.Hash
	Value() interface{}
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
	TransactionHash() e_common.Hash
	Valid() bool
}

type ExecuteResultsVote interface {
	GroupId() *uint256.Int
	Value() interface{}
	Hash() e_common.Hash
	PublicKey() common.PublicKey
	Address() e_common.Address
	Sign() common.Sign
}
