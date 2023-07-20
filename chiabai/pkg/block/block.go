package block

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	e "gitlab.com/meta-node/meta-node/pkg/entry"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IBlock interface {
	// general
	Marshal() ([]byte, error)
	GetProto() protoreflect.ProtoMessage
	String() string

	// getter
	GetHash() common.Hash
	GetNumber() *uint256.Int
	GetLastEntryHash() common.Hash
	GetLeaderAddress() common.Address
	GetAccountStatesRoot() common.Hash
	GetReceiptRoot() common.Hash
	GetPrevrandao() uint64
	GetBaseFee() uint64
	GetGasLimit() uint64
	GetType() pb.BLOCK_TYPE
	GetTimeStamp() uint64

	CalculateHash() (common.Hash, error)

	// setter
	SetTimeStamp(uint64)
}

type Block struct {
	proto *pb.Block
}

func NewBlock(bProto *pb.Block) *Block {
	return &Block{
		proto: bProto,
	}
}

func Unmarshal(bytes []byte) (IBlock, error) {
	p := &pb.Block{}
	err := proto.Unmarshal(bytes, p)
	if err != nil {
		return nil, err
	}
	if proto.Equal(p, &pb.Block{}) {
		return nil, nil
	}
	return NewBlock(p), nil
}

func (b *Block) GetHash() common.Hash {
	return common.BytesToHash(b.proto.Hash)
}

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
	}

	bData, err := proto.Marshal(blockHashData)
	if err != nil {
		return common.Hash{}, err
	}
	hash := crypto.Keccak256Hash(bData)
	return hash, nil
}

func (b *Block) GetNumber() *uint256.Int {
	return uint256.NewInt(0).SetBytes(b.proto.Number)
}

func (b *Block) GetLastEntryHash() common.Hash {
	return common.BytesToHash(b.proto.LastEntryHash)
}

func (b *Block) GetAccountStatesRoot() common.Hash {
	return common.BytesToHash(b.proto.AccountStatesRoot)
}

func (b *Block) IsVirtual() bool {
	return b.proto.Hash != nil && (len(b.proto.AccountStatesRoot) == 0 && len(b.proto.ReceiptRoot) == 0)
}

func (b *Block) Marshal() ([]byte, error) {
	return proto.Marshal(b.proto)
}

func (b *Block) GetProto() protoreflect.ProtoMessage {
	return b.proto
}

func (b *Block) SetTimeStamp(timestamp uint64) {
	b.proto.TimeStamp = timestamp
}

func CheckBlockHash(block IBlock) bool {
	correctHash, err := block.CalculateHash()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error when calculate hash %v", err))
		return false
	}
	return block.GetHash() == correctHash
}

func NewVirtualBlock(
	lastBlock IBlock,
	leaderAddress common.Address,
	entriesPerSlot uint64,
	hashesPerEntry uint64,
	entriesPerSecond uint64,
	accountStatesRoot common.Hash,
) IBlock {
	number := uint256.NewInt(0).AddUint64(lastBlock.GetNumber(), 1)

	virtualBlockProto := &pb.Block{
		Number:            number.Bytes(),
		Type:              pb.BLOCK_TYPE_FAIL,
		LeaderAddress:     leaderAddress.Bytes(),
		BaseFee:           lastBlock.GetBaseFee(),
		GasLimit:          lastBlock.GetGasLimit(),
		AccountStatesRoot: accountStatesRoot.Bytes(),
		ReceiptRoot:       nil,
	}

	var lastEntry e.IEntry
	lastEntryHash := lastBlock.GetLastEntryHash()
	for i := uint64(0); i < entriesPerSlot; i++ {
		lastEntry = e.NewEntry(
			number,
			lastEntryHash,
			hashesPerEntry,
			nil,
		)
		lastEntryHash = lastEntry.GetHash()
	}
	virtualBlockProto.LastEntryHash = lastEntry.GetHash().Bytes()

	virtualBlock := NewBlock(virtualBlockProto)
	hash, err := virtualBlock.CalculateHash()
	if err != nil {
		logger.Warn(fmt.Sprintf("Error when hash virtual block %v", err))
	}
	virtualBlock.SetHash(hash)
	return virtualBlock
}

func (b *Block) String() string {
	str := fmt.Sprintf(`
	Hash: %v
	Number: %v
	Last entry hash: %v
	Account states root: %v
	Receipt root: %v
	Timestamp: %v
`,
		common.BytesToHash(b.proto.Hash),
		uint256.NewInt(0).SetBytes(b.proto.Number),
		common.BytesToHash(b.proto.LastEntryHash),
		common.BytesToHash(b.proto.AccountStatesRoot),
		common.BytesToHash(b.proto.ReceiptRoot),
		b.proto.TimeStamp,
	)
	return str
}

func (b *Block) GetTimeStamp() uint64 {
	return b.proto.TimeStamp
}

func (b *Block) GetLeaderAddress() common.Address {
	return common.BytesToAddress(b.proto.LeaderAddress)
}

func (b *Block) GetPrevrandao() uint64 {
	// calculate by using uint64 keccak(blockhash)
	return uint256.NewInt(0).SetBytes(
		crypto.Keccak256(b.GetHash().Bytes()),
	).Uint64()
}

func (b *Block) GetBaseFee() uint64 {
	return b.proto.BaseFee
}

func (b *Block) GetGasLimit() uint64 {
	return b.proto.GasLimit
}

func (b *Block) GetType() pb.BLOCK_TYPE {
	return b.proto.Type
}

func (b *Block) GetReceiptRoot() common.Hash {
	return common.BytesToHash(b.proto.ReceiptRoot)
}
