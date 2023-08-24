package poh

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const MAX_PERCENTAGE = 1000
const SLOT_PER_LEADER = 4

var (
	ErrTotalSlotNotDivisible = errors.New("total slot not divisible")
	ErrInvalidStakeStates    = errors.New("invalid stake infos")
)

type LeaderSchedule struct {
	seed *uint256.Int

	fromSlot *uint256.Int
	toSlot   *uint256.Int

	slots map[uint256.Int]common.Address
}

func NewLeaderSchedule(
	seed *uint256.Int,
	stakeStatesManager types.StakeStatesManager,
	validatorStakeRequire *uint256.Int,
	fromSlot *uint256.Int,
	toSlot *uint256.Int,
) (*LeaderSchedule, error) {

	validatorWithStakeState, err := stakeStatesManager.StakeStates(
		p_common.VALIDATOR_STAKE_POOL_ADDRESS,
	).MapStakeState(
		pb.STAKE_TYPE_VALIDATOR,
		p_common.MAX_VALIDATOR,
		p_common.MINIMUM_VALIDATOR_STAKE_AMOUNT,
	)
	logger.DebugP("ZZZZ", validatorWithStakeState)
	if err != nil {
		return nil, err
	}
	validatorsWithStake := make(map[common.Address]*uint256.Int, len(validatorWithStakeState))
	for address, stakeState := range validatorWithStakeState {
		validatorsWithStake[address] = stakeState.Amount()
	}

	stakePercentages := calculateValidatorStakePercentage(validatorsWithStake)
	totalSlot := uint256.NewInt(0).Sub(toSlot, fromSlot)

	if (totalSlot.Uint64()+1)%SLOT_PER_LEADER != 0 {
		logger.Error(fromSlot)
		return nil, ErrTotalSlotNotDivisible
	}
	logger.DebugP("ZZZZ 2")

	slots := calculateSlot(stakePercentages, seed, fromSlot, totalSlot.Uint64())
	// calculate slots
	ls := &LeaderSchedule{
		seed:     seed,
		fromSlot: fromSlot,
		toSlot:   toSlot,
		slots:    slots,
	}
	return ls, nil
}

// general
func (ls *LeaderSchedule) Marshal() ([]byte, error) {
	return proto.Marshal(ls.Proto())
}

func (ls *LeaderSchedule) Unmarshal(b []byte) error {
	protoLs := &pb.LeaderSchedule{}
	err := proto.Unmarshal(b, protoLs)
	if err != nil {
		return err
	}
	ls.FromProto(protoLs)
	return nil
}

func (ls *LeaderSchedule) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbSchedule := pbMessage.(*pb.LeaderSchedule)
	ls.seed = uint256.NewInt(0).SetBytes(pbSchedule.Seed)
	ls.fromSlot = uint256.NewInt(0).SetBytes(pbSchedule.FromSlot)
	ls.toSlot = uint256.NewInt(0).SetBytes(pbSchedule.ToSlot)
	ls.slots = make(map[uint256.Int]common.Address)
	for i, v := range pbSchedule.Slots {
		ls.slots[*uint256.NewInt(0).SetBytes(common.FromHex(i))] = common.BytesToAddress(v)
	}
}

func (ls *LeaderSchedule) Proto() protoreflect.ProtoMessage {
	protoLs := &pb.LeaderSchedule{
		Seed:     ls.seed.Bytes(),
		FromSlot: ls.fromSlot.Bytes(),
		ToSlot:   ls.toSlot.Bytes(),
	}
	slots := make(map[string][]byte, len(ls.slots))
	for i, v := range ls.slots {
		slots[hex.EncodeToString(i.Bytes())] = v.Bytes()
	}
	protoLs.Slots = slots
	return protoLs
}

func (ls *LeaderSchedule) String() string {
	str := "Slot: \n"
	for i, v := range ls.slots {
		str += fmt.Sprintf("%v %v\n", &i, v)
	}
	return str
}

// getter
func (ls *LeaderSchedule) LeaderAtSlot(slot *uint256.Int) common.Address {
	return ls.slots[*slot]
}

func (ls *LeaderSchedule) ToSlot() *uint256.Int {
	return ls.toSlot
}

func (ls *LeaderSchedule) SetSlots(slots map[uint256.Int]common.Address) {
	ls.slots = slots
}

func (ls *LeaderSchedule) SetToSlot(toSlot *uint256.Int) {
	ls.toSlot = toSlot
}

//

type Range struct {
	From uint64
	To   uint64
}

func calculateSlot(
	stakePercentage map[common.Address]uint64,
	seed *uint256.Int,
	fromSlot *uint256.Int,
	totalSlot uint64,
) map[uint256.Int]common.Address {
	rs := map[uint256.Int]common.Address{}
	// range from 0 to 1000
	rangeAddress := map[*Range]common.Address{}
	track := uint64(0)
	var lastRange *Range

	addresses := make([]common.Address, 0, len(stakePercentage))
	for a := range stakePercentage {
		addresses = append(addresses, a)
	}
	sort.Slice(addresses, func(i, j int) bool {
		return hex.EncodeToString(addresses[i].Bytes()) < hex.EncodeToString(addresses[j].Bytes())
	})

	for _, address := range addresses {
		rangeV := &Range{track, track + stakePercentage[address]}
		rangeAddress[rangeV] = address
		track += stakePercentage[address]
		lastRange = rangeV
	}
	if lastRange != nil && lastRange.To != MAX_PERCENTAGE {
		lastRange.To = MAX_PERCENTAGE
	}

	// for each slot use random from seed to get leader address, then update seed
	slotCount := fromSlot.Clone()
	toSlot := uint256.NewInt(0).AddUint64(slotCount, totalSlot)
	logger.Info("to slot, slot count", toSlot, slotCount)
	logger.Info("stakePercentage", stakePercentage)
	for toSlot.Gt(slotCount) {
		r := rand.New(rand.NewSource(int64(seed.Uint64())))
		randValue := uint64(r.Intn(1000))
		for rangeV, address := range rangeAddress {
			if randValue >= rangeV.From && randValue < rangeV.To {
				for u := uint64(0); u < SLOT_PER_LEADER; u++ {
					rs[*slotCount] = address
					slotCount = slotCount.AddUint64(slotCount, 1)
				}
				break
			}
		}
		seed = seed.AddUint64(seed, 1)
	}

	return rs
}

func calculateValidatorStakePercentage(validatorsWithStake map[common.Address]*uint256.Int) map[common.Address]uint64 {
	// calculate total staked of validators
	totalStake := uint256.NewInt(0)
	for _, v := range validatorsWithStake {
		totalStake = totalStake.Add(totalStake, v)
	}
	stakePercentage := map[common.Address]uint64{}
	// zero staked flow
	if totalStake.IsZero() {
		for k := range validatorsWithStake {
			percentage := float64(MAX_PERCENTAGE) / float64(len(validatorsWithStake))
			stakePercentage[k] = uint64(percentage)
		}
	} else {
		// calculate /1000 staked of each validator
		for k, v := range validatorsWithStake {
			percentage := uint256.NewInt(0).Mul(v, uint256.NewInt(MAX_PERCENTAGE))
			percentage = percentage.Div(percentage, totalStake)
			stakePercentage[k] = percentage.Uint64()
		}
	}
	return stakePercentage
}
