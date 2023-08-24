package pack

import (
	"time"

	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
)

type PendingPack struct {
	maxTransaction int
	transactions   []types.Transaction
}

func NewPendingPack(maxTransaction int) *PendingPack {
	return &PendingPack{
		maxTransaction: maxTransaction,
	}
}

func (pp *PendingPack) AddTransaction(t types.Transaction) bool {
	pp.transactions = append(pp.transactions, t)
	return len(pp.transactions) >= pp.maxTransaction
}

func (pp *PendingPack) GetTotalTransaction() int {
	return len(pp.transactions)
}

func (pp *PendingPack) GetPack() types.Pack {
	pbTransactions := make([]*pb.Transaction, len(pp.transactions))
	for i, v := range pp.transactions {
		pbTransactions[i] = v.Proto().(*pb.Transaction)
	}
	pack := NewPack(pp.transactions, time.Now().UnixMicro())
	pack.SetHash(pack.CalculateHash())
	pack.SetAggregateSign(pack.CalculateAggregateSign())
	return pack
}
