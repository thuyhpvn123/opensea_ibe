package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/pack"
)

type IVerifyPackSignResultVote interface {
	GetHash() common.Hash
	GetValue() interface{}
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
	GetSign() cm.Sign
	GetPackHash() common.Hash
	Valid() bool
}

type VerifyPackSignResultVote struct {
	verifyMinerAddress   common.Address
	verifyMinerPublicKey cm.PublicKey
	verifyMinerSign      cm.Sign
	verifyResult         pack.IVerifyPackSignResult
}

func NewVerifyPackSignResultVote(
	verifyMinerAddress common.Address,
	verifyMinerPublicKey cm.PublicKey,
	verifyMinerSign cm.Sign,
	verifyResult pack.IVerifyPackSignResult,
) IVerifyPackSignResultVote {
	return &VerifyPackSignResultVote{
		verifyMinerAddress:   verifyMinerAddress,
		verifyMinerPublicKey: verifyMinerPublicKey,
		verifyMinerSign:      verifyMinerSign,
		verifyResult:         verifyResult,
	}
}

func (v *VerifyPackSignResultVote) GetHash() common.Hash {
	return v.verifyResult.GetHash()
}

func (v *VerifyPackSignResultVote) GetValue() interface{} {
	return v.verifyResult
}

func (v *VerifyPackSignResultVote) GetPublicKey() cm.PublicKey {
	return v.verifyMinerPublicKey
}

func (v *VerifyPackSignResultVote) GetAddress() common.Address {
	return v.verifyMinerAddress
}

func (v *VerifyPackSignResultVote) GetSign() cm.Sign {
	return v.verifyMinerSign
}

func (v *VerifyPackSignResultVote) GetPackHash() common.Hash {
	return v.verifyResult.GetPackHash()
}

func (v *VerifyPackSignResultVote) Valid() bool {
	return v.verifyResult.Valid()
}
