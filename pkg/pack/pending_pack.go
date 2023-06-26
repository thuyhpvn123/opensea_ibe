package pack

import (
	"time"

	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	t "gitlab.com/meta-node/meta-node/pkg/transaction"
)

type PendingPack struct {
	maxTransaction int
	transactions   []t.ITransaction
}

func NewPendingPack(maxTransaction int) *PendingPack {
	return &PendingPack{
		maxTransaction: maxTransaction,
	}
}

func (pp *PendingPack) AddTransaction(t t.ITransaction) bool {
	pp.transactions = append(pp.transactions, t)
	return len(pp.transactions) >= pp.maxTransaction
}

func (pp *PendingPack) GetTotalTransaction() int {
	return len(pp.transactions)
}

func (pp *PendingPack) GetPack() IPack {
	pbTransactions := make([]*pb.Transaction, len(pp.transactions))
	for i, v := range pp.transactions {
		pbTransactions[i] = v.GetProto().(*pb.Transaction)
	}
	pack := NewPack(pp.transactions, time.Now().UnixMicro())
	pack.SetHash(pack.CalculateHash())
	pack.SetAggregateSign(pack.CalculateAggregateSign())
	return pack
}
