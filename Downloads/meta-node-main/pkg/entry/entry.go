package entry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	p "gitlab.com/meta-node/meta-node/pkg/pack"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IEntry interface {
	GetBlockNumber() *uint256.Int
	GetHash() common.Hash
	GetProto() protoreflect.ProtoMessage
	GetProtoWithBlockNumber() protoreflect.ProtoMessage
	GetPacks() []p.IPack
	Marshal() []byte
}

type Entry struct {
	blockNumber *uint256.Int
	numHashes   uint64 // number of hash since previous entry
	hash        common.Hash
	packs       []p.IPack
	proto       *pb.PohEntry
}

func NewEntry(
	blockNumber *uint256.Int,
	lastHash common.Hash,
	numHashes uint64,
	packs []p.IPack) IEntry {
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

func UnmarshalWithBlockNumber(b []byte) (IEntry, error) {
	pb := &pb.EntryWithBlockNumber{}
	err := proto.Unmarshal(b, pb)
	if err != nil {
		return nil, err
	}
	blockNumber := uint256.NewInt(0).SetBytes(pb.BlockNumber)
	packs := make([]p.IPack, len(pb.Entry.Packs))
	for i, v := range pb.Entry.Packs {
		packs[i] = p.PackFromProto(v)
	}
	return &Entry{
		blockNumber: blockNumber,
		hash:        common.BytesToHash(pb.Entry.Hash),
		numHashes:   pb.Entry.NumHashes,
		packs:       packs,
	}, nil
}

func (e *Entry) GetBlockNumber() *uint256.Int {
	return e.blockNumber
}

func (e *Entry) GetHash() common.Hash {
	return e.hash
}

func createHash(lastHash common.Hash, packs []p.IPack) common.Hash {
	packHashes := [][]byte{}

	for _, v := range packs {
		packHashes = append(packHashes, v.GetHash().Bytes())
	}

	hashData := &pb.PohHashData{
		PreHash:    lastHash.Bytes(),
		PackHashes: packHashes,
	}

	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}

func (e *Entry) GetProto() protoreflect.ProtoMessage {
	e.proto = &pb.PohEntry{
		Hash:      e.hash.Bytes(),
		NumHashes: e.numHashes,
	}
	for _, v := range e.packs {
		e.proto.Packs = append(e.proto.Packs, v.GetProto().(*pb.Pack))
	}

	return e.proto
}

func (e *Entry) GetProtoWithBlockNumber() protoreflect.ProtoMessage {
	pbWithBlockNumber := &pb.EntryWithBlockNumber{
		Entry:       e.GetProto().(*pb.PohEntry),
		BlockNumber: e.blockNumber.Bytes(),
	}
	return pbWithBlockNumber
}

func (e *Entry) Marshal() []byte {
	b, err := proto.Marshal(e.proto)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return b
}

func (e *Entry) GetPacks() []p.IPack {
	return e.packs
}
