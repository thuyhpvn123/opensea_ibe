package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/pack"
)

type IVerifyPacksSignResultVote interface {
	GetHash() common.Hash
	GetValue() interface{}
	GetPublicKey() cm.PublicKey
	GetAddress() common.Address
	GetSign() cm.Sign
	GetRequestHash() common.Hash
	Valid() bool
}

type VerifyPacksSignResultVote struct {
	address      common.Address
	publicKey    cm.PublicKey
	sign         cm.Sign
	verifyResult pack.IVerifyPacksSignResult
}

func NewVerifyPacksSignResultVote(
	address common.Address,
	publicKey cm.PublicKey,
	sign cm.Sign,
	verifyResult pack.IVerifyPacksSignResult,
) IVerifyPacksSignResultVote {
	return &VerifyPacksSignResultVote{
		address:      address,
		publicKey:    publicKey,
		sign:         sign,
		verifyResult: verifyResult,
	}
}

func (v *VerifyPacksSignResultVote) GetHash() common.Hash {
	return v.verifyResult.GetHash()
}

func (v *VerifyPacksSignResultVote) GetValue() interface{} {
	return v.verifyResult
}

func (v *VerifyPacksSignResultVote) GetPublicKey() cm.PublicKey {
	return v.publicKey
}

func (v *VerifyPacksSignResultVote) GetAddress() common.Address {
	return v.address
}

func (v *VerifyPacksSignResultVote) GetSign() cm.Sign {
	return v.sign
}

func (v *VerifyPacksSignResultVote) GetRequestHash() common.Hash {
	return v.verifyResult.GetRequestHash()
}

func (v *VerifyPacksSignResultVote) Valid() bool {
	return v.verifyResult.Valid()
}
