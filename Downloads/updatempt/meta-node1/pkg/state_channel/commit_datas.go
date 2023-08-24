package state_channel

import (
	"fmt"

	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/receipt"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// /////
type StateChannelCommitDatas struct {
	signs    []types.StateChannelCommitSign
	datas    []types.StateChannelCommitData
	receipts []types.Receipt
}

func NewStateChannelCommitDatas(
	signs []types.StateChannelCommitSign,
	datas []types.StateChannelCommitData,
	receipts []types.Receipt,
) types.StateChannelCommitDatas {
	return &StateChannelCommitDatas{
		signs:    signs,
		datas:    datas,
		receipts: receipts,
	}
}

// general
func (c *StateChannelCommitDatas) Marshal() ([]byte, error) {
	return proto.Marshal(c.Proto())
}

func (c *StateChannelCommitDatas) Unmarshal(b []byte) error {
	cd := &pb.StateChannelCommitDatas{}
	err := proto.Unmarshal(b, cd)
	if err != nil {
		return err
	}
	c.FromProto(cd)
	return nil
}

func (c *StateChannelCommitDatas) CommitRoot() e_common.Hash {
	//TODO
	logger.Info("TODO")
	return e_common.Hash{}
}

func (c *StateChannelCommitDatas) Proto() protoreflect.ProtoMessage {
	signs := make([]*pb.StateChannelCommitSign, len(c.signs))
	for i, v := range c.signs {
		signs[i] = v.Proto().(*pb.StateChannelCommitSign)
	}
	datas := make([]*pb.StateChannelCommitData, len(c.datas))
	for i, v := range c.datas {
		datas[i] = v.Proto().(*pb.StateChannelCommitData)
	}
	receipts := make([]*pb.Receipt, len(c.receipts))
	for i, v := range c.receipts {
		receipts[i] = v.Proto().(*pb.Receipt)
	}
	return &pb.StateChannelCommitDatas{
		Signs:    signs,
		Datas:    datas,
		Receipts: receipts,
	}
}

func (c *StateChannelCommitDatas) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbCommitDatas := pbMessage.(*pb.StateChannelCommitDatas)
	signs := make([]types.StateChannelCommitSign, len(pbCommitDatas.Signs))
	for i, v := range pbCommitDatas.Signs {
		signs[i] = &StateChannelCommitSign{}
		signs[i].FromProto(v)
	}
	c.signs = signs
	datas := make([]types.StateChannelCommitData, len(pbCommitDatas.Datas))
	for i, v := range pbCommitDatas.Datas {
		datas[i] = &StateChannelCommitData{}
		datas[i].FromProto(v)
	}
	c.datas = datas

	receipts := make([]types.Receipt, len(pbCommitDatas.Receipts))
	for i, v := range pbCommitDatas.Receipts {
		receipts[i] = &receipt.Receipt{}
		receipts[i].FromProto(v)
	}
	c.receipts = receipts
}

func (c *StateChannelCommitDatas) String() string {
	str := "Commit Datas: \n"
	for _, v := range c.datas {
		str += fmt.Sprintf("%v\n", v.String())
	}
	str += "Validator signs: \n"
	for _, v := range c.signs {
		str += fmt.Sprintf("%v\n", v.String())
	}
	return str
}

// getter
func (c *StateChannelCommitDatas) Signs() []types.StateChannelCommitSign {
	return c.signs
}

func (c *StateChannelCommitDatas) SignOfAddress(address e_common.Address) types.StateChannelCommitSign {
	logger.DebugP("StateChannelCommitDatas Signs", c.signs)
	for _, v := range c.signs {
		logger.DebugP("StateChannelCommitDatas Signs", v.Address())
		if v.Address() == address {
			return v
		}
	}
	return nil
}

func (c *StateChannelCommitDatas) Datas() []types.StateChannelCommitData {
	return c.datas
}

func (c *StateChannelCommitDatas) Receipts() []types.Receipt {
	return c.receipts
}

func (c *StateChannelCommitDatas) TotalGas() uint64 {
	rs := uint64(0)
	for _, v := range c.receipts {
		rs += v.GasUsed()
	}
	return rs
}

// setter()
func (c *StateChannelCommitDatas) AddSign(sign types.StateChannelCommitSign) {
	c.signs = append(c.signs, sign)
}
