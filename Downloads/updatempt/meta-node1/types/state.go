package types

import (
	e_common "github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type AccountState interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	Copy() AccountState
	String() string

	// getter
	Address() e_common.Address
	LastHash() e_common.Hash
	Balance() *uint256.Int
	PendingBalance() *uint256.Int
	TotalBalance() *uint256.Int
	SmartContractState() SmartContractState
	DeviceKey() e_common.Hash
	StateChannelState() StateChannelState

	// setter
	SetSmartContractState(smState SmartContractState)
	SetStateChannelState(scs StateChannelState)
	SetBalance(newBalance *uint256.Int)
	SetPendingBalance(newBalance *uint256.Int)
	SetNewDeviceKey(newDeviceKey e_common.Hash)
	SetLastHash(newLastHash e_common.Hash)
	AddPendingBalance(amount *uint256.Int)
	SubPendingBalance(amount *uint256.Int) error
	SubBalance(amount *uint256.Int) error
	SubTotalBalance(amount *uint256.Int) error
	AddBalance(amount *uint256.Int)
	SetCodeHash(hash e_common.Hash)
	SetStorageHost(storageHost string)
	SetStorageRoot(hash e_common.Hash)
	SetLogsHash(hash e_common.Hash)
}

type StakeState interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	Copy() StakeState
	String() string

	// getter
	Address() e_common.Address
	Amount() *uint256.Int
	Type() pb.STAKE_TYPE
	PublicConnectionAddress() string

	// setter
	SetAddress(e_common.Address)
	AddAmount(*uint256.Int)
	SubAmount(*uint256.Int) error
	SetAmount(*uint256.Int)
	SetType(pb.STAKE_TYPE)
	SetPublicConnectionAddress(string)
}

type StakeStates interface {
	// general
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)
	Copy() StakeStates
	// getter
	MapStakeState(
		_type pb.STAKE_TYPE,
		maxStaker int,
		minStakeAmount *uint256.Int,
	) (map[e_common.Address]StakeState, error)
	StakeState(address e_common.Address, _type pb.STAKE_TYPE) StakeState
	// setter
	AddStakingBalance(StakeState)
	SubStakingBalance(StakeState) error
	ChangePublicConnectionAddress(StakeState)
}

type SmartContractState interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string

	// getter
	Proto() protoreflect.ProtoMessage
	CreatorPublicKey() common.PublicKey
	StorageHost() string
	StorageAddress() e_common.Address
	CodeHash() e_common.Hash
	StorageRoot() e_common.Hash
	LogsHash() e_common.Hash
	RelatedAddress() []e_common.Address
	LockingStateChannel() e_common.Address

	// setter
	SetCreatorPublicKey(common.PublicKey)
	SetStorageHost(string)
	SetCodeHash(e_common.Hash)
	SetStorageRoot(e_common.Hash)
	SetLogsHash(e_common.Hash)
	SetRelatedAddress([]e_common.Address)
	SetLockingStateChannel(e_common.Address)
}

type SmartContractStateConfirm interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string
	Address() e_common.Address
	SmartContractState() SmartContractState
	BlockNumber() *uint256.Int
}

type StateChannelState interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string
	Proto() protoreflect.ProtoMessage
	FromProto(protoreflect.ProtoMessage)

	ValidatorAddresses() []e_common.Address
}
