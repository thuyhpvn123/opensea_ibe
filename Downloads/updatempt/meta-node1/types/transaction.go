package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Transaction interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	String() string

	// getter
	CalculateHash() e_common.Hash
	Hash() e_common.Hash
	NewDeviceKey() e_common.Hash
	LastDeviceKey() e_common.Hash
	FromAddress() e_common.Address
	ToAddress() e_common.Address
	Pubkey() common.PublicKey
	LastHash() e_common.Hash
	Sign() common.Sign
	Amount() *uint256.Int
	PendingUse() *uint256.Int
	Action() pb.ACTION
	BRelatedAddresses() [][]byte
	RelatedAddresses() []e_common.Address
	Data() []byte
	Fee(currentGasPrice uint64) *uint256.Int
	DeployData() DeployData
	CallData() CallData
	OpenStateChannelData() OpenStateChannelData
	CommitAccountStateChannelData() CommitAccountStateChannelData
	CommissionSign() common.Sign
	StateChannelCommitDatas() StateChannelCommitDatas
	StakeData() StakeState
	UnStakeData() StakeState
	MaxGas() uint64
	MaxGasPrice() uint64
	MaxTimeUse() uint64
	MaxFee() *uint256.Int

	// setter
	SetSign(privateKey common.PrivateKey)
	SetCommissionSign(privateKey common.PrivateKey)
	SetHash(e_common.Hash)

	// verifiers
	ValidTransactionHash() bool
	ValidSign() bool
	ValidLastHash(fromAccountState AccountState) bool
	ValidDeviceKey(fromAccountState AccountState) bool
	ValidMaxGas() bool
	ValidMaxGasPrice(currentGasPrice uint64) bool
	ValidAmount(fromAccountState AccountState, currentGasPrice uint64) bool
	ValidPendingUse(fromAccountState AccountState) bool
	ValidDeploySmartContractToAccount(fromAccountState AccountState) bool
	ValidOpenChannelToAccount(fromAccountState AccountState) bool
	ValidCallSmartContractToAccount(toAccountState AccountState) bool
}

type CallData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// geter
	Input() []byte
}

type DeployData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// getter
	Code() []byte
	StorageHost() string
	StorageAddress() e_common.Address
}

type OpenStateChannelData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	// geter
	ValidatorAddresses() []e_common.Address
}

type CommitAccountStateChannelData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	// geter
	Address() e_common.Address
	CloseSmartContract() bool
	Amount() *uint256.Int
}

type VerifyTransactionSignResult interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	TransactionHash() e_common.Hash
	Valid() bool
	ResultHash() e_common.Hash
}

type VerifyTransactionSignRequest interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	TransactionHash() e_common.Hash
	SenderPublicKey() common.PublicKey
	SenderSign() common.Sign
}
