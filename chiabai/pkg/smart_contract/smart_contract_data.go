package smart_contract

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	c_merkle_patricia_trie "gitlab.com/meta-node/meta-node/pkg/merkle_patricia_trie/c_version"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ISmartContractData interface {
	// general
	LoadFromProto(fbProto *pb.SmartContractData)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	Copy() ISmartContractData
	String() string
	// getter
	GetProto() protoreflect.ProtoMessage
	GetLogs() []IEventLogs
	GetCode() []byte
	GetStorage() map[string][]byte
	GetCodeHash() common.Hash
	GetStorageRoot() common.Hash
	GetLogsHash(common.Hash) common.Hash
	// setter
	SetCode([]byte)
	SetStorage(string, []byte)
	SetStorages(map[string][]byte)
	// changer
	AddLogs(IEventLogs)
	ClearUpdatedLog()
}

type SmartContractData struct {
	proto       *pb.SmartContractData
	updatedLogs []IEventLogs
}

func NewSmartContractData(
	code []byte,
	storage map[string][]byte,
) ISmartContractData {
	if storage == nil {
		return &SmartContractData{
			proto: &pb.SmartContractData{
				Code:    code,
				Storage: make(map[string][]byte),
			},
		}
	}
	return &SmartContractData{
		proto: &pb.SmartContractData{
			Code:    code,
			Storage: storage,
		},
	}
}

// general
// Copy() ISmartContractData

func (s *SmartContractData) GetLogs() []IEventLogs {
	return s.updatedLogs
}

func (s *SmartContractData) ClearUpdatedLog() {
	s.updatedLogs = make([]IEventLogs, 0)
}

func (s *SmartContractData) LoadFromProto(dataPb *pb.SmartContractData) {
	s.proto = dataPb
}

func (s *SmartContractData) Marshal() ([]byte, error) {
	return proto.Marshal(s.proto)
}

func (s *SmartContractData) Unmarshal(b []byte) error {
	dataPb := &pb.SmartContractData{}
	err := proto.Unmarshal(b, dataPb)
	if err != nil {
		return err
	}
	s.LoadFromProto(dataPb)
	return nil
}

func (s *SmartContractData) Copy() ISmartContractData {
	cp := &SmartContractData{
		proto:       proto.Clone(s.proto).(*pb.SmartContractData),
		updatedLogs: make([]IEventLogs, len(s.updatedLogs)),
	}
	copy(cp.updatedLogs, s.updatedLogs)
	return cp
}

func (s *SmartContractData) String() string {
	str := fmt.Sprintf(
		`	Code: %v
	Storage Map:
		`,
		hex.EncodeToString(s.proto.Code),
	)
	for k, v := range s.proto.Storage {
		str += fmt.Sprintf("%v: %v \n\t\t", k, common.Bytes2Hex(v))
	}

	for _, v := range s.updatedLogs {
		for _, a := range v.GetEventLogList() {
			str += a.String()
		}
	}

	return str
}

func (s *SmartContractData) GetCode() []byte {
	return s.proto.Code
}

func (s *SmartContractData) GetStorage() map[string][]byte {
	return s.proto.Storage
}

func (s *SmartContractData) SetCode(code []byte) {
	s.proto.Code = code
}

func (s *SmartContractData) SetStorage(k string, v []byte) {
	s.proto.Storage[k] = v
}

func (s *SmartContractData) SetStorages(storages map[string][]byte) {
	s.proto.Storage = storages
}

func (s *SmartContractData) AddLogs(logs IEventLogs) {
	s.updatedLogs = append(s.updatedLogs, logs)
}

func (s *SmartContractData) GetCodeHash() common.Hash {
	return crypto.Keccak256Hash(s.GetCode())
}

func (s *SmartContractData) GetProto() protoreflect.ProtoMessage {
	return s.proto
}

// must commit befor get storage root
func (s *SmartContractData) GetStorageRoot() common.Hash {
	return c_merkle_patricia_trie.GetStorageRoot(s.proto.Storage)
}

// must commit befor get storage root
func (s *SmartContractData) GetLogsHash(lastLogHash common.Hash) common.Hash {
	return GetLogsHash(lastLogHash, s.GetLogs())
}

type ISmartContractUpdateData interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	LoadFromProto(fbProto *pb.SmartContractUpdateData)
	GetAddress() common.Address
	GetCode() []byte
	GetStorage() map[string][]byte
	GetLogs() []IEventLogs
	GetCodeHash() common.Hash
	GetStorageRoot() common.Hash
	GetLogsHash(lastLogHash common.Hash) common.Hash
	String() string
	GetBytesHash() []byte
}

func NewSmartContractUpdateData(
	address common.Address,
	dirtyCode []byte,
	dirtyStorage map[string][]byte,
	dirtyLogs []IEventLogs,
) ISmartContractUpdateData {

	smartContractUpdateData := &pb.SmartContractUpdateData{}
	smartContractUpdateData.Address = address.Bytes()

	if dirtyCode != nil {
		smartContractUpdateData.Code = dirtyCode
	}

	if dirtyStorage != nil {
		smartContractUpdateData.Storage = dirtyStorage
	}

	if len(dirtyLogs) > 0 {
		pbEventLogs := make([]*pb.EventLogs, len(dirtyLogs))
		for i, logs := range dirtyLogs {
			pbEventLogs[i] = logs.GetProto()
		}
		smartContractUpdateData.EventLogs = pbEventLogs
	}

	return &SmartContractUpdateData{
		proto: smartContractUpdateData,
	}
}

type SmartContractUpdateData struct {
	proto *pb.SmartContractUpdateData
}

func (su *SmartContractUpdateData) GetBytesHash() []byte {
	sortedStorage := mapToSortedList(su.proto.Storage)
	hashData := &pb.SmartContractUpdateDataHash{
		Address:          su.proto.Address,
		Code:             su.proto.Code,
		SortedMapStorage: sortedStorage,
		EventLogs:        su.proto.EventLogs,
	}
	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256(b)
}

func (su *SmartContractUpdateData) Marshal() ([]byte, error) {
	return proto.Marshal(su.proto)
}

func (s *SmartContractUpdateData) GetLogsHash(lastLogHash common.Hash) common.Hash {
	return GetLogsHash(lastLogHash, s.GetLogs())
}

func (s *SmartContractUpdateData) GetCodeHash() common.Hash {
	return crypto.Keccak256Hash(s.GetCode())
}

func (su *SmartContractUpdateData) GetAddress() common.Address {
	return common.BytesToAddress(su.proto.Address)
}

func (su *SmartContractUpdateData) GetCode() []byte {
	return su.proto.Code
}

func (su *SmartContractUpdateData) GetLogs() []IEventLogs {
	logs := make([]IEventLogs, len(su.proto.EventLogs))
	for i, eventLog := range su.proto.EventLogs {
		log := &EventLogs{}
		log.LoadFromProto(eventLog)
		logs[i] = log
	}
	return logs
}
func (su *SmartContractUpdateData) GetStorage() map[string][]byte {
	return su.proto.Storage
}
func (s *SmartContractUpdateData) GetStorageRoot() common.Hash {
	return c_merkle_patricia_trie.GetStorageRoot(s.proto.Storage)
}
func (su *SmartContractUpdateData) Unmarshal(b []byte) error {
	dataPb := &pb.SmartContractUpdateData{}
	err := proto.Unmarshal(b, dataPb)
	if err != nil {
		return err
	}
	su.LoadFromProto(dataPb)
	return nil
}

func (su *SmartContractUpdateData) LoadFromProto(dataPb *pb.SmartContractUpdateData) {
	su.proto = dataPb
}

func SmartContractDataFromProto(ssPb *pb.SmartContractData) ISmartContractData {
	if proto.Equal(ssPb, &pb.SmartContractData{}) {
		return nil
	}
	return &SmartContractData{
		proto: ssPb,
	}
}

func GetLogsHash(lastLogHash common.Hash, logsArray []IEventLogs) common.Hash {
	for _, logs := range logsArray {
		logList := logs.GetEventLogList()
		logsHashData := &pb.LogsHashData{
			LastLogHash: lastLogHash.Bytes(),
			LogHashes:   make([][]byte, len(logList)),
		}
		for i, log := range logList {
			logsHashData.LogHashes[i] = log.GetHash().Bytes()
		}
		b, _ := proto.Marshal(logsHashData)
		lastLogHash = crypto.Keccak256Hash(b)
	}
	return lastLogHash
}

func (su *SmartContractUpdateData) String() string {
	str := fmt.Sprintf(`
	Address: %v
	Code: %v
	Storage Map:
	`,
		common.BytesToAddress(su.proto.Address),
		hex.EncodeToString(su.proto.Code),
	)
	for k, v := range su.proto.Storage {
		str += fmt.Sprintf("%v: %v \n", k, common.Bytes2Hex(v))
	}

	for _, v := range su.proto.EventLogs {
		for _, a := range v.EventLogs {
			eventLog := &EventLog{}
			eventLog.LoadFromProto(a)
			str += fmt.Sprintf(eventLog.String())
		}
	}

	return str
}
