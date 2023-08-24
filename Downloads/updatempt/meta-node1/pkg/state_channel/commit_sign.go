package state_channel

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ////////////
type StateChannelCommitSign struct {
	publicKey common.PublicKey
	sign      common.Sign
}

func NewStateChannelCommitSign(
	publicKey common.PublicKey,
	sign common.Sign,
) types.StateChannelCommitSign {
	return &StateChannelCommitSign{
		publicKey: publicKey,
		sign:      sign,
	}
}

// general
func (c *StateChannelCommitSign) Marshal() ([]byte, error) {
	return proto.Marshal(c.Proto())
}

func (c *StateChannelCommitSign) Unmarshal(b []byte) error {
	pbSign := &pb.StateChannelCommitSign{}
	err := proto.Unmarshal(b, pbSign)
	if err != nil {
		return err
	}
	c.FromProto(pbSign)
	return nil
}

func (c *StateChannelCommitSign) Proto() protoreflect.ProtoMessage {
	return &pb.StateChannelCommitSign{
		PublicKey: c.publicKey.Bytes(),
		Sign:      c.sign.Bytes(),
	}
}

func (c *StateChannelCommitSign) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbSign := pbMessage.(*pb.StateChannelCommitSign)
	c.publicKey = common.PubkeyFromBytes(pbSign.PublicKey)
	c.sign = common.SignFromBytes(pbSign.Sign)
}

func (c *StateChannelCommitSign) String() string {
	return "TODO"
}

// getter

func (c *StateChannelCommitSign) Address() e_common.Address {
	return common.AddressFromPubkey(c.publicKey)
}

func (c *StateChannelCommitSign) PublicKey() common.PublicKey {
	return c.publicKey
}

func (c *StateChannelCommitSign) Sign() common.Sign {
	return c.sign
}
