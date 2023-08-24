package smart_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type IEventLog interface {
	LoadFromProto(logPb *pb.EventLog)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	GetHash() common.Hash
	GetAddress() common.Address
	GetProto() *pb.EventLog
	GetBlockNumber() string
	GetTransactionHash() string
	GetData() string
	GetTopics() []string
	String() string
}

type EventLog struct {
	proto *pb.EventLog
}

func NewEventLog(
	blockNumber *uint256.Int,
	transactionHash common.Hash,
	address common.Address,
	data []byte,
	topics [][]byte,
) IEventLog {
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

func (l *EventLog) GetProto() *pb.EventLog {
	return l.proto
}
func (l *EventLog) LoadFromProto(logPb *pb.EventLog) {
	l.proto = logPb
}

func (l *EventLog) Unmarshal(b []byte) error {
	logPb := &pb.EventLog{}
	err := proto.Unmarshal(b, logPb)
	if err != nil {
		return err
	}
	l.LoadFromProto(logPb)
	return nil
}

func (l *EventLog) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

func (l *EventLog) GetHash() common.Hash {
	b, _ := l.Marshal()
	return crypto.Keccak256Hash(b)
}

func (l *EventLog) GetAddress() common.Address {
	return common.BytesToAddress(l.proto.Address)
}

func (l *EventLog) GetBlockNumber() string {
	return common.Bytes2Hex(l.proto.BlockNumber)
}

func (l *EventLog) GetTransactionHash() string {
	return common.Bytes2Hex(l.proto.TransactionHash)
}

func (l *EventLog) GetData() string {
	return common.Bytes2Hex(l.proto.Data)
}

func (l *EventLog) GetTopics() []string {
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

func GetNewLogHash(lastLogHash common.Hash, newLogs []IEventLog) common.Hash {
	logHashes := make([][]byte, len(newLogs))
	for i, v := range newLogs {
		logHashes[i] = v.GetHash().Bytes()
	}
	logHashData := &pb.LogsHashData{
		LastLogHash: lastLogHash.Bytes(),
		LogHashes:   logHashes,
	}
	b, _ := proto.Marshal(logHashData)
	return crypto.Keccak256Hash(b)
}

type IEventLogs interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	LoadFromProto(logPb *pb.EventLogs)
	GetProto() *pb.EventLogs
	GetEventLogList() []IEventLog
}

type EventLogs struct {
	proto *pb.EventLogs
}

func (l *EventLogs) GetEventLogList() []IEventLog {
	eventLogsPb := l.proto.GetEventLogs()
	eventLogList := make([]IEventLog, len(eventLogsPb))
	for idx, eventLog := range eventLogsPb {
		eventLogList[idx] = &EventLog{}
		eventLogList[idx].LoadFromProto(eventLog)
	}
	return eventLogList
}

func (l *EventLogs) LoadFromProto(logPb *pb.EventLogs) {
	l.proto = logPb
}

func NewEventLogs(eventLogs []IEventLog) IEventLogs {
	pbEventLogs := make([]*pb.EventLog, len(eventLogs))
	for idx, eventLog := range eventLogs {
		pbEventLogs[idx] = eventLog.GetProto()
	}
	return &EventLogs{
		proto: &pb.EventLogs{
			EventLogs: pbEventLogs,
		},
	}
}
func (l *EventLogs) Unmarshal(b []byte) error {
	logsPb := &pb.EventLogs{}
	err := proto.Unmarshal(b, logsPb)
	if err != nil {
		return err
	}
	l.LoadFromProto(logsPb)
	return nil
}

func (l *EventLogs) Marshal() ([]byte, error) {
	return proto.Marshal(l.proto)
}

func (l *EventLogs) GetProto() *pb.EventLogs {
	return l.proto
}
