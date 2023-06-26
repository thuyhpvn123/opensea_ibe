package pack

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IPack interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	GetTimestamp() uint64
	CalculateHash() common.Hash
	GetHash() common.Hash
	SetHash(common.Hash)
	GetProto() protoreflect.ProtoMessage
	GetBytes() []byte
	GetTransactions() []transaction.ITransaction
	CalculateAggregateSign() p_common.Sign
	SetAggregateSign(p_common.Sign)
	GetAggregateSign() p_common.Sign
	GetAggregateSignData() (pubArr [][]byte, hashArr [][]byte, sign []byte)
	ValidData() bool
	ValidSign() bool
}

type Pack struct {
	proto *pb.Pack
}

func NewPack(transactions []transaction.ITransaction, timeStamp int64) IPack {
	proto := &pb.Pack{
		Transactions: transaction.TransactionsToProto(transactions),
		TimeStamp:    timeStamp,
	}
	return &Pack{
		proto: proto,
	}
}

func PackFromProto(packPb *pb.Pack) IPack {
	return &Pack{
		proto: packPb,
	}
}

func PacksFromProto(packPb *pb.Pack) IPack {
	return &Pack{
		proto: packPb,
	}
}

func (p *Pack) Unmarshal(b []byte) error {
	protoPack := &pb.Pack{}
	err := proto.Unmarshal(b, protoPack)
	if err != nil {
		return err
	}
	p.proto = protoPack
	return nil
}

func (p *Pack) Marshal() ([]byte, error) {
	return proto.Marshal(p.proto)
}

func (p *Pack) CalculateHash() common.Hash {
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

func (p *Pack) GetTimestamp() uint64 {
	return uint64(p.proto.TimeStamp)
}

func (p *Pack) GetHash() common.Hash {
	return common.BytesToHash(p.proto.Hash)
}

func (p *Pack) SetHash(hash common.Hash) {
	p.proto.Hash = hash.Bytes()
}

func (p *Pack) GetProto() protoreflect.ProtoMessage {
	return p.proto
}

func (p *Pack) GetBytes() []byte {
	b, err := proto.Marshal(p.proto)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return b
}

func (p *Pack) GetTransactions() []transaction.ITransaction {
	return transaction.TransactionsFromProto(p.proto.Transactions)
}

func (p *Pack) CalculateAggregateSign() p_common.Sign {
	transactions := p.GetTransactions()
	signatures := make([][]byte, len(transactions))
	for i, v := range transactions {
		sign := v.GetSign()
		signatures[i] = sign.Bytes()
	}
	aggSign := bls.CreateAggregateSign(signatures)
	return p_common.SignFromBytes(aggSign)
}

func (p *Pack) SetAggregateSign(sign p_common.Sign) {
	p.proto.AggregateSign = sign.Bytes()
}

func (p *Pack) GetAggregateSign() p_common.Sign {
	return p_common.SignFromBytes(p.proto.AggregateSign)
}

func (p *Pack) ValidData() bool {
	return p.CalculateHash() == p.GetHash()
}

func (p *Pack) ValidSign() bool {
	pubArr, hashArr, sign := p.GetAggregateSignData()
	return bls.VerifyAggregateSign(pubArr, sign, hashArr)
}

func (p *Pack) GetAggregateSignData() ([][]byte, [][]byte, []byte) {
	transactions := p.GetTransactions()
	totalTransaction := len(transactions)
	pubArr := make([][]byte, totalTransaction)
	hashArr := make([][]byte, totalTransaction)
	for index, t := range transactions {
		hashArr[index] = t.GetHash().Bytes()
		pubArr[index] = t.GetPubkey().Bytes()
	}
	return pubArr, hashArr, p.GetAggregateSign().Bytes()
}
