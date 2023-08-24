package receipt

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
)

var (
	ErrorReceiptNotFound = errors.New("receipt not found")
)

type Receipts struct {
	trie     *merkle_patricia_trie.Trie
	receipts map[common.Hash]types.Receipt
}

func NewReceipts() types.Receipts {
	trie := merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode(), storage.NewMemoryDb())
	return &Receipts{
		trie:     trie,
		receipts: make(map[common.Hash]types.Receipt),
	}
}

func (r *Receipts) ReceiptsRoot() (common.Hash, error) {
	_, hash, err := r.trie.HashRoot()
	return hash, err
}

func (r *Receipts) AddReceipt(receipt types.Receipt) error {
	b, err := receipt.Marshal()
	if err != nil {
		return err
	}
	r.receipts[receipt.TransactionHash()] = receipt
	r.trie.Set(receipt.TransactionHash().Bytes(), b)
	return nil
}

func (r *Receipts) ReceiptsMap() map[common.Hash]types.Receipt {
	return r.receipts
}

func (r *Receipts) UpdateExecuteResultToReceipt(
	hash common.Hash,
	status pb.RECEIPT_STATUS,
	returnValue []byte,
	exception pb.EXCEPTION,
	gasUsed uint64,
) error {
	receipt := r.receipts[hash]
	if receipt == nil {
		return ErrorReceiptNotFound
	}
	receipt.UpdateExecuteResult(
		status,
		returnValue,
		exception,
		gasUsed,
	)
	err := r.AddReceipt(receipt)
	return err
}

func (r *Receipts) GasUsed() uint64 {
	gasUsed := uint64(0)
	if r.receipts == nil {
		return gasUsed
	} else {
		for _, v := range r.receipts {
			gasUsed += v.GasUsed()
		}
	}
	return gasUsed
}
