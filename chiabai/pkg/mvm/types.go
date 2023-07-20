package mvm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract"
)

type LogsJson struct {
	Logs []map[string]interface{}
}

type MVMExecuteResult struct {
	MapAddBalance    map[string][]byte
	MapSubBalance    map[string][]byte
	MapCodeChange    map[string][]byte
	MapCodeHash      map[string][]byte
	MapStorageChange map[string][][]byte
	MapStorageRoot   map[string][]byte
	JEventLogs       LogsJson
	Status           pb.RECEIPT_STATUS
	Exception        pb.EXCEPTION
	Exmsg            string
	Return           []byte
	GasUsed          uint64
}

func (er *MVMExecuteResult) String() string {
	str := fmt.Sprintf(`
	Exit Reason: %v
	Exception: %v
	Exception Message: %v
	Output: %v
	Gas Used: %v
	Add Balance Change:
	`, er.Status, er.Exception, er.Exmsg, common.Bytes2Hex(er.Return), er.GasUsed)
	for i, v := range er.MapAddBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Sub Balance Change: ")
	for i, v := range er.MapSubBalance {
		str += fmt.Sprintf("%v: %v \n", i, uint256.NewInt(0).SetBytes(v))
	}
	str += fmt.Sprintln("Code Hash: ")
	for i, v := range er.MapCodeHash {
		str += fmt.Sprintf("%v: %v \n", i, common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Code Change: ")
	for i, v := range er.MapCodeChange {
		str += fmt.Sprintf("%v: %v \n", i, common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Storage change: ")
	for i, v := range er.MapStorageChange {
		str += fmt.Sprintf("%v:\n", i)
		for _, s := range v {
			str += fmt.Sprintf("	%v:%v\n", common.Bytes2Hex(s[:32]), common.Bytes2Hex(s[32:]))
		}
	}
	str += fmt.Sprintln("Storage root: ")
	for i, v := range er.MapStorageRoot {
		str += fmt.Sprintf("%v: %v \n", i, common.Bytes2Hex(v))
	}
	str += fmt.Sprintln("Logs: ")
	for _, v := range er.JEventLogs.Logs {
		str += fmt.Sprintf("Address %v: \n", v["address"].(string))
		str += fmt.Sprintf("Data: %v \n", v["data"].(string))
		str += fmt.Sprintln("Topic:")
		for _, t := range v["topics"].([]interface{}) {
			str += fmt.Sprintf("	%v\n", t)
		}
	}
	return str
}

func (lj *LogsJson) GetCompleteEventLogs(
	blockNumber *uint256.Int,
	transactionHash common.Hash,
) []smart_contract.IEventLog {
	rs := make([]smart_contract.IEventLog, len(lj.Logs))
	for i, v := range lj.Logs {
		rawTopics := v["topics"].([]interface{})
		topics := make([][]byte, len(rawTopics))
		for i, v := range rawTopics {
			topics[i] = common.FromHex(v.(string))
		}
		rs[i] = smart_contract.NewEventLog(
			blockNumber,
			transactionHash,
			common.HexToAddress(v["address"].(string)),
			common.FromHex(v["data"].(string)),
			topics,
		)
	}
	return rs
}
