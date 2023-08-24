package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Pack interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Timestamp() uint64
	Hash() e_common.Hash

	Transactions() []Transaction
	CalculateAggregateSign() common.Sign
	SetAggregateSign(common.Sign)
	AggregateSign() common.Sign
	AggregateSignData() (pubArr [][]byte, hashArr [][]byte, sign []byte)
	ValidData() bool
	ValidSign() bool

	// setter
	CalculateHash() e_common.Hash
	SetHash(e_common.Hash)
}

type PackPool interface {
	Size() int
	AddPack(pack Pack)
	AddPacks(packs []Pack)
	TakePack(numberOfPack uint64) []Pack
	Copy() PackPool
}

type VerifyPacksSignRequest interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	AggregateSignDatas() []AggregateSignData
	Hash() e_common.Hash
}

type AggregateSignData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	PackHash() e_common.Hash
	Publickeys() [][]byte
	Hashes() [][]byte
	Sign() []byte
}

type VerifyPackSignResult interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	PackHash() e_common.Hash
	Hash() e_common.Hash
	Proto() protoreflect.ProtoMessage
	Valid() bool
}

type VerifyPacksSignResult interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	Results() []VerifyPackSignResult
	TotalPack() int
	RequestHash() e_common.Hash
	Hash() e_common.Hash
	Valid() bool
}
