package transaction

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type IVerifyTransactionSignResult interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	GetTransactionHash() common.Hash
	GetValid() bool
	GetResultHash() common.Hash
}

type IVerifyTransactionSignRequest interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	GetTransactionHash() common.Hash
	GetSenderPublicKey() p_common.PublicKey
	GetSenderSign() p_common.Sign
}

type VerifyTransactionSignRequest struct {
	proto *pb.VerifyTransactionSignRequest
}

func NewVerifyTransactionRequest(
	transactionHash common.Hash,
	senderPubkey p_common.PublicKey,
	senderSign p_common.Sign,
) IVerifyTransactionSignRequest {
	return &VerifyTransactionSignRequest{
		proto: &pb.VerifyTransactionSignRequest{
			Hash:   transactionHash.Bytes(),
			Pubkey: senderPubkey.Bytes(),
			Sign:   senderSign.Bytes(),
		},
	}
}

func (request *VerifyTransactionSignRequest) Unmarshal(bytes []byte) error {
	requestPb := &pb.VerifyTransactionSignRequest{}
	err := proto.Unmarshal(bytes, requestPb)
	if err != nil {
		return err
	}
	request.proto = requestPb
	return nil
}

func (request *VerifyTransactionSignRequest) Marshal() ([]byte, error) {
	return proto.Marshal(request.proto)
}

func (request *VerifyTransactionSignRequest) GetTransactionHash() common.Hash {
	return common.BytesToHash(request.proto.Hash)
}

func (request *VerifyTransactionSignRequest) GetSenderPublicKey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(request.proto.Pubkey)
}

func (request *VerifyTransactionSignRequest) GetSenderSign() p_common.Sign {
	return p_common.SignFromBytes(request.proto.Sign)
}

type VerifyTransactionSignResult struct {
	proto *pb.VerifyTransactionSignResult
}

func NewVerifyTransactionResult(
	transactionHash common.Hash,
	valid bool,
) IVerifyTransactionSignResult {
	return &VerifyTransactionSignResult{
		proto: &pb.VerifyTransactionSignResult{
			Hash:  transactionHash.Bytes(),
			Valid: valid,
		},
	}
}
func (result *VerifyTransactionSignResult) Unmarshal(bytes []byte) error {
	resultPb := &pb.VerifyTransactionSignResult{}
	err := proto.Unmarshal(bytes, resultPb)
	if err != nil {
		return err
	}
	result.proto = resultPb
	return nil
}

func (result *VerifyTransactionSignResult) Marshal() ([]byte, error) {
	return proto.Marshal(result.proto)
}

func (result *VerifyTransactionSignResult) GetTransactionHash() common.Hash {
	return common.BytesToHash(result.proto.Hash)
}

func (result *VerifyTransactionSignResult) GetValid() bool {
	return result.proto.Valid
}

func (result *VerifyTransactionSignResult) GetResultHash() common.Hash {
	b, _ := proto.Marshal(result.proto)
	return crypto.Keccak256Hash(b)
}
