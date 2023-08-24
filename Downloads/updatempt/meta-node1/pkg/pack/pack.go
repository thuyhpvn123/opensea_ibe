package pack

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/transaction"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Pack struct {
	proto *pb.Pack
}

func NewPack(transactions []types.Transaction, timeStamp int64) types.Pack {
	proto := &pb.Pack{
		Transactions: transaction.TransactionsToProto(transactions),
		TimeStamp:    timeStamp,
	}
	return &Pack{
		proto: proto,
	}
}

func PackFromProto(packPb *pb.Pack) types.Pack {
	return &Pack{
		proto: packPb,
	}
}

func PacksFromProto(packPb *pb.Pack) types.Pack {
	return &Pack{
		proto: packPb,
	}
}

// general
func (p *Pack) Unmarshal(b []byte) error {
	protoPack := &pb.Pack{}
	err := proto.Unmarshal(b, protoPack)
	if err != nil {
		return err
	}
	p.FromProto(protoPack)
	return nil
}

func (p *Pack) Marshal() ([]byte, error) {
	return proto.Marshal(p.Proto())
}

func (p *Pack) Proto() protoreflect.ProtoMessage {
	return p.proto
}

func (p *Pack) FromProto(pbMessage protoreflect.ProtoMessage) {
	protoPack := pbMessage.(*pb.Pack)
	p.proto = protoPack
}

func (p *Pack) String() string {
	return "TODO"
}

// getter
func (p *Pack) Transactions() []types.Transaction {
	return transaction.TransactionsFromProto(p.proto.Transactions)
}

func (p *Pack) Timestamp() uint64 {
	return uint64(p.proto.TimeStamp)
}

func (p *Pack) Hash() e_common.Hash {
	return e_common.BytesToHash(p.proto.Hash)
}

func (p *Pack) CalculateAggregateSign() common.Sign {
	transactions := p.Transactions()
	signatures := make([][]byte, len(transactions))
	for i, v := range transactions {
		sign := v.Sign()
		signatures[i] = sign.Bytes()
	}
	aggSign := bls.CreateAggregateSign(signatures)
	return common.SignFromBytes(aggSign)
}

func (p *Pack) SetAggregateSign(sign common.Sign) {
	p.proto.AggregateSign = sign.Bytes()
}

func (p *Pack) AggregateSign() common.Sign {
	return common.SignFromBytes(p.proto.AggregateSign)
}

func (p *Pack) ValidData() bool {
	return p.CalculateHash() == p.Hash()
}

func (p *Pack) ValidSign() bool {
	pubArr, hashArr, sign := p.AggregateSignData()
	return bls.VerifyAggregateSign(pubArr, sign, hashArr)
}

func (p *Pack) AggregateSignData() ([][]byte, [][]byte, []byte) {
	transactions := p.Transactions()
	totalTransaction := len(transactions)
	pubArr := make([][]byte, totalTransaction)
	hashArr := make([][]byte, totalTransaction)
	for index, t := range transactions {
		hashArr[index] = t.Hash().Bytes()
		pubArr[index] = t.Pubkey().Bytes()
	}
	return pubArr, hashArr, p.AggregateSign().Bytes()
}

// setter
func (p *Pack) SetHash(hash e_common.Hash) {
	p.proto.Hash = hash.Bytes()
}

func (p *Pack) CalculateHash() e_common.Hash {
	txtHashes := make([][]byte, len(p.proto.Transactions))
	for i, v := range p.proto.Transactions {
		txtHashes[i] = v.Hash
	}
	packHashData := &pb.PackHashData{
		TransactionHashes: txtHashes,
	}
	bHashData, _ := proto.Marshal(packHashData)
	hash := crypto.Keccak256Hash(bHashData)
	return hash
}
