package block

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/receipt"
	"github.com/meta-node-blockchain/meta-node/pkg/state"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FullBlock struct {
	block    types.Block
	receipts types.Receipts
	// account state data
	accountStateChanges map[common.Address]types.AccountState
	stakeStateChanges   map[common.Address]types.StakeStates
	//
	validatorSigns map[cm.PublicKey]cm.Sign

	nextLeaderAddress common.Address
}

func NewFullBlock(
	block types.Block,
	receipts types.Receipts,
	validatorSigns map[cm.PublicKey]cm.Sign,
	accountStateChanges map[common.Address]types.AccountState,
	nextLeaderAddress common.Address,
	stakeStateChanges map[common.Address]types.StakeStates,
) types.FullBlock {
	return &FullBlock{
		block:               block,
		receipts:            receipts,
		accountStateChanges: accountStateChanges,
		nextLeaderAddress:   nextLeaderAddress,
		stakeStateChanges:   stakeStateChanges,
		validatorSigns:      validatorSigns,
	}
}

// general
func (fb *FullBlock) Marshal() ([]byte, error) {
	return proto.Marshal(fb.Proto())
}

func (fb *FullBlock) Proto() protoreflect.ProtoMessage {
	validatorSigns := make(map[string][]byte, len(fb.validatorSigns))
	for i, v := range fb.validatorSigns {
		validatorSigns[hex.EncodeToString(i.Bytes())] = v.Bytes()
	}
	var receipts []*pb.Receipt
	if fb.receipts != nil {
		receiptsMap := fb.receipts.ReceiptsMap()
		receipts = make([]*pb.Receipt, len(receiptsMap))
		i := 0
		for _, v := range receiptsMap {
			receipts[i] = v.Proto().(*pb.Receipt)
			i++
		}
	}

	accountStateChanges := make([]*pb.AccountState, len(fb.accountStateChanges))
	i := 0
	for _, v := range fb.accountStateChanges {
		accountStateChanges[i] = v.Proto().(*pb.AccountState)
		i++
	}

	stakeStateChanges := make(map[string]*pb.StakeStates, len(fb.stakeStateChanges))
	i = 0
	for address, v := range fb.stakeStateChanges {
		stakeStateChanges[hex.EncodeToString(address.Bytes())] = v.Proto().(*pb.StakeStates)
		i++
	}
	fbProto := &pb.FullBlock{
		Block:               fb.block.Proto().(*pb.Block),
		Receipts:            receipts,
		AccountStateChanges: accountStateChanges,
		ValidatorSigns:      validatorSigns,
		StakeStateChanges:   stakeStateChanges,
	}
	return fbProto
}

func (fb *FullBlock) Unmarshal(b []byte) error {
	fbProto := &pb.FullBlock{}
	err := proto.Unmarshal(b, fbProto)
	if err != nil {
		return err
	}
	fb.FromProto(fbProto)
	return nil
}

func (fb *FullBlock) FromProto(pbMessage protoreflect.ProtoMessage) {
	fbProto := pbMessage.(*pb.FullBlock)
	block := NewBlock(fbProto.Block)

	receipts := receipt.NewReceipts()
	for _, v := range fbProto.Receipts {
		receipts.AddReceipt(receipt.ReceiptFromProto(v))
	}

	accountStateChanges := map[common.Address]types.AccountState{}
	for _, v := range fbProto.AccountStateChanges {
		accountStateChanges[common.BytesToAddress(v.Address)] = state.AccountStateFromProto(v)
	}

	stakeStateChanges := make(map[common.Address]types.StakeStates, len(fbProto.StakeStateChanges))
	for address, v := range fbProto.StakeStateChanges {
		stakeStateChanges[common.HexToAddress(address)] = &state.StakeStates{}
		stakeStateChanges[common.HexToAddress(address)].FromProto(v)
	}

	validatorSigns := make(map[cm.PublicKey]cm.Sign)
	for i, v := range fbProto.ValidatorSigns {
		validatorSigns[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}

	*fb = FullBlock{
		block:               block,
		receipts:            receipts,
		accountStateChanges: accountStateChanges,
		stakeStateChanges:   stakeStateChanges,
		validatorSigns:      validatorSigns,
	}
}

func (fb *FullBlock) String() string {
	return "TODO"
}

// getter
func (fb *FullBlock) Block() types.Block {
	return fb.block
}

func (fb *FullBlock) ValidatorSigns() map[cm.PublicKey]cm.Sign {
	return fb.validatorSigns
}

func (fb *FullBlock) AccountStateChanges() map[common.Address]types.AccountState {
	return fb.accountStateChanges
}

func (fb *FullBlock) Receipts() types.Receipts {
	return fb.receipts
}

func (fb *FullBlock) GasUsed() uint64 {
	gasUsed := uint64(0)
	if fb.receipts == nil {
		return gasUsed
	} else {
		for _, v := range fb.receipts.ReceiptsMap() {
			gasUsed += v.GasUsed()
		}
	}
	return gasUsed
}

func (fb *FullBlock) StakeStateChanges() map[common.Address]types.StakeStates {
	return fb.stakeStateChanges
}

func (fb *FullBlock) NextLeaderAddress() common.Address {
	return fb.nextLeaderAddress
}

func (fb *FullBlock) Fee() *uint256.Int {
	return uint256.NewInt(0).Mul(
		uint256.NewInt(fb.Block().BaseFee()),
		uint256.NewInt(fb.GasUsed()),
	)
}

// setter
func (fb *FullBlock) SetBlock(b types.Block) {
	fb.block = b
}

func (fb *FullBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	fb.validatorSigns[pk] = sign
}

func (fb *FullBlock) SetValidatorSigns(signs map[cm.PublicKey]cm.Sign) {
	fb.validatorSigns = signs
}

func (fb *FullBlock) SetTimeStamp(timestamp uint64) {
	fb.block.SetTimeStamp(timestamp)
}

func (fb *FullBlock) SetNextLeaderAddress(address common.Address) {
	fb.nextLeaderAddress = address
}
