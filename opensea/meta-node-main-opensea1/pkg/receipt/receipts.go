package receipt

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/merkle_patricia_trie"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
)

var (
	ErrorReceiptNotFound = errors.New("receipt not found")
)

type IReceipts interface {
	GetReceiptsRoot() (common.Hash, error)
	AddReceipt(IReceipt) error
	UpdateExecuteResultToReceipt(
		common.Hash,
		pb.RECEIPT_STATUS,
		[]byte,
		pb.EXCEPTION,
		uint64,
	) error
	GetReceiptsMap() map[common.Hash]IReceipt
	GetGasUsed() uint64
}

type Receipts struct {
	trie     *merkle_patricia_trie.Trie
	receipts map[common.Hash]IReceipt
}

func NewReceipts() IReceipts {
	trie := merkle_patricia_trie.New(merkle_patricia_trie.NewEmtyFullNode())
	return &Receipts{
		trie:     trie,
		receipts: make(map[common.Hash]IReceipt),
	}
}

func (r *Receipts) GetReceiptsRoot() (common.Hash, error) {
	_, hash, err := r.trie.HashRoot()
	return hash, err
}

func (r *Receipts) AddReceipt(receipt IReceipt) error {
	b, err := receipt.Marshal()
	if err != nil {
		return err
	}
	r.receipts[receipt.GetTransactionHash()] = receipt
	r.trie.Set(receipt.GetTransactionHash().Bytes(), b)
	return nil
}

func (r *Receipts) GetReceiptsMap() map[common.Hash]IReceipt {
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

func (r *Receipts) GetGasUsed() uint64 {
	gasUsed := uint64(0)
	if r.receipts == nil {
		return gasUsed
	} else {
		for _, v := range r.receipts {
			gasUsed += v.GetGasUsed()
		}
	}
	return gasUsed
}
