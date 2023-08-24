package pos

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
)

var (
	ErrNotExistsInAddresses = errors.New("not exist in addresses")
	ErrAlreadyVoted         = errors.New("already voted")
	ErrInvalidSign          = errors.New("invalid sign")
)

// vote pool using stake weight
type VotePool struct {
	approveRateMul uint64
	approveRateDiv uint64

	requireStakedAmount *uint256.Int
	addresses           map[common.Address]*uint256.Int          // use to track participators and amount
	votes               map[common.Hash]map[cm.PublicKey]cm.Sign // vote hash => addresses
	mapAddressVote      map[common.Address]common.Hash
	voteValues          map[common.Hash]interface{}
	result              *common.Hash

	finished bool
	voteMu   sync.RWMutex
}

func NewVotePool(
	stakeStates types.StakeStates,
	_type pb.STAKE_TYPE,
	maxStaker int,
	minStakeAmount *uint256.Int,
	approveRateMul uint64,
	approveRateDiv uint64,
) *VotePool {
	addressWithStakeState, _ := stakeStates.MapStakeState(
		_type,
		maxStaker,
		minStakeAmount,
	)
	addresses := make(map[common.Address]*uint256.Int, len(addressWithStakeState))
	totalStaked := uint256.NewInt(0)
	for address, info := range addressWithStakeState {
		amount := info.Amount()
		totalStaked = totalStaked.Add(totalStaked, amount)
		addresses[address] = amount
	}
	// zero staked flow
	if totalStaked.IsZero() {
		for k := range addresses {
			addresses[k] = uint256.NewInt(1) // change to 1 instead of 0 for equal calculate require
		}
		totalStaked = uint256.NewInt(uint64(len(addresses)))
	}

	requireStakedAmount := totalStaked
	requireStakedAmount = requireStakedAmount.Mul(requireStakedAmount, uint256.NewInt(approveRateMul))
	requireStakedAmount = requireStakedAmount.Div(requireStakedAmount, uint256.NewInt(approveRateDiv))

	return &VotePool{
		requireStakedAmount: requireStakedAmount,
		addresses:           addresses,
		approveRateMul:      approveRateMul,
		approveRateDiv:      approveRateDiv,
		votes:               make(map[common.Hash]map[cm.PublicKey]cm.Sign),
		mapAddressVote:      make(map[common.Address]common.Hash),
		voteValues:          make(map[common.Hash]interface{}),
		result:              nil,
	}
}
func (p *VotePool) AddVote(vote types.Vote) error {
	p.voteMu.Lock()
	defer p.voteMu.Unlock()
	pubkey := vote.PublicKey()
	sign := vote.Sign()
	hash := vote.Hash()
	value := vote.Value()
	address := vote.Address()

	if !bls.VerifySign(pubkey, sign, hash.Bytes()) {
		return ErrInvalidSign
	}
	if v, ok := p.addresses[address]; !ok || v == nil {
		return ErrNotExistsInAddresses
	}

	if _, ok := p.mapAddressVote[address]; ok {
		return ErrAlreadyVoted
	}

	p.mapAddressVote[address] = hash
	if p.votes[hash] == nil {
		p.votes[hash] = make(map[cm.PublicKey]cm.Sign)
	}
	p.votes[hash][pubkey] = sign

	if value != nil {
		p.voteValues[hash] = value
	}

	p.checkVote(hash)
	return nil
}

func (p *VotePool) AddVoteValue(voteHash common.Hash, voteValue interface{}) {
	p.voteMu.Lock()
	defer p.voteMu.Unlock()
	p.voteValues[voteHash] = voteValue
}

func (p *VotePool) AddressVote(address common.Address) common.Hash {
	p.voteMu.Lock()
	defer p.voteMu.Unlock()
	return p.mapAddressVote[address]
}

func (p *VotePool) checkVote(voteHash common.Hash) {
	totalStakedForVote := uint256.NewInt(0)
	//
	for k := range p.votes[voteHash] {
		stakedAmount := p.addresses[cm.AddressFromPubkey(k)]
		totalStakedForVote = totalStakedForVote.Add(totalStakedForVote, stakedAmount)
	}
	logger.Debug("totalStakedForVote / require", totalStakedForVote, p.requireStakedAmount)
	if totalStakedForVote.Gt(p.requireStakedAmount) || totalStakedForVote.Eq(p.requireStakedAmount) {
		p.result = &voteHash
	}
}

func (p *VotePool) Addresses() map[common.Address]*uint256.Int {
	return p.addresses
}

func (p *VotePool) Result() *common.Hash {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.result
}

func (p *VotePool) ResultValue() interface{} {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.voteValues[*p.result]
}

func (p *VotePool) WrongVoteAddresses() []common.Address {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()

	rs := []common.Address{}
	for k := range p.addresses {
		if p.mapAddressVote[k] != *p.result {
			rs = append(rs, k)
		}
	}
	return rs
}

func (p *VotePool) Signs(voteHash common.Hash) map[cm.PublicKey]cm.Sign {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.votes[voteHash]
}

func (p *VotePool) SetFinished(finished bool) {
	p.voteMu.Lock()
	defer p.voteMu.Unlock()
	p.finished = finished
}

func (p *VotePool) Finished() bool {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.finished
}
