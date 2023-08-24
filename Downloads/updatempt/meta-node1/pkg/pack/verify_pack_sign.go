package pack

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

func NewVerifyPacksSignRequest(packs []types.AggregateSignData) types.VerifyPacksSignRequest {
	requestPb := &pb.VerifyPacksSignRequest{
		PacksData: make([]*pb.PackAggregateSignData, len(packs)),
	}
	packHashes := make([][]byte, len(packs))
	for i, v := range packs {
		requestPb.PacksData[i] = &pb.PackAggregateSignData{
			PackHash:   v.PackHash().Bytes(),
			PublicKeys: v.Publickeys(),
			Hashes:     v.Hashes(),
			Sign:       v.Sign(),
		}
		packHashes[i] = v.PackHash().Bytes()
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

func (request *VerifyPacksSignRequest) AggregateSignDatas() []types.AggregateSignData {
	rs := make([]types.AggregateSignData, len(request.proto.PacksData))
	for i, v := range request.proto.PacksData {
		rs[i] = AggregateSignDataFromProto(v)
	}
	return rs
}

func (request *VerifyPacksSignRequest) Hash() common.Hash {
	return common.BytesToHash(request.proto.Hash)
}

// ==========

func NewAggregateSignData(pack types.Pack) types.AggregateSignData {
	pubArr, hashArr, sign := pack.AggregateSignData()
	signDataPb := &pb.PackAggregateSignData{
		PackHash:   pack.Hash().Bytes(),
		PublicKeys: pubArr,
		Hashes:     hashArr,
		Sign:       sign,
	}

	return &AggregateSignData{
		proto: signDataPb,
	}
}

func AggregateSignDataFromProto(proto *pb.PackAggregateSignData) types.AggregateSignData {
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

func (ad *AggregateSignData) PackHash() common.Hash {
	return common.BytesToHash(ad.proto.PackHash)
}

func (ad *AggregateSignData) Publickeys() [][]byte {
	return ad.proto.PublicKeys
}

func (ad *AggregateSignData) Hashes() [][]byte {
	return ad.proto.Hashes
}

func (ad *AggregateSignData) Sign() []byte {
	return ad.proto.Sign
}

// ===========
func NewVerifyPackSignResult(
	packHash common.Hash,
	valid bool,
) types.VerifyPackSignResult {
	rsPb := &pb.VerifyPackSignResult{
		PackHash: packHash.Bytes(),
		Valid:    valid,
	}
	return &VerifyPackSignResult{
		proto: rsPb,
	}
}

func VerifyPackSignResultFromProto(proto *pb.VerifyPackSignResult) types.VerifyPackSignResult {
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

func (rs *VerifyPackSignResult) PackHash() common.Hash {
	return common.BytesToHash(rs.proto.PackHash)
}

func (rs *VerifyPackSignResult) Hash() common.Hash {
	b, _ := proto.Marshal(rs.proto)
	return crypto.Keccak256Hash(b)
}

func (rs *VerifyPackSignResult) Proto() protoreflect.ProtoMessage {
	return rs.proto
}

func (rs *VerifyPackSignResult) Valid() bool {
	return rs.proto.Valid
}

// ===========
func NewVerifyPacksSignResult(
	requestHash common.Hash,
	results []types.VerifyPackSignResult,
) types.VerifyPacksSignResult {
	pbResults := make([]*pb.VerifyPackSignResult, len(results))
	for i, v := range results {
		pbResults[i] = v.Proto().(*pb.VerifyPackSignResult)
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

func (rs *VerifyPacksSignResult) Results() []types.VerifyPackSignResult {
	rss := make([]types.VerifyPackSignResult, len(rs.proto.Results))
	for i, v := range rs.proto.Results {
		rss[i] = VerifyPackSignResultFromProto(v)
	}
	return rss
}

func (rs *VerifyPacksSignResult) TotalPack() int {
	return len(rs.proto.Results)
}

func (rs *VerifyPacksSignResult) Valid() bool {
	for _, v := range rs.Results() {
		if !v.Valid() {
			return false
		}
	}
	return true
}

func (rs *VerifyPacksSignResult) RequestHash() common.Hash {
	return common.BytesToHash(rs.proto.RequestHash)
}

func (rs *VerifyPacksSignResult) Hash() common.Hash {
	b, _ := proto.Marshal(rs.proto)
	return crypto.Keccak256Hash(b)
}
