package smart_contract

import (
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	t "gitlab.com/meta-node/meta-node/pkg/transaction"
	"google.golang.org/protobuf/proto"
)

type IExecuteTransactions interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)

	// getter
	GetTransaction() []t.ITransaction
	GetTotalTransactions() int
	GetGroupId() *uint256.Int
}

type ExecuteTransactions struct {
	proto *pb.ExecuteTransactions
}

func NewExecuteTransactions(transactions []t.ITransaction, groupId *uint256.Int, blockNumber *uint256.Int) IExecuteTransactions {
	etPb := &pb.ExecuteTransactions{
		Transactions: make([]*pb.Transaction, len(transactions)),
		GroupId:      groupId.Bytes(),
		BlockNumber:  blockNumber.Bytes(),
	}

	for i, v := range transactions {
		etPb.Transactions[i] = v.GetProto().(*pb.Transaction)
	}
	return &ExecuteTransactions{
		proto: etPb,
	}
}

func (et *ExecuteTransactions) Unmarshal(b []byte) error {
	etPb := &pb.ExecuteTransactions{}
	err := proto.Unmarshal(b, etPb)
	if err != nil {
		return err
	}
	et.proto = etPb
	return nil
}

func (et *ExecuteTransactions) Marshal() ([]byte, error) {
	return proto.Marshal(et.proto)
}

// getter
func (et *ExecuteTransactions) GetTransaction() []t.ITransaction {
	rs := make([]t.ITransaction, len(et.proto.Transactions))
	for i, v := range et.proto.Transactions {
		rs[i] = t.TransactionFromProto(v)
	}
	return rs
}

func (et *ExecuteTransactions) GetGroupId() *uint256.Int {
	return uint256.NewInt(0).SetBytes(et.proto.GroupId)
}

func (et *ExecuteTransactions) GetTotalTransactions() int {
	return len(et.proto.Transactions)
}
