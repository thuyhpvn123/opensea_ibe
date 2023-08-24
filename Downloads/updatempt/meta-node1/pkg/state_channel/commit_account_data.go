package state_channel

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CommitAccountStateChannelData struct {
	address            e_common.Address
	closeSmartContract bool
	amount             *uint256.Int
}

func NewCommitAccountStateChannelData(
	address e_common.Address,
	closeSmartContract bool,
	amount *uint256.Int,
) types.CommitAccountStateChannelData {
	return &CommitAccountStateChannelData{
		address:            address,
		closeSmartContract: closeSmartContract,
		amount:             amount,
	}
}

func (c *CommitAccountStateChannelData) Unmarshal(b []byte) error {
	cdPb := &pb.CommitAccountStateChannelData{}
	err := proto.Unmarshal(b, cdPb)
	if err != nil {
		return err
	}
	c.FromProto(cdPb)
	return nil
}

func (c *CommitAccountStateChannelData) Marshal() ([]byte, error) {
	return proto.Marshal(c.Proto())
}

func (c *CommitAccountStateChannelData) Proto() protoreflect.ProtoMessage {
	return &pb.CommitAccountStateChannelData{
		Address:            c.address.Bytes(),
		CloseSmartContract: c.closeSmartContract,
		Amount:             c.amount.Bytes(),
	}
}

func (c *CommitAccountStateChannelData) FromProto(pbMessage protoreflect.ProtoMessage) {
	cdPb := pbMessage.(*pb.CommitAccountStateChannelData)
	c.address = e_common.BytesToAddress(cdPb.Address)
	c.closeSmartContract = cdPb.CloseSmartContract
	c.amount = uint256.NewInt(0).SetBytes(cdPb.Amount)
}

// getter
func (c *CommitAccountStateChannelData) Address() e_common.Address {
	return c.address
}

func (c *CommitAccountStateChannelData) CloseSmartContract() bool {
	return c.closeSmartContract
}

func (c *CommitAccountStateChannelData) Amount() *uint256.Int {
	return c.amount
}
