package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract"
)

type ExecuteResultsVote struct {
	executeResults *smart_contract.ExecuteResults
	sign           cm.Sign
	pubkey         cm.PublicKey
}

func NewExecuteResultsVote(
	executeResults *smart_contract.ExecuteResults,
	sign cm.Sign,
	pubkey cm.PublicKey,
) *ExecuteResultsVote {
	return &ExecuteResultsVote{
		executeResults: executeResults,
		sign:           sign,
		pubkey:         pubkey,
	}
}

func (v *ExecuteResultsVote) GetGroupId() *uint256.Int {
	return v.executeResults.GetGroupId()
}

func (v *ExecuteResultsVote) GetValue() interface{} {
	return v.executeResults
}

func (v *ExecuteResultsVote) GetHash() common.Hash {
	return v.executeResults.GetHash()
}

func (v *ExecuteResultsVote) GetPublicKey() cm.PublicKey {
	return v.pubkey
}

func (v *ExecuteResultsVote) GetAddress() common.Address {
	return cm.GetAddressFromPubkey(v.pubkey)
}

func (v *ExecuteResultsVote) GetSign() cm.Sign {
	return v.sign
}
