package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	t "gitlab.com/meta-node/meta-node/pkg/transaction"
)

type IVerifyTransactionSignVote interface {
	GetHash() common.Hash
	GetValue() interface{}
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
	GetSign() cm.Sign
	GetTransactionHash() common.Hash
	GetValid() bool
}

type VerifyTransactionSignVote struct {
	verifyMinerAddress   common.Address
	verifyMinerPublicKey cm.PublicKey
	verifyMinerSign      cm.Sign
	verifyResult         t.IVerifyTransactionSignResult
}

func NewVerifyTransactionSignVote(
	verifyMinerAddress common.Address,
	verifyMinerPublicKey cm.PublicKey,
	verifyMinerSign cm.Sign,
	verifyResult t.IVerifyTransactionSignResult,
) IVerifyTransactionSignVote {
	return &VerifyTransactionSignVote{
		verifyMinerAddress:   verifyMinerAddress,
		verifyMinerPublicKey: verifyMinerPublicKey,
		verifyMinerSign:      verifyMinerSign,
		verifyResult:         verifyResult,
	}
}

func (v *VerifyTransactionSignVote) GetHash() common.Hash {
	return v.verifyResult.GetResultHash()
}

func (v *VerifyTransactionSignVote) GetValue() interface{} {
	return v.verifyResult
}

func (v *VerifyTransactionSignVote) GetPublicKey() cm.PublicKey {
	return v.verifyMinerPublicKey
}

func (v *VerifyTransactionSignVote) GetAddress() common.Address {
	return v.verifyMinerAddress
}

func (v *VerifyTransactionSignVote) GetSign() cm.Sign {
	return v.verifyMinerSign
}

func (v *VerifyTransactionSignVote) GetTransactionHash() common.Hash {
	return v.verifyResult.GetTransactionHash()
}

func (v *VerifyTransactionSignVote) GetValid() bool {
	return v.verifyResult.GetValid()
}
