package block

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	cm "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ConfirmBlock struct {
	proto *pb.ConfirmBlock
}

func NewConfirmBlock(proto *pb.ConfirmBlock) types.ConfirmBlock {
	return &ConfirmBlock{
		proto: proto,
	}
}

func UnmarshalConfirmBlock(b []byte) (types.ConfirmBlock, error) {
	pbConfirmBlock := &pb.ConfirmBlock{}
	err := proto.Unmarshal(b, pbConfirmBlock)
	if err != nil {
		return nil, err
	}
	return NewConfirmBlock(pbConfirmBlock), nil
}

func ConfirmBlockFromFullBlock(b types.FullBlock) *ConfirmBlock {
	mapPkValidatorSigns := b.ValidatorSigns()
	mapStringValidatorSigns := make(map[string][]byte, len(mapPkValidatorSigns))
	for i, v := range mapPkValidatorSigns {
		mapStringValidatorSigns[hex.EncodeToString(i.Bytes())] = v.Bytes()
	}
	block := b.Block()
	return &ConfirmBlock{
		proto: &pb.ConfirmBlock{
			Hash:              block.Hash().Bytes(),
			Number:            block.Number().Bytes(),
			AccountStatesRoot: block.AccountStatesRoot().Bytes(),
			ValidatorSigns:    mapStringValidatorSigns,
			TimeStamp:         (uint64)(time.Now().Unix()),
			NextLeaderAddress: b.NextLeaderAddress().Bytes(),
		},
	}
}

// general
func (cb *ConfirmBlock) Marshal() ([]byte, error) {
	return proto.Marshal(cb.proto)
}

func (cb *ConfirmBlock) Proto() protoreflect.ProtoMessage {
	return cb.proto
}

func (cb *ConfirmBlock) FromProto(pbBlock protoreflect.ProtoMessage) {
	cb.proto = pbBlock.(*pb.ConfirmBlock)
}

func (cb *ConfirmBlock) Unmarshal(b []byte) error {
	pbConfirmBlock := &pb.ConfirmBlock{}
	err := proto.Unmarshal(b, pbConfirmBlock)
	if err != nil {
		return err
	}
	cb.FromProto(pbConfirmBlock)
	return nil
}

func (cb *ConfirmBlock) String() string {
	return "TODO"
}

// getter
func (cb *ConfirmBlock) Hash() common.Hash {
	return common.BytesToHash(cb.proto.Hash)
}

func (cb *ConfirmBlock) Number() *uint256.Int {
	return uint256.NewInt(0).SetBytes(cb.proto.Number)
}

func (cb *ConfirmBlock) NextLeaderAddress() common.Address {
	return common.BytesToAddress(cb.proto.NextLeaderAddress)
}

func (cb *ConfirmBlock) ValidatorSigns() map[cm.PublicKey]cm.Sign {
	rs := make(map[cm.PublicKey]cm.Sign)
	for i, v := range cb.proto.ValidatorSigns {
		rs[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}
	return rs
}

func (cb *ConfirmBlock) TimeStamp() uint64 {
	return cb.proto.TimeStamp
}

func (cb *ConfirmBlock) AccountStatesRoot() common.Hash {
	return common.BytesToHash(cb.proto.AccountStatesRoot)
}

// setter

func (cb *ConfirmBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	cb.proto.ValidatorSigns[hex.EncodeToString(pk.Bytes())] = sign.Bytes()
}

// other
func CheckBlockValidatorSigns(block types.ConfirmBlock) bool {
	validatorSigns := block.ValidatorSigns()
	for pubkey, sign := range validatorSigns {
		if !bls.VerifySign(pubkey, sign, block.Hash().Bytes()) {
			logger.Debug(
				"CheckBlockValidatorSigns",
				hex.EncodeToString(pubkey.Bytes()),
				block.Hash(),
				hex.EncodeToString(sign.Bytes()),
			)
			return false
		}
	}
	return true
}
