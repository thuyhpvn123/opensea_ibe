package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/types"
)

type VerifyPacksSignResultVote struct {
	address      common.Address
	publicKey    cm.PublicKey
	sign         cm.Sign
	verifyResult types.VerifyPacksSignResult
}

func NewVerifyPacksSignResultVote(
	address common.Address,
	publicKey cm.PublicKey,
	sign cm.Sign,
	verifyResult types.VerifyPacksSignResult,
) types.VerifyPacksSignResultVote {
	return &VerifyPacksSignResultVote{
		address:      address,
		publicKey:    publicKey,
		sign:         sign,
		verifyResult: verifyResult,
	}
}

func (v *VerifyPacksSignResultVote) Hash() common.Hash {
	return v.verifyResult.Hash()
}

func (v *VerifyPacksSignResultVote) Value() interface{} {
	return v.verifyResult
}

func (v *VerifyPacksSignResultVote) PublicKey() cm.PublicKey {
	return v.publicKey
}

func (v *VerifyPacksSignResultVote) Address() common.Address {
	return v.address
}

func (v *VerifyPacksSignResultVote) Sign() cm.Sign {
	return v.sign
}

func (v *VerifyPacksSignResultVote) RequestHash() common.Hash {
	return v.verifyResult.RequestHash()
}

func (v *VerifyPacksSignResultVote) Valid() bool {
	return v.verifyResult.Valid()
}
