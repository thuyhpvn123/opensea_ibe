package smart_contract

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	c_merkle_patricia_trie "github.com/meta-node-blockchain/meta-node/pkg/merkle_patricia_trie/c_version"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type SmartContractData struct {
	proto       *pb.SmartContractData
	updatedLogs []types.EventLogs
}

func NewSmartContractData(
	code []byte,
	storage map[string][]byte,
) types.SmartContractData {
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

func SmartContractDataFromProto(ssPb *pb.SmartContractData) types.SmartContractData {
	if proto.Equal(ssPb, &pb.SmartContractData{}) {
		return nil
	}
	return &SmartContractData{
		proto: ssPb,
	}
}

// general
func (s *SmartContractData) FromProto(dataPb *pb.SmartContractData) {
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
	s.FromProto(dataPb)
	return nil
}

func (s *SmartContractData) Proto() protoreflect.ProtoMessage {
	return s.proto
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
		for _, a := range v.EventLogList() {
			str += a.String()
		}
	}

	return str
}
func (s *SmartContractData) Copy() types.SmartContractData {
	cp := &SmartContractData{
		proto:       proto.Clone(s.proto).(*pb.SmartContractData),
		updatedLogs: make([]types.EventLogs, len(s.updatedLogs)),
	}
	copy(cp.updatedLogs, s.updatedLogs)
	return cp
}

// getter

func (s *SmartContractData) Logs() []types.EventLogs {
	return s.updatedLogs
}

func (s *SmartContractData) Code() []byte {
	return s.proto.Code
}

func (s *SmartContractData) Storage() map[string][]byte {
	if s.proto.Storage != nil {
		return s.proto.Storage
	}
	return make(map[string][]byte)
}

func (s *SmartContractData) CodeHash() common.Hash {
	return crypto.Keccak256Hash(s.Code())
}

// setter
func (s *SmartContractData) SetCode(code []byte) {
	s.proto.Code = code
}

func (s *SmartContractData) SetStorage(k string, v []byte) {
	s.proto.Storage[k] = v
}

func (s *SmartContractData) SetStorages(storages map[string][]byte) {
	s.proto.Storage = storages
}

func (s *SmartContractData) AddLogs(logs types.EventLogs) {
	s.updatedLogs = append(s.updatedLogs, logs)
}

// must commit befor get storage root
func (s *SmartContractData) StorageRoot() common.Hash {
	return c_merkle_patricia_trie.GetStorageRoot(s.proto.Storage)
}

// must commit befor get storage root
func (s *SmartContractData) LogsHash(lastLogHash common.Hash) common.Hash {
	return LogsHash(lastLogHash, s.Logs())
}

func (s *SmartContractData) ClearUpdatedLog() {
	s.updatedLogs = make([]types.EventLogs, 0)
}

type SmartContractUpdateData struct {
	proto *pb.SmartContractUpdateData
}

func NewSmartContractUpdateData(
	address common.Address,
	dirtyCode []byte,
	dirtyStorage map[string][]byte,
	dirtyLogs []types.EventLogs,
	blockNumber *uint256.Int,
) types.SmartContractUpdateData {

	smartContractUpdateData := &pb.SmartContractUpdateData{}
	smartContractUpdateData.Address = address.Bytes()
	smartContractUpdateData.BlockNumber = blockNumber.Bytes()

	if dirtyCode != nil {
		smartContractUpdateData.Code = dirtyCode
	}

	if dirtyStorage != nil {
		smartContractUpdateData.Storage = dirtyStorage
	}

	if len(dirtyLogs) > 0 {
		pbEventLogs := make([]*pb.EventLogs, len(dirtyLogs))
		for i, logs := range dirtyLogs {
			pbEventLogs[i] = logs.Proto()
		}
		smartContractUpdateData.EventLogs = pbEventLogs
	}

	return &SmartContractUpdateData{
		proto: smartContractUpdateData,
	}
}

// general
func (su *SmartContractUpdateData) BytesHash() []byte {
	sortedStorage := mapToSortedList(su.proto.Storage)
	hashData := &pb.SmartContractUpdateDataHash{
		Address:          su.proto.Address,
		Code:             su.proto.Code,
		SortedMapStorage: sortedStorage,
		EventLogs:        su.proto.EventLogs,
		BlockNumber:      su.proto.BlockNumber,
	}
	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256(b)
}

func (su *SmartContractUpdateData) Marshal() ([]byte, error) {
	return proto.Marshal(su.proto)
}

func (su *SmartContractUpdateData) Unmarshal(b []byte) error {
	dataPb := &pb.SmartContractUpdateData{}
	err := proto.Unmarshal(b, dataPb)
	if err != nil {
		return err
	}
	su.FromProto(dataPb)
	return nil
}

func (su *SmartContractUpdateData) FromProto(dataPb *pb.SmartContractUpdateData) {
	su.proto = dataPb
}

func (su *SmartContractUpdateData) String() string {
	str := fmt.Sprintf(`
	BlockNumber: %v
	Address: %v
	Code: %v
	Storage Map:
	`,
		uint256.NewInt(0).SetBytes(su.proto.BlockNumber),
		common.BytesToAddress(su.proto.Address),
		hex.EncodeToString(su.proto.Code),
	)
	for k, v := range su.proto.Storage {
		str += fmt.Sprintf("%v: %v \n", k, common.Bytes2Hex(v))
	}

	for _, v := range su.proto.EventLogs {
		for _, a := range v.EventLogs {
			eventLog := &EventLog{}
			eventLog.FromProto(a)
			str += fmt.Sprintf(eventLog.String())
		}
	}

	return str
}

func (s *SmartContractUpdateData) LogsHash(lastLogHash common.Hash) common.Hash {
	return LogsHash(lastLogHash, s.Logs())
}

func (s *SmartContractUpdateData) CodeHash() common.Hash {
	return crypto.Keccak256Hash(s.Code())
}

func (su *SmartContractUpdateData) Address() common.Address {
	return common.BytesToAddress(su.proto.Address)
}

func (su *SmartContractUpdateData) Code() []byte {
	return su.proto.Code
}

func (su *SmartContractUpdateData) BlockNumber() *uint256.Int {
	return uint256.NewInt(0).SetBytes(su.proto.BlockNumber)
}

func (su *SmartContractUpdateData) Logs() []types.EventLogs {
	logs := make([]types.EventLogs, len(su.proto.EventLogs))
	for i, eventLog := range su.proto.EventLogs {
		log := &EventLogs{}
		log.FromProto(eventLog)
		logs[i] = log
	}
	return logs
}
func (su *SmartContractUpdateData) Storage() map[string][]byte {
	return su.proto.Storage
}
func (s *SmartContractUpdateData) StorageRoot() common.Hash {
	return c_merkle_patricia_trie.GetStorageRoot(s.proto.Storage)
}

func LogsHash(lastLogHash common.Hash, logsArray []types.EventLogs) common.Hash {
	logger.Info("TRACE GetLogsHash", lastLogHash, logsArray)
	for _, logs := range logsArray {
		logList := logs.EventLogList()
		logsHashData := &pb.LogsHashData{
			LastLogHash: lastLogHash.Bytes(),
			LogHashes:   make([][]byte, len(logList)),
		}
		for i, log := range logList {
			logsHashData.LogHashes[i] = log.Hash().Bytes()
		}
		b, _ := proto.Marshal(logsHashData)
		lastLogHash = crypto.Keccak256Hash(b)
	}
	return lastLogHash
}
