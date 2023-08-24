package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ExecuteTransactions interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)

	// getter
	Transactions() []Transaction
	TotalTransactions() int
	GroupId() *uint256.Int
}

type ExecuteResult interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	String() string
	// getter
	Hash() common.Hash
	TransactionHash() common.Hash
	Action() pb.ACTION
	MapAddBalance() map[string][]byte
	MapSubBalance() map[string][]byte
	MapCodeHash() map[string][]byte
	MapStorageRoot() map[string][]byte
	MapLogsHash() map[string][]byte
	ReceiptStatus() pb.RECEIPT_STATUS
	Exception() pb.EXCEPTION
	Return() []byte
	GasUsed() uint64
	EventLogs() []EventLog
	// setter
}

type ExecuteResults interface {
	// general
	Unmarshal(b []byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	String() string
	// getter
	Hash() common.Hash
	GroupId() *uint256.Int
	Results() []ExecuteResult
	TotalExecute() int
}

type SmartContractData interface {
	// general
	FromProto(fbProto *pb.SmartContractData)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Proto() protoreflect.ProtoMessage
	Copy() SmartContractData
	String() string
	// getter
	Logs() []EventLogs
	Code() []byte
	Storage() map[string][]byte
	CodeHash() common.Hash
	StorageRoot() common.Hash
	LogsHash(common.Hash) common.Hash
	// setter
	SetCode([]byte)
	SetStorage(string, []byte)
	SetStorages(map[string][]byte)
	// changer
	AddLogs(EventLogs)
	ClearUpdatedLog()
}

type SmartContractDatas interface {
	// getter
	SmartContractData(common.Address) (SmartContractData, error)
	StorageIterator() storage.IIterator
	Storage() storage.Storage

	// setter
	SetSmartContractData(common.Address, SmartContractData)
	SetStorages(address common.Address, storages map[string][]byte)

	// other
	OpenStorage() error
	CloseStorage() error

	Dirty() map[common.Address]SmartContractData
	Commit() error
	Cancel()
	CopyToNewPath(newPath string) (SmartContractDatas, error)
}

type SmartContractUpdateData interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	FromProto(fbProto *pb.SmartContractUpdateData)
	String() string
	// getter
	Address() common.Address
	Code() []byte
	Storage() map[string][]byte
	Logs() []EventLogs
	CodeHash() common.Hash
	StorageRoot() common.Hash
	LogsHash(lastLogHash common.Hash) common.Hash
	BytesHash() []byte
	BlockNumber() *uint256.Int
}
