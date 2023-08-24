package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type UpdateStateField interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Field() string
	Value() []byte
}

type StateChannelCommitData interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	Address() e_common.Address
	UpdateStateFields() []UpdateStateField

	// setter()
	AddUpdateField(UpdateStateField)
}

type StateChannelCommitSign interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string
	// getter

	Address() e_common.Address
	PublicKey() common.PublicKey
	Sign() common.Sign
}

type StateChannelCommitDatas interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	CommitRoot() e_common.Hash
	Signs() []StateChannelCommitSign
	SignOfAddress(address e_common.Address) StateChannelCommitSign
	Datas() []StateChannelCommitData
	Receipts() []Receipt
	TotalGas() uint64
	// setter()
	AddSign(StateChannelCommitSign)
}
