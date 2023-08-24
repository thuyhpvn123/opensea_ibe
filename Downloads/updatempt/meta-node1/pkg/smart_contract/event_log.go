package smart_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
)

type EventLog struct {
	proto *pb.EventLog
}

func NewEventLog(
	blockNumber *uint256.Int,
	transactionHash common.Hash,
	address common.Address,
	data []byte,
	topics [][]byte,
) types.EventLog {
	return &EventLog{
		proto: &pb.EventLog{
			BlockNumber:     blockNumber.Bytes(),
			TransactionHash: transactionHash.Bytes(),
			Address:         address.Bytes(),
			Data:            data,
			Topics:          topics,
		},
	}
}

// general
func (l *EventLog) Proto() *pb.EventLog {
	return l.proto
}
func (l *EventLog) FromProto(logPb *pb.EventLog) {
	l.proto = logPb
}

func (l *EventLog) Unmarshal(b []byte) error {
	logPb := &pb.EventLog{}
	err := proto.Unmarshal(b, logPb)
	if err != nil {
		return err
	}
	l.FromProto(logPb)
	return nil
}

func (l *EventLog) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

// getter
func (l *EventLog) Hash() common.Hash {
	b, _ := l.Marshal()
	return crypto.Keccak256Hash(b)
}

func (l *EventLog) Address() common.Address {
	return common.BytesToAddress(l.proto.Address)
}

func (l *EventLog) BlockNumber() string {
	return common.Bytes2Hex(l.proto.BlockNumber)
}

func (l *EventLog) TransactionHash() string {
	return common.Bytes2Hex(l.proto.TransactionHash)
}

func (l *EventLog) Data() string {
	return common.Bytes2Hex(l.proto.Data)
}

func (l *EventLog) Topics() []string {
	topics := make([]string, 0)
	for _, topic := range l.proto.Topics {
		topics = append(topics, common.Bytes2Hex(topic))
	}
	return topics
}

func (l *EventLog) String() string {
	str := fmt.Sprintf(`
	Block Count: %v
	Transaction Hash: %v
	Address: %v
	Data: %v
	Topics: 
	`,
		uint256.NewInt(0).SetBytes(l.proto.BlockNumber),
		common.BytesToHash(l.proto.TransactionHash),
		common.BytesToAddress(l.proto.Address),
		common.Bytes2Hex(l.proto.Data),
	)

	for i, t := range l.proto.Topics {
		str += fmt.Sprintf("\tTopic %v: %v\n", i, common.Bytes2Hex(t))
	}
	return str
}

func NewLogHash(lastLogHash common.Hash, newLogs []types.EventLog) common.Hash {
	logHashes := make([][]byte, len(newLogs))
	for i, v := range newLogs {
		logHashes[i] = v.Hash().Bytes()
	}
	logHashData := &pb.LogsHashData{
		LastLogHash: lastLogHash.Bytes(),
		LogHashes:   logHashes,
	}
	b, _ := proto.Marshal(logHashData)
	return crypto.Keccak256Hash(b)
}

//

type EventLogs struct {
	proto *pb.EventLogs
}

func NewEventLogs(eventLogs []types.EventLog) types.EventLogs {
	pbEventLogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventLogs[idx] = eventLog.Proto()
	}
	return &EventLogs{
		proto: &pb.EventLogs{
			EventLogs: pbEventLogs,
		},
	}
}

// general
func (l *EventLogs) FromProto(logPb *pb.EventLogs) {
	l.proto = logPb
}

func (l *EventLogs) Unmarshal(b []byte) error {
	logsPb := &pb.EventLogs{}
	err := proto.Unmarshal(b, logsPb)
	if err != nil {
		return err
	}
	l.FromProto(logsPb)
	return nil
}

func (l *EventLogs) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

func (l *EventLogs) Proto() *pb.EventLogs {
	return l.proto
}

// getter
func (l *EventLogs) EventLogList() []types.EventLog {
	eventLogsPb := l.proto.EventLogs
	eventLogList := make([]types.EventLog, len(eventLogsPb))
	for idx, eventLog := range eventLogsPb {
		eventLogList[idx] = &EventLog{}
		eventLogList[idx].FromProto(eventLog)
	}
	return eventLogList
}
