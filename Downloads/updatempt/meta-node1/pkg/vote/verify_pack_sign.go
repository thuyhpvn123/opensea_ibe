package vote

import (
	"github.com/ethereum/go-ethereum/common"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/types"
)

type VerifyPackSignResultVote struct {
	verifyMinerAddress   common.Address
	verifyMinerPublicKey cm.PublicKey
	verifyMinerSign      cm.Sign
	verifyResult         types.VerifyPackSignResult
}

func NewVerifyPackSignResultVote(
	verifyMinerAddress common.Address,
	verifyMinerPublicKey cm.PublicKey,
	verifyMinerSign cm.Sign,
	verifyResult types.VerifyPackSignResult,
) types.VerifyPackSignResultVote {
	return &VerifyPackSignResultVote{
		verifyMinerAddress:   verifyMinerAddress,
		verifyMinerPublicKey: verifyMinerPublicKey,
		verifyMinerSign:      verifyMinerSign,
		verifyResult:         verifyResult,
	}
}

func (v *VerifyPackSignResultVote) Hash() common.Hash {
	return v.verifyResult.Hash()
}

func (v *VerifyPackSignResultVote) Value() interface{} {
	return v.verifyResult
}

func (v *VerifyPackSignResultVote) PublicKey() cm.PublicKey {
	return v.verifyMinerPublicKey
}

func (v *VerifyPackSignResultVote) Address() common.Address {
	return v.verifyMinerAddress
}

func (v *VerifyPackSignResultVote) Sign() cm.Sign {
	return v.verifyMinerSign
}

func (v *VerifyPackSignResultVote) PackHash() common.Hash {
	return v.verifyResult.PackHash()
}

func (v *VerifyPackSignResultVote) Valid() bool {
	return v.verifyResult.Valid()
}
