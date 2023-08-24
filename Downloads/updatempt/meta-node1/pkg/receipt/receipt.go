package receipt

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Receipt struct {
	proto *pb.Receipt
}

func NewReceipt(
	transactionHash common.Hash,
	fromAddress common.Address,
	toAddress common.Address,
	amount *uint256.Int,
	action pb.ACTION,
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gastFee uint64,
	gasUsed uint64,
) types.Receipt {
	proto := &pb.Receipt{
		TransactionHash: transactionHash.Bytes(),
		FromAddress:     fromAddress.Bytes(),
		ToAddress:       toAddress.Bytes(),
		Amount:          amount.Bytes(),
		Action:          action,
		Status:          status,
		Return:          returnValue,
		Exception:       exception,
		GasUsed:         gasUsed,
		GasFee:          gastFee,
	}
	return ReceiptFromProto(proto)
}

func ReceiptFromProto(proto *pb.Receipt) types.Receipt {
	return &Receipt{
		proto: proto,
	}
}

// general
func (r *Receipt) FromProto(proto *pb.Receipt) {
	r.proto = proto
}

func (r *Receipt) Unmarshal(b []byte) error {
	receiptPb := &pb.Receipt{}
	err := proto.Unmarshal(b, receiptPb)
	if err != nil {
		return err
	}
	r.proto = receiptPb
	return nil
}

func (r *Receipt) Marshal() ([]byte, error) {
	return proto.Marshal(r.proto)
}

func (r *Receipt) Proto() protoreflect.ProtoMessage {
	return r.proto
}

func (r *Receipt) String() string {
	str := fmt.Sprintf(`
	Transaction hash: %v
	From address: %v
	To address: %v
	Amount: %v
	Action: %v
	Status: %v
	Return: %v
	Exception: %v
	GasUsed: %v
	GasFee: %v
`,
		common.BytesToHash(r.proto.TransactionHash),
		common.BytesToAddress(r.proto.FromAddress),
		common.BytesToAddress(r.proto.ToAddress),
		uint256.NewInt(0).SetBytes(r.proto.Amount),
		r.proto.Action,
		r.proto.Status,
		common.Bytes2Hex(r.proto.Return),
		r.proto.Exception,
		r.proto.GasUsed,
		r.proto.GasFee,
	)
	return str
}

// getter
func (r *Receipt) TransactionHash() common.Hash {
	return common.BytesToHash(r.proto.TransactionHash)
}

func (r *Receipt) FromAddress() common.Address {
	return common.BytesToAddress(r.proto.FromAddress)
}

func (r *Receipt) ToAddress() common.Address {
	return common.BytesToAddress(r.proto.ToAddress)
}

func (r *Receipt) GasUsed() uint64 {
	return r.proto.GasUsed
}

func (r *Receipt) Amount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(r.proto.Amount)
}

func (r *Receipt) Return() []byte {
	return r.proto.Return
}

func (r *Receipt) Status() pb.RECEIPT_STATUS {
	return r.proto.Status
}

func (r *Receipt) Action() pb.ACTION {
	return r.proto.Action
}

// setter
func (r *Receipt) UpdateExecuteResult(
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gasUsed uint64,
) {
	r.proto.Status = status
	r.proto.Return = returnValue
	r.proto.Exception = exception
	r.proto.GasUsed = gasUsed
}
