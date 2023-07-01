package network

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(protoStruct protoreflect.ProtoMessage) error
	GetBody() []byte
	GetToAddress() common.Address
	GetPubkey() cm.PublicKey
	GetSign() cm.Sign
	GetCommand() string
}

type Message struct {
	proto *pb.Message
}

func NewMessage(pbMessage *pb.Message) IMessage {
	return &Message{
		proto: pbMessage,
	}
}

func (m *Message) Marshal() ([]byte, error) {
	return proto.Marshal(m.proto)
}

func (m *Message) Unmarshal(protoStruct protoreflect.ProtoMessage) error {
	err := proto.Unmarshal(m.proto.Body, protoStruct)
	return err
}

func (m *Message) GetCommand() string {
	return m.proto.Header.Command
}

func (m *Message) GetBody() []byte {
	return m.proto.Body
}

func (m *Message) GetPubkey() cm.PublicKey {
	return cm.PubkeyFromBytes(m.proto.Header.Pubkey)
}

func (m *Message) GetSign() cm.Sign {
	return cm.SignFromBytes(m.proto.Header.Sign)
}

func (m *Message) GetToAddress() common.Address {
	return common.BytesToAddress(m.proto.Header.ToAddress)
}

func (m *Message) String() string {
	str := fmt.Sprintf(`
	Header:
		Command: %v
		Pubkey: %v
		ToAddress: %v
		Sign: %v
		Version: %v
	Body: %v
`,
		m.proto.Header.Command,
		hex.EncodeToString(m.proto.Header.Pubkey),
		hex.EncodeToString(m.proto.Header.ToAddress),
		hex.EncodeToString(m.proto.Header.Sign),
		m.proto.Header.Version,
		hex.EncodeToString(m.proto.Body),
	)
	return str
}
