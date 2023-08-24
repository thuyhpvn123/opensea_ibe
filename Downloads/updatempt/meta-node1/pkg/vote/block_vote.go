package vote

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type BlockVote struct {
	proto *pb.BlockVote
}

func NewBlockVote(
	proto *pb.BlockVote,
) types.BlockVote {
	return &BlockVote{
		proto,
	}
}

// general
func (v *BlockVote) Marshal() ([]byte, error) {
	return proto.Marshal(v.Proto())
}

func (v *BlockVote) Unmarshal(
	b []byte,
	pubkey cm.PublicKey,
	sign cm.Sign,
) error {
	pbBlockVote := &pb.BlockVote{}
	err := proto.Unmarshal(b, pbBlockVote)
	if err != nil {
		return err
	}
	pbBlockVote.Pubkey = pubkey.Bytes()
	pbBlockVote.Sign = sign.Bytes()
	v.FromProto(pbBlockVote)
	return nil
}

func (v *BlockVote) FromProto(pbMessage protoreflect.ProtoMessage) {
	v.proto = pbMessage.(*pb.BlockVote)
}

func (v *BlockVote) Proto() protoreflect.ProtoMessage {
	return v.proto
}

// getter
func (v *BlockVote) BlockNumber() *uint256.Int {
	return uint256.NewInt(0).SetBytes(v.proto.Number)
}

func (v *BlockVote) Hash() common.Hash {
	return common.BytesToHash(v.proto.Hash)
}

func (v *BlockVote) Value() interface{} {
	if len(v.proto.BlockData) == 0 {
		return nil
	}
	return v.proto.BlockData
}

func (v *BlockVote) PublicKey() cm.PublicKey {
	return cm.PubkeyFromBytes(v.proto.Pubkey)
}

func (v *BlockVote) Address() common.Address {
	return cm.AddressFromPubkey(v.PublicKey())
}

func (v *BlockVote) Sign() cm.Sign {
	return cm.SignFromBytes(v.proto.Sign)
}
