package entry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	p "github.com/meta-node-blockchain/meta-node/pkg/pack"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Entry struct {
	blockNumber *uint256.Int
	numHashes   uint64 // number of hash since previous entry
	hash        common.Hash
	packs       []types.Pack
	proto       *pb.PohEntry
}

func NewEntry(
	blockNumber *uint256.Int,
	lastHash common.Hash,
	numHashes uint64,
	packs []types.Pack) types.Entry {
	var hash common.Hash
	for i := uint64(0); i < numHashes-1; i++ {
		hash = createHash(lastHash, nil)
		lastHash = hash
	}
	hash = createHash(lastHash, packs)
	return &Entry{
		blockNumber: blockNumber,
		numHashes:   numHashes,
		hash:        hash,
		packs:       packs,
	}
}

// general
func (e *Entry) Marshal() ([]byte, error) {
	return proto.Marshal(e.ProtoWithBlockNumber())
}

func (e *Entry) Unmarshal(b []byte) error {
	pbEntryWithBlockNumber := &pb.EntryWithBlockNumber{}
	err := proto.Unmarshal(b, pbEntryWithBlockNumber)
	if err != nil {
		return err
	}
	e.FromProto(pbEntryWithBlockNumber)
	return nil
}

func (e *Entry) Proto() protoreflect.ProtoMessage {
	e.proto = &pb.PohEntry{
		Hash:      e.hash.Bytes(),
		NumHashes: e.numHashes,
	}
	for _, v := range e.packs {
		e.proto.Packs = append(e.proto.Packs, v.Proto().(*pb.Pack))
	}

	return e.proto
}

func (e *Entry) ProtoWithBlockNumber() protoreflect.ProtoMessage {
	pbWithBlockNumber := &pb.EntryWithBlockNumber{
		Entry:       e.Proto().(*pb.PohEntry),
		BlockNumber: e.blockNumber.Bytes(),
	}
	return pbWithBlockNumber
}

func (e *Entry) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbEntryWithBlockNumber := pbMessage.(*pb.EntryWithBlockNumber)
	e.blockNumber = uint256.NewInt(0).SetBytes(pbEntryWithBlockNumber.BlockNumber)
	packs := make([]types.Pack, len(pbEntryWithBlockNumber.Entry.Packs))
	for i, v := range pbEntryWithBlockNumber.Entry.Packs {
		packs[i] = p.PackFromProto(v)
	}
	e.hash = common.BytesToHash(pbEntryWithBlockNumber.Entry.Hash)
	e.numHashes = pbEntryWithBlockNumber.Entry.NumHashes
	e.packs = packs
}

func (e *Entry) String() string {
	return "TODO"
}

// getter

func (e *Entry) BlockNumber() *uint256.Int {
	return e.blockNumber
}

func (e *Entry) Hash() common.Hash {
	return e.hash
}

func (e *Entry) Packs() []types.Pack {
	return e.packs
}

// setter
func createHash(lastHash common.Hash, packs []types.Pack) common.Hash {
	packHashes := [][]byte{}

	for _, v := range packs {
		packHashes = append(packHashes, v.Hash().Bytes())
	}

	hashData := &pb.PohHashData{
		PreHash:    lastHash.Bytes(),
		PackHashes: packHashes,
	}

	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}
