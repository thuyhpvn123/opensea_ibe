package pack

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IVerifyPacksSignRequest interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	GetAggregateSignDatas() []IAggregateSignData
	GetHash() common.Hash
}

type IAggregateSignData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	GetPackHash() common.Hash
	GetPublickeys() [][]byte
	GetHashes() [][]byte
	GetSign() []byte
}

type IVerifyPackSignResult interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	GetPackHash() common.Hash
	GetHash() common.Hash
	GetProto() protoreflect.ProtoMessage
	Valid() bool
}

type IVerifyPacksSignResult interface {
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	GetResults() []IVerifyPackSignResult
	GetTotalPack() int
	GetRequestHash() common.Hash
	GetHash() common.Hash
	Valid() bool
}

type AggregateSignData struct {
	proto *pb.PackAggregateSignData
}

type VerifyPacksSignRequest struct {
	proto *pb.VerifyPacksSignRequest
}

type VerifyPackSignResult struct {
	proto *pb.VerifyPackSignResult
}

type VerifyPacksSignResult struct {
	proto *pb.VerifyPacksSignResult
}

func NewVerifyPacksSignRequest(packs []IAggregateSignData) IVerifyPacksSignRequest {
	requestPb := &pb.VerifyPacksSignRequest{
		PacksData: make([]*pb.PackAggregateSignData, len(packs)),
	}
	packHashes := make([][]byte, len(packs))
	for i, v := range packs {
		requestPb.PacksData[i] = &pb.PackAggregateSignData{
			PackHash:   v.GetPackHash().Bytes(),
			PublicKeys: v.GetPublickeys(),
			Hashes:     v.GetHashes(),
			Sign:       v.GetSign(),
		}
		packHashes[i] = v.GetPackHash().Bytes()
	}

	hashData := &pb.VerifyPacksSignRequestHashData{
		PackHashes: packHashes,
	}
	bHashData, _ := proto.Marshal(hashData)
	requestHash := crypto.Keccak256(bHashData)
	requestPb.Hash = requestHash
	return &VerifyPacksSignRequest{
		proto: requestPb,
	}
}

func (request *VerifyPacksSignRequest) Unmarshal(b []byte) error {
	requestPb := &pb.VerifyPacksSignRequest{}
	err := proto.Unmarshal(b, requestPb)
	if err != nil {
		return err
	}
	request.proto = requestPb
	return nil
}

func (request *VerifyPacksSignRequest) Marshal() ([]byte, error) {
	return proto.Marshal(request.proto)
}

func (request *VerifyPacksSignRequest) GetAggregateSignDatas() []IAggregateSignData {
	rs := make([]IAggregateSignData, len(request.proto.PacksData))
	for i, v := range request.proto.PacksData {
		rs[i] = AggregateSignDataFromProto(v)
	}
	return rs
}

func (request *VerifyPacksSignRequest) GetHash() common.Hash {
	return common.BytesToHash(request.proto.Hash)
}

// ==========

func NewAggregateSignData(pack IPack) IAggregateSignData {
	pubArr, hashArr, sign := pack.GetAggregateSignData()
	signDataPb := &pb.PackAggregateSignData{
		PackHash:   pack.GetHash().Bytes(),
		PublicKeys: pubArr,
		Hashes:     hashArr,
		Sign:       sign,
	}

	return &AggregateSignData{
		proto: signDataPb,
	}
}

func AggregateSignDataFromProto(proto *pb.PackAggregateSignData) IAggregateSignData {
	return &AggregateSignData{
		proto: proto,
	}
}

func (ad *AggregateSignData) Unmarshal(b []byte) error {
	adProto := &pb.PackAggregateSignData{}
	err := proto.Unmarshal(b, adProto)
	if err != nil {
		return err
	}
	ad.proto = adProto
	return nil
}

func (ad *AggregateSignData) Marshal() ([]byte, error) {
	return proto.Marshal(ad.proto)
}

func (ad *AggregateSignData) GetPackHash() common.Hash {
	return common.BytesToHash(ad.proto.PackHash)
}

func (ad *AggregateSignData) GetPublickeys() [][]byte {
	return ad.proto.PublicKeys
}

func (ad *AggregateSignData) GetHashes() [][]byte {
	return ad.proto.Hashes
}

func (ad *AggregateSignData) GetSign() []byte {
	return ad.proto.Sign
}

// ===========
func NewVerifyPackSignResult(
	packHash common.Hash,
	valid bool,
) IVerifyPackSignResult {
	rsPb := &pb.VerifyPackSignResult{
		PackHash: packHash.Bytes(),
		Valid:    valid,
	}
	return &VerifyPackSignResult{
		proto: rsPb,
	}
}

func VerifyPackSignResultFromProto(proto *pb.VerifyPackSignResult) IVerifyPackSignResult {
	return &VerifyPackSignResult{
		proto: proto,
	}
}

func (rs *VerifyPackSignResult) Unmarshal(b []byte) error {
	rsPb := &pb.VerifyPackSignResult{}
	err := proto.Unmarshal(b, rsPb)
	if err != nil {
		return err
	}
	rs.proto = rsPb
	return nil
}

func (rs *VerifyPackSignResult) Marshal() ([]byte, error) {
	return proto.Marshal(rs.proto)
}

func (rs *VerifyPackSignResult) GetPackHash() common.Hash {
	return common.BytesToHash(rs.proto.PackHash)
}

func (rs *VerifyPackSignResult) GetHash() common.Hash {
	b, _ := proto.Marshal(rs.proto)
	return crypto.Keccak256Hash(b)
}

func (rs *VerifyPackSignResult) GetProto() protoreflect.ProtoMessage {
	return rs.proto
}

func (rs *VerifyPackSignResult) Valid() bool {
	return rs.proto.Valid
}

// ===========
func NewVerifyPacksSignResult(
	requestHash common.Hash,
	results []IVerifyPackSignResult,
) IVerifyPacksSignResult {
	pbResults := make([]*pb.VerifyPackSignResult, len(results))
	for i, v := range results {
		pbResults[i] = v.GetProto().(*pb.VerifyPackSignResult)
	}
	rsPb := &pb.VerifyPacksSignResult{
		RequestHash: requestHash.Bytes(),
		Results:     pbResults,
	}
	return &VerifyPacksSignResult{
		proto: rsPb,
	}
}

func (rs *VerifyPacksSignResult) Unmarshal(b []byte) error {
	rsPb := &pb.VerifyPacksSignResult{}
	err := proto.Unmarshal(b, rsPb)
	if err != nil {
		return err
	}
	rs.proto = rsPb
	return nil
}

func (rs *VerifyPacksSignResult) Marshal() ([]byte, error) {
	return proto.Marshal(rs.proto)
}

func (rs *VerifyPacksSignResult) GetResults() []IVerifyPackSignResult {
	rss := make([]IVerifyPackSignResult, len(rs.proto.Results))
	for i, v := range rs.proto.Results {
		rss[i] = VerifyPackSignResultFromProto(v)
	}
	return rss
}

func (rs *VerifyPacksSignResult) GetTotalPack() int {
	return len(rs.proto.Results)
}

func (rs *VerifyPacksSignResult) Valid() bool {
	for _, v := range rs.GetResults() {
		if !v.Valid() {
			return false
		}
	}
	return true
}

func (rs *VerifyPacksSignResult) GetRequestHash() common.Hash {
	return common.BytesToHash(rs.proto.RequestHash)
}

func (rs *VerifyPacksSignResult) GetHash() common.Hash {
	b, _ := proto.Marshal(rs.proto)
	return crypto.Keccak256Hash(b)
}
