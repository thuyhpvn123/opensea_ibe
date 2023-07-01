package block

import (
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IConfirmBlock interface {
	GetHash() common.Hash
	GetNumber() *uint256.Int
	GetNextLeaderAddress() common.Address
	GetValidatorSigns() map[cm.PublicKey]cm.Sign
	GetTimeStamp() uint64
	AddValidatorSign(cm.PublicKey, cm.Sign)
	GetProto() protoreflect.ProtoMessage
	GetAccountStatesRoot() common.Hash
	Marshal() ([]byte, error)
}

type ConfirmBlock struct {
	proto *pb.ConfirmBlock
}

func NewConfirmBlock(proto *pb.ConfirmBlock) IConfirmBlock {
	return &ConfirmBlock{
		proto: proto,
	}
}

func UnmarshalConfirmBlock(b []byte) (IConfirmBlock, error) {
	pbConfirmBlock := &pb.ConfirmBlock{}
	err := proto.Unmarshal(b, pbConfirmBlock)
	if err != nil {
		return nil, err
	}
	return NewConfirmBlock(pbConfirmBlock), nil
}

func ConfirmBlockFromFullBlock(b IFullBlock) *ConfirmBlock {
	mapPkValidatorSigns := b.GetValidatorSigns()
	mapStringValidatorSigns := make(map[string][]byte, len(mapPkValidatorSigns))
	for i, v := range mapPkValidatorSigns {
		mapStringValidatorSigns[hex.EncodeToString(i.Bytes())] = v.Bytes()
	}
	block := b.GetBlock()
	return &ConfirmBlock{
		proto: &pb.ConfirmBlock{
			Hash:              block.GetHash().Bytes(),
			Number:            block.GetNumber().Bytes(),
			AccountStatesRoot: block.GetAccountStatesRoot().Bytes(),
			ValidatorSigns:    mapStringValidatorSigns,
			TimeStamp:         (uint64)(time.Now().Unix()),
			NextLeaderAddress: b.GetNextLeaderAddress().Bytes(),
		},
	}
}
func (cb *ConfirmBlock) GetHash() common.Hash {
	return common.BytesToHash(cb.proto.Hash)
}

func (cb *ConfirmBlock) GetNextLeaderAddress() common.Address {
	return common.BytesToAddress(cb.proto.NextLeaderAddress)
}

func (cb *ConfirmBlock) GetNumber() *uint256.Int {
	return uint256.NewInt(0).SetBytes(cb.proto.Number)
}

func (cb *ConfirmBlock) GetValidatorSigns() map[cm.PublicKey]cm.Sign {
	rs := make(map[cm.PublicKey]cm.Sign)
	for i, v := range cb.proto.ValidatorSigns {
		rs[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}
	return rs
}

func (cb *ConfirmBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	cb.proto.ValidatorSigns[hex.EncodeToString(pk.Bytes())] = sign.Bytes()
}

func (cb *ConfirmBlock) GetProto() protoreflect.ProtoMessage {
	return cb.proto
}

func (cb *ConfirmBlock) Marshal() ([]byte, error) {
	return proto.Marshal(cb.proto)
}

func (cb *ConfirmBlock) GetAccountStatesRoot() common.Hash {
	return common.BytesToHash(cb.proto.AccountStatesRoot)
}

func (cb *ConfirmBlock) GetTimeStamp() uint64 {
	return cb.proto.TimeStamp
}

func CheckBlockValidatorSigns(block IConfirmBlock) bool {
	validatorSigns := block.GetValidatorSigns()
	for pubkey, sign := range validatorSigns {
		if !bls.VerifySign(pubkey, sign, block.GetHash().Bytes()) {
			logger.Debug(
				"CheckBlockValidatorSigns",
				hex.EncodeToString(pubkey.Bytes()),
				block.GetHash(),
				hex.EncodeToString(sign.Bytes()),
			)
			return false
		}
	}
	return true
}
