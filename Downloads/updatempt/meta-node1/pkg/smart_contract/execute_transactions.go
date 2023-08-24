package smart_contract

import (
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	t "github.com/meta-node-blockchain/meta-node/pkg/transaction"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
)

type ExecuteTransactions struct {
	proto *pb.ExecuteTransactions
}

func NewExecuteTransactions(transactions []types.Transaction, groupId *uint256.Int, blockNumber *uint256.Int) types.ExecuteTransactions {
	etPb := &pb.ExecuteTransactions{
		Transactions: make([]*pb.Transaction, len(transactions)),
		GroupId:      groupId.Bytes(),
		BlockNumber:  blockNumber.Bytes(),
	}

	for i, v := range transactions {
		etPb.Transactions[i] = v.Proto().(*pb.Transaction)
	}
	return &ExecuteTransactions{
		proto: etPb,
	}
}

// general

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
func (et *ExecuteTransactions) Transactions() []types.Transaction {
	rs := make([]types.Transaction, len(et.proto.Transactions))
	for i, v := range et.proto.Transactions {
		rs[i] = t.TransactionFromProto(v)
	}
	return rs
}

func (et *ExecuteTransactions) GroupId() *uint256.Int {
	return uint256.NewInt(0).SetBytes(et.proto.GroupId)
}

func (et *ExecuteTransactions) TotalTransactions() int {
	return len(et.proto.Transactions)
}
