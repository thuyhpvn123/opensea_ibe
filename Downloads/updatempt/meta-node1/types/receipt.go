package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Receipt interface {
	// general
	FromProto(proto *pb.Receipt)
	Proto() protoreflect.ProtoMessage
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	String() string

	// getter
	TransactionHash() e_common.Hash
	FromAddress() e_common.Address
	ToAddress() e_common.Address
	Amount() *uint256.Int
	Status() pb.RECEIPT_STATUS
	Action() pb.ACTION
	GasUsed() uint64
	Return() []byte

	// setter
	UpdateExecuteResult(
		status pb.RECEIPT_STATUS,
		output []byte,
		exception pb.EXCEPTION,
		gasUsed uint64,
	)
}

type Receipts interface {
	// getter
	ReceiptsRoot() (e_common.Hash, error)
	ReceiptsMap() map[e_common.Hash]Receipt
	GasUsed() uint64

	// setter
	AddReceipt(Receipt) error
	UpdateExecuteResultToReceipt(
		e_common.Hash,
		pb.RECEIPT_STATUS,
		[]byte,
		pb.EXCEPTION,
		uint64,
	) error
}
