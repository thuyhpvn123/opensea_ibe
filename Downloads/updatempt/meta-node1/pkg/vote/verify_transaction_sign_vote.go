package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/types"
)

type VerifyTransactionSignVote struct {
	verifyMinerAddress   common.Address
	verifyMinerPublicKey cm.PublicKey
	verifyMinerSign      cm.Sign
	verifyResult         types.VerifyTransactionSignResult
}

func NewVerifyTransactionSignVote(
	verifyMinerAddress common.Address,
	verifyMinerPublicKey cm.PublicKey,
	verifyMinerSign cm.Sign,
	verifyResult types.VerifyTransactionSignResult,
) types.VerifyTransactionSignVote {
	return &VerifyTransactionSignVote{
		verifyMinerAddress:   verifyMinerAddress,
		verifyMinerPublicKey: verifyMinerPublicKey,
		verifyMinerSign:      verifyMinerSign,
		verifyResult:         verifyResult,
	}
}

func (v *VerifyTransactionSignVote) Hash() common.Hash {
	return v.verifyResult.ResultHash()
}

func (v *VerifyTransactionSignVote) Value() interface{} {
	return v.verifyResult
}

func (v *VerifyTransactionSignVote) PublicKey() cm.PublicKey {
	return v.verifyMinerPublicKey
}

func (v *VerifyTransactionSignVote) Address() common.Address {
	return v.verifyMinerAddress
}

func (v *VerifyTransactionSignVote) Sign() cm.Sign {
	return v.verifyMinerSign
}

func (v *VerifyTransactionSignVote) TransactionHash() common.Hash {
	return v.verifyResult.TransactionHash()
}

func (v *VerifyTransactionSignVote) Valid() bool {
	return v.verifyResult.Valid()
}
