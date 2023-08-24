package transaction

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
)

type VerifyTransactionSignRequest struct {
	proto *pb.VerifyTransactionSignRequest
}

func NewVerifyTransactionRequest(
	transactionHash common.Hash,
	senderPubkey p_common.PublicKey,
	senderSign p_common.Sign,
	commissionPubkey p_common.PublicKey,
	commissionSign p_common.Sign,
) types.VerifyTransactionSignRequest {
	return &VerifyTransactionSignRequest{
		proto: &pb.VerifyTransactionSignRequest{
			Hash:             transactionHash.Bytes(),
			Pubkey:           senderPubkey.Bytes(),
			Sign:             senderSign.Bytes(),
			CommissionPubkey: commissionPubkey.Bytes(),
			CommissionSign:   commissionSign.Bytes(),
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

func (request *VerifyTransactionSignRequest) TransactionHash() common.Hash {
	return common.BytesToHash(request.proto.Hash)
}

func (request *VerifyTransactionSignRequest) SenderPublicKey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(request.proto.Pubkey)
}

func (request *VerifyTransactionSignRequest) SenderSign() p_common.Sign {
	return p_common.SignFromBytes(request.proto.Sign)
}

func (request *VerifyTransactionSignRequest) CommissionPublicKey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(request.proto.CommissionPubkey)
}

func (request *VerifyTransactionSignRequest) CommissionSign() p_common.Sign {
	return p_common.SignFromBytes(request.proto.CommissionSign)
}

//

type VerifyTransactionSignResult struct {
	proto *pb.VerifyTransactionSignResult
}

func NewVerifyTransactionResult(
	transactionHash common.Hash,
	valid bool,
) types.VerifyTransactionSignResult {
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

func (result *VerifyTransactionSignResult) TransactionHash() common.Hash {
	return common.BytesToHash(result.proto.Hash)
}

func (result *VerifyTransactionSignResult) Valid() bool {
	return result.proto.Valid
}

func (result *VerifyTransactionSignResult) ResultHash() common.Hash {
	b, _ := proto.Marshal(result.proto)
	return crypto.Keccak256Hash(b)
}
