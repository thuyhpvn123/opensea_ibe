package smart_contract

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IExecuteResult interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string
	// getter
	GetProto() protoreflect.ProtoMessage
	GetHash() common.Hash
	GetTransactionHash() common.Hash
	GetAction() pb.ACTION
	GetMapAddBalance() map[string][]byte
	GetMapSubBalance() map[string][]byte
	GetMapCodeHash() map[string][]byte
	GetMapStorageRoot() map[string][]byte
	GetMapLogsHash() map[string][]byte
	GetReceiptStatus() pb.RECEIPT_STATUS
	GetException() pb.EXCEPTION
	GetReturn() []byte
	GetGasUsed() uint64
	GetEventLogs() []IEventLog

	// setter
}

type ExecuteResult struct {
	proto *pb.ExecuteResult
	//
	mapCodeChange    map[string][]byte
	mapStorageChange map[string][][]byte
	eventLogs        []IEventLog
}

type ExecuteResults struct {
	proto   *pb.ExecuteResults
	results []IExecuteResult
}

func NewExecuteResult(
	transactionHash common.Hash,
	action pb.ACTION,
	mapAddBalance map[string][]byte,
	mapSubBalance map[string][]byte,
	mapCodeChange map[string][]byte,
	mapCodeHash map[string][]byte,
	mapStorageChange map[string][][]byte,
	mapStorageRoot map[string][]byte,
	eventLogs []IEventLog,
	mapLogsHash map[string][]byte,
	status pb.RECEIPT_STATUS,
	exception pb.EXCEPTION,
	rt []byte,
	gasUsed uint64,
) IExecuteResult {
	sortedMapAddBalance := mapToSortedList(mapAddBalance)
	sortedMapSubBalance := mapToSortedList(mapSubBalance)
	sortedMapCodeHash := mapToSortedList(mapCodeHash)
	sortedMapStorageRoot := mapToSortedList(mapStorageRoot)
	sortedMapLogsHash := mapToSortedList(mapLogsHash)

	hashData := &pb.ExecuteResultHashData{
		TransactionHash:      transactionHash.Bytes(),
		Action:               action,
		SortedMapAddBalance:  sortedMapAddBalance,
		SortedMapSubBalance:  sortedMapSubBalance,
		SortedMapCodeHash:    sortedMapCodeHash,
		SortedMapStorageRoot: sortedMapStorageRoot,
		SortedMapLogsHash:    sortedMapLogsHash,
		Status:               status,
		Exception:            exception,
		Return:               rt,
		GasUsed:              gasUsed,
	}
	b, _ := proto.Marshal(hashData)
	hash := crypto.Keccak256(b)
	rs := &ExecuteResult{
		proto: &pb.ExecuteResult{
			TransactionHash: transactionHash.Bytes(),
			Action:          action,
			MapAddBalance:   mapAddBalance,
			MapSubBalance:   mapSubBalance,
			MapCodeHash:     mapCodeHash,
			MapStorageRoot:  mapStorageRoot,
			MapLogsHash:     mapLogsHash,
			Status:          status,
			Exception:       exception,
			Return:          rt,
			GasUsed:         gasUsed,
			Hash:            hash,
		},
		mapCodeChange:    mapCodeChange,
		mapStorageChange: mapStorageChange,
		eventLogs:        eventLogs,
	}
	return rs
}

func ExecuteResultFromProto(erPb *pb.ExecuteResult) IExecuteResult {
	return &ExecuteResult{
		proto: erPb,
	}
}

// general
func (r *ExecuteResult) Unmarshal(b []byte) error {
	pbRequest := &pb.ExecuteResult{}
	err := proto.Unmarshal(b, pbRequest)
	if err != nil {
		return err
	}
	r.proto = pbRequest
	return nil
}

func (r *ExecuteResult) Marshal() ([]byte, error) {
	return proto.Marshal(r.proto)
}

func (ex *ExecuteResult) String() string {
	str := fmt.Sprintf(`
	Transaction Hash: %v
	Action: %v
	Add Balance Change:
	`,
		common.Bytes2Hex(ex.proto.TransactionHash),
		ex.proto.Action,
	)
	for i, v := range ex.proto.MapAddBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Sub Balance Change: ")
	for i, v := range ex.proto.MapSubBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Code Hash: ")
	for i, v := range ex.proto.MapCodeHash {
		str += fmt.Sprintf("%v: %v \n", common.HexToAddress(i), common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Storage roots: ")
	for i, v := range ex.proto.MapStorageRoot {
		str += fmt.Sprintf("%v: %v \n", common.HexToAddress(i), common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Log hashes: ")
	for i, v := range ex.proto.MapLogsHash {
		str += fmt.Sprintf("%v: %v \n", common.HexToAddress(i), common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Logs: ")
	for _, v := range ex.eventLogs {
		str += fmt.Sprintf("	%v\n", v)
	}

	str += fmt.Sprintf(`
	Status: %v
	Exception: %v
	Return: %v
	GasUsed: %v
	Hash: %v
	`,
		ex.proto.Status,
		ex.proto.Exception,
		hex.EncodeToString(ex.proto.Return),
		ex.proto.GasUsed,
		hex.EncodeToString(ex.proto.Hash),
	)
	return str
}

// getter
func (r *ExecuteResult) GetProto() protoreflect.ProtoMessage {
	return r.proto
}

func (r *ExecuteResult) GetHash() common.Hash {
	return common.BytesToHash(r.proto.Hash)
}

func (r *ExecuteResult) GetTransactionHash() common.Hash {
	return common.BytesToHash(r.proto.TransactionHash)
}

func (r *ExecuteResult) GetAction() pb.ACTION {
	return r.proto.Action
}

func (r *ExecuteResult) GetMapAddBalance() map[string][]byte {
	return r.proto.MapAddBalance
}

func (r *ExecuteResult) GetMapSubBalance() map[string][]byte {
	return r.proto.MapSubBalance
}

func (r *ExecuteResult) GetMapCodeHash() map[string][]byte {
	return r.proto.MapCodeHash
}

func (r *ExecuteResult) GetMapStorageRoot() map[string][]byte {
	return r.proto.MapStorageRoot
}

func (r *ExecuteResult) GetMapLogsHash() map[string][]byte {
	return r.proto.MapLogsHash
}

func (r *ExecuteResult) GetReceiptStatus() pb.RECEIPT_STATUS {
	return r.proto.Status
}

func (r *ExecuteResult) GetException() pb.EXCEPTION {
	return r.proto.Exception
}

func (r *ExecuteResult) GetReturn() []byte {
	return r.proto.Return
}

func (r *ExecuteResult) GetGasUsed() uint64 {
	return r.proto.GasUsed
}

func MarshalResults(requests []ExecuteResult, groupId *uint256.Int) ([]byte, error) {
	pbResults := &pb.ExecuteResults{
		Results: make([]*pb.ExecuteResult, len(requests)),
		GroupId: groupId.Bytes(),
	}

	for i, v := range requests {
		pbResults.Results[i] = v.GetProto().(*pb.ExecuteResult)
	}
	return proto.Marshal(pbResults)
}

func (er *ExecuteResults) Unmarshal(b []byte) error {
	pbExecuteResults := &pb.ExecuteResults{}
	err := proto.Unmarshal(b, pbExecuteResults)
	if err != nil {
		return err
	}
	er.proto = pbExecuteResults
	for _, v := range pbExecuteResults.Results {
		er.results = append(er.results, ExecuteResultFromProto(v))
	}
	return nil
}

func (er *ExecuteResults) Marshal() ([]byte, error) {
	return proto.Marshal(er.proto)
}

func (er *ExecuteResults) GetProto() protoreflect.ProtoMessage {
	return er.proto
}

func (er *ExecuteResults) GetHash() common.Hash {
	return common.BytesToHash(er.proto.Hash)
}

func (er *ExecuteResults) GetGroupId() *uint256.Int {
	return uint256.NewInt(0).SetBytes(er.proto.GroupId)
}

func (er *ExecuteResults) GetResults() []IExecuteResult {
	return er.results
}

func (er *ExecuteResults) GetTotalExecute() int {
	return len(er.results)
}

func NewExecuteResults(
	results []IExecuteResult,
	groupId *uint256.Int,
) (*ExecuteResults, error) {
	pbErs := &pb.ExecuteResults{
		GroupId: groupId.Bytes(),
		Results: make([]*pb.ExecuteResult, len(results)),
	}
	hashes := make([][]byte, len(results))
	for i, v := range results {
		pbErs.Results[i] = v.GetProto().(*pb.ExecuteResult)
		hashes[i] = v.GetHash().Bytes()
	}

	hashData := &pb.ExecuteResultsHashData{
		GroupId:      groupId.Bytes(),
		ResultHashes: hashes,
	}
	bHashData, err := proto.Marshal(hashData)
	if err != nil {
		return nil, err
	}
	pbErs.Hash = crypto.Keccak256(bHashData)

	rs := &ExecuteResults{
		proto:   pbErs,
		results: results,
	}
	return rs, nil
}

func (r *ExecuteResult) GetEventLogs() []IEventLog {
	return r.eventLogs
}

func mapToSortedList(dataMap map[string][]byte) [][]byte {
	if len(dataMap) == 0 {
		return nil
	}
	rs := make([]string, len(dataMap))
	count := 0
	for i, v := range dataMap {
		rs[count] = i + hex.EncodeToString(v)
		count++
	}

	sort.Strings(rs)
	rsList := make([][]byte, len(dataMap))
	for i, v := range rs {
		rsList[i] = common.FromHex(v)
	}
	return rsList
}
