package receipt

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IReceipt interface {
	GetProto() protoreflect.ProtoMessage
	GetTransactionHash() common.Hash
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	UpdateExecuteResult(
		status pb.RECEIPT_STATUS,
		output []byte,
		exception pb.EXCEPTION,
		gasUsed uint64,
	)
	GetFromAddress() common.Address
	GetToAddress() common.Address
	GetAmount() *uint256.Int
	GetReturn() []byte
	GetStatus() pb.RECEIPT_STATUS
	GetAction() pb.ACTION
	GetGasUsed() uint64
	String() string
}

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
) IReceipt {
	proto := &pb.Receipt{
		TransactionHash: transactionHash.Bytes(),
		FromAddress:     fromAddress.Bytes(),
		ToAddress:       toAddress.Bytes(),
		Amount:          amount.Bytes(),
		Action:          action,
		Status:          status,
		Return:          returnValue,
		Exception:       exception,
		GasUsed:         p_common.TRANSFER_GAS_COST,
		GasFee:          gastFee,
	}
	return ReceiptFromProto(proto)
}

func ReceiptFromProto(proto *pb.Receipt) IReceipt {
	return &Receipt{
		proto: proto,
	}
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

func (r *Receipt) GetProto() protoreflect.ProtoMessage {
	return r.proto
}

func (r *Receipt) GetTransactionHash() common.Hash {
	return common.BytesToHash(r.proto.TransactionHash)
}

func (r *Receipt) Marshal() ([]byte, error) {
	return proto.Marshal(r.proto)
}

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

func (r *Receipt) GetFromAddress() common.Address {
	return common.BytesToAddress(r.proto.FromAddress)
}

func (r *Receipt) GetToAddress() common.Address {
	return common.BytesToAddress(r.proto.ToAddress)
}

func (r *Receipt) GetGasUsed() uint64 {
	return r.proto.GasUsed
}

func (r *Receipt) GetAmount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(r.proto.Amount)
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

func (r *Receipt) GetStatus() pb.RECEIPT_STATUS {
	return r.proto.Status
}

func (r *Receipt) GetAction() pb.ACTION {
	return r.proto.Action
}
func (r *Receipt) GetReturn() []byte {
	return  (r.proto.Return)
}

