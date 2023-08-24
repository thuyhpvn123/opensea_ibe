package block

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	e "github.com/meta-node-blockchain/meta-node/pkg/entry"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Block struct {
	proto *pb.Block
}

func NewBlock(bProto *pb.Block) *Block {
	return &Block{
		proto: bProto,
	}
}

// general
func (b *Block) Marshal() ([]byte, error) {
	return proto.Marshal(b.proto)
}

func (b *Block) Unmarshal(bytes []byte) error {
	p := &pb.Block{}
	err := proto.Unmarshal(bytes, p)
	if err != nil {
		return err
	}
	if proto.Equal(p, &pb.Block{}) {
		return nil
	}
	b.proto = p
	return nil
}

func (b *Block) Proto() protoreflect.ProtoMessage {
	return b.proto
}

func (b *Block) FromProto(pbBlock protoreflect.ProtoMessage) {
	b.proto = pbBlock.(*pb.Block)
}

func (b *Block) String() string {
	str := fmt.Sprintf(`
	Hash: %v
	Number: %v
	Last entry hash: %v
	Account states root: %v
	Receipt root: %v
	Timestamp: %v
	Stake states root: %v
`,
		common.BytesToHash(b.proto.Hash),
		uint256.NewInt(0).SetBytes(b.proto.Number),
		common.BytesToHash(b.proto.LastEntryHash),
		common.BytesToHash(b.proto.AccountStatesRoot),
		common.BytesToHash(b.proto.ReceiptRoot),
		b.proto.TimeStamp,
		common.BytesToHash(b.proto.StakeStatesRoot),
	)
	return str
}

// getter
func (b *Block) Hash() common.Hash {
	return common.BytesToHash(b.proto.Hash)
}

func (b *Block) Number() *uint256.Int {
	return uint256.NewInt(0).SetBytes(b.proto.Number)
}

func (b *Block) LastEntryHash() common.Hash {
	return common.BytesToHash(b.proto.LastEntryHash)
}

func (b *Block) AccountStatesRoot() common.Hash {
	return common.BytesToHash(b.proto.AccountStatesRoot)
}

func (b *Block) IsVirtual() bool {
	return b.proto.Hash != nil && (len(b.proto.AccountStatesRoot) == 0 && len(b.proto.ReceiptRoot) == 0)
}

func (b *Block) TimeStamp() uint64 {
	return b.proto.TimeStamp
}

func (b *Block) LeaderAddress() common.Address {
	return common.BytesToAddress(b.proto.LeaderAddress)
}

func (b *Block) Prevrandao() uint64 {
	// calculate by using uint64 keccak(blockhash)
	return uint256.NewInt(0).SetBytes(
		crypto.Keccak256(b.Hash().Bytes()),
	).Uint64()
}

func (b *Block) BaseFee() uint64 {
	return b.proto.BaseFee
}

func (b *Block) GasLimit() uint64 {
	return b.proto.GasLimit
}

func (b *Block) Type() pb.BLOCK_TYPE {
	return b.proto.Type
}

func (b *Block) ReceiptRoot() common.Hash {
	return common.BytesToHash(b.proto.ReceiptRoot)
}

func (b *Block) StakeStatesRoot() common.Hash {
	return common.BytesToHash(b.proto.StakeStatesRoot)
}

// setter
func (b *Block) SetHash(hash common.Hash) {
	b.proto.Hash = hash.Bytes()
}

func (b *Block) CalculateHash() (common.Hash, error) {
	blockHashData := &pb.BlockHashData{
		Number:            b.proto.Number,
		Type:              b.proto.Type,
		LastEntryHash:     b.proto.LastEntryHash,
		LeaderAddress:     b.proto.LeaderAddress,
		AccountStatesRoot: b.proto.AccountStatesRoot,
		ReceiptRoot:       b.proto.ReceiptRoot,
		BaseFee:           b.proto.BaseFee,
		GasLimit:          b.proto.GasLimit,
		TimeStamp:         b.proto.TimeStamp,
		StakeStatesRoot:   b.proto.StakeStatesRoot,
	}

	bData, err := proto.Marshal(blockHashData)
	if err != nil {
		return common.Hash{}, err
	}
	hash := crypto.Keccak256Hash(bData)
	return hash, nil
}

func (b *Block) SetTimeStamp(timestamp uint64) {
	b.proto.TimeStamp = timestamp
}

func CheckBlockHash(block types.Block) bool {
	correctHash, err := block.CalculateHash()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error when calculate hash %v", err))
		return false
	}
	return block.Hash() == correctHash
}

func NewVirtualBlock(
	lastBlock types.Block,
	leaderAddress common.Address,
	entriesPerSlot uint64,
	hashesPerEntry uint64,
	entriesPerSecond uint64,
	accountStatesRoot common.Hash,
	stakeStatesRoot common.Hash,
) types.Block {
	number := uint256.NewInt(0).AddUint64(lastBlock.Number(), 1)

	virtualBlockProto := &pb.Block{
		Number:            number.Bytes(),
		Type:              pb.BLOCK_TYPE_FAIL,
		LeaderAddress:     leaderAddress.Bytes(),
		BaseFee:           lastBlock.BaseFee(),
		GasLimit:          lastBlock.GasLimit(),
		AccountStatesRoot: accountStatesRoot.Bytes(),
		ReceiptRoot:       nil,
		StakeStatesRoot:   stakeStatesRoot.Bytes(),
	}

	var lastEntry types.Entry
	lastEntryHash := lastBlock.LastEntryHash()
	for i := uint64(0); i < entriesPerSlot; i++ {
		lastEntry = e.NewEntry(
			number,
			lastEntryHash,
			hashesPerEntry,
			nil,
		)
		lastEntryHash = lastEntry.Hash()
	}
	virtualBlockProto.LastEntryHash = lastEntry.Hash().Bytes()

	virtualBlock := NewBlock(virtualBlockProto)
	hash, err := virtualBlock.CalculateHash()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error when hash virtual block %v", err))
	}
	virtualBlock.SetHash(hash)
	return virtualBlock
}
