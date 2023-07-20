package pos

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/state"
)

var (
	ErrInvalidStakeStates = errors.New("invalid stake infos")
)

func GetAddressWithStakeStates(
	stakeAddress common.Address,
	_type pb.STAKE_TYPE,
	states state.IAccountStates,
	maxStaker int,
	minStakeAmount *uint256.Int,
) (map[common.Address]state.IStakeState, error) {
	stakeAccountState, err := states.GetAccountState(stakeAddress)
	if err != nil {
		logger.Error("Error when get stake info", err)
		return nil, err
	}
	stakeStates := stakeAccountState.GetStakeStatesByType(_type, minStakeAmount)
	if stakeStates == nil {
		logger.Error("Error when get stake info", ErrInvalidStakeStates)
		return nil, ErrInvalidStakeStates
	}
	sortedAddressByAmount := state.SortAdressStakeStatesByAmount(stakeStates)
	totalStaker := len(sortedAddressByAmount)
	if totalStaker > maxStaker {
		totalStaker = maxStaker
	}
	addressWithStakeStates := make(map[common.Address]state.IStakeState, totalStaker)
	for i := 0; i < totalStaker; i++ {
		addressWithStakeStates[sortedAddressByAmount[i]] = stakeStates[sortedAddressByAmount[i]]
	}
	return addressWithStakeStates, nil
}
