package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/smart_contract"
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

func (v *ExecuteResultsVote) GroupId() *uint256.Int {
	return v.executeResults.GroupId()
}

func (v *ExecuteResultsVote) Value() interface{} {
	return v.executeResults
}

func (v *ExecuteResultsVote) Hash() common.Hash {
	return v.executeResults.Hash()
}

func (v *ExecuteResultsVote) PublicKey() cm.PublicKey {
	return v.pubkey
}

func (v *ExecuteResultsVote) Address() common.Address {
	return cm.AddressFromPubkey(v.pubkey)
}

func (v *ExecuteResultsVote) Sign() cm.Sign {
	return v.sign
}
