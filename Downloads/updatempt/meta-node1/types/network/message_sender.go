package network

import (
	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	"github.com/meta-node-blockchain/meta-node/pkg/common"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type MessageSender interface {
	SetKeyPair(keyPair *bls.KeyPair)

	SendMessage(
		connection Connection,
		command string,
		pbMessage protoreflect.ProtoMessage,
		sign common.Sign,
	) error

	SendBytes(
		connection Connection,
		command string,
		b []byte,
		sign common.Sign,
	) error
}
