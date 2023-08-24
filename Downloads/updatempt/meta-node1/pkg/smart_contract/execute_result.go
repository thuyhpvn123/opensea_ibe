package smart_contract

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ExecuteResult struct {
	proto *pb.ExecuteResult
	//
	mapCodeChange    map[string][]byte
	mapStorageChange map[string][][]byte
	eventLogs        []types.EventLog
}

type ExecuteResults struct {
	proto   *pb.ExecuteResults
	results []types.ExecuteResult
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
	eventLogs []types.EventLog,
	mapLogsHash map[string][]byte,
	status pb.RECEIPT_STATUS,
	exception pb.EXCEPTION,
	rt []byte,
	gasUsed uint64,
) types.ExecuteResult {
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

func ExecuteResultFromProto(erPb *pb.ExecuteResult) types.ExecuteResult {
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
func (r *ExecuteResult) Proto() protoreflect.ProtoMessage {
	return r.proto
}

func (r *ExecuteResult) Hash() common.Hash {
	return common.BytesToHash(r.proto.Hash)
}

func (r *ExecuteResult) TransactionHash() common.Hash {
	return common.BytesToHash(r.proto.TransactionHash)
}

func (r *ExecuteResult) Action() pb.ACTION {
	return r.proto.Action
}

func (r *ExecuteResult) MapAddBalance() map[string][]byte {
	return r.proto.MapAddBalance
}

func (r *ExecuteResult) MapSubBalance() map[string][]byte {
	return r.proto.MapSubBalance
}

func (r *ExecuteResult) MapCodeHash() map[string][]byte {
	return r.proto.MapCodeHash
}

func (r *ExecuteResult) MapStorageRoot() map[string][]byte {
	return r.proto.MapStorageRoot
}

func (r *ExecuteResult) MapLogsHash() map[string][]byte {
	return r.proto.MapLogsHash
}

func (r *ExecuteResult) ReceiptStatus() pb.RECEIPT_STATUS {
	return r.proto.Status
}

func (r *ExecuteResult) Exception() pb.EXCEPTION {
	return r.proto.Exception
}

func (r *ExecuteResult) Return() []byte {
	return r.proto.Return
}

func (r *ExecuteResult) GasUsed() uint64 {
	return r.proto.GasUsed
}

func (r *ExecuteResult) EventLogs() []types.EventLog {
	return r.eventLogs
}

// ExecuteResults

func NewExecuteResults(
	results []types.ExecuteResult,
	groupId *uint256.Int,
) (*ExecuteResults, error) {
	pbErs := &pb.ExecuteResults{
		GroupId: groupId.Bytes(),
		Results: make([]*pb.ExecuteResult, len(results)),
	}
	hashes := make([][]byte, len(results))
	for i, v := range results {
		pbErs.Results[i] = v.Proto().(*pb.ExecuteResult)
		hashes[i] = v.Hash().Bytes()
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

func (er *ExecuteResults) Proto() protoreflect.ProtoMessage {
	return er.proto
}

func (er *ExecuteResults) String() string {
	return "TODO"
}

func (er *ExecuteResults) Hash() common.Hash {
	return common.BytesToHash(er.proto.Hash)
}

func (er *ExecuteResults) GroupId() *uint256.Int {
	return uint256.NewInt(0).SetBytes(er.proto.GroupId)
}

func (er *ExecuteResults) Results() []types.ExecuteResult {
	return er.results
}

func (er *ExecuteResults) TotalExecute() int {
	return len(er.results)
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
