package network

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type IMessageSender interface {
	SetKeyPair(keyPair *bls.KeyPair)

	SendMessage(
		connection IConnection,
		command string,
		pbMessage proto.Message,
		sign p_common.Sign,
	) error

	SendBytes(
		connection IConnection,
		command string,
		b []byte,
		sign p_common.Sign,
	) error
}

type MessageSender struct {
	keyPair *bls.KeyPair
	version string
}

func NewMessageSender(
	keyPair *bls.KeyPair,
	version string,
) IMessageSender {
	return &MessageSender{
		keyPair: keyPair,
		version: version,
	}
}

func (s *MessageSender) SetKeyPair(keyPair *bls.KeyPair) {
	s.keyPair = keyPair
}

func (s *MessageSender) SendMessage(
	connection IConnection,
	command string,
	pbMessage proto.Message,
	sign p_common.Sign,
) error {
	return SendMessage(
		connection,
		s.keyPair,
		command,
		pbMessage,
		sign,
		s.version,
	)
}

func (s *MessageSender) SendBytes(
	connection IConnection,
	command string,
	b []byte,
	sign p_common.Sign,
) error {
	return SendBytes(
		connection,
		s.keyPair,
		command,
		b,
		sign,
		s.version,
	)
}

func getHeaderForCommand(
	pubkey p_common.PublicKey,
	command string,
	toAddress common.Address,
	sign p_common.Sign,
	version string,
) *pb.Header {
	return &pb.Header{
		Command:   command,
		Pubkey:    pubkey.Bytes(),
		Sign:      sign.Bytes(),
		Version:   version,
		ToAddress: toAddress.Bytes(),
	}
}

func SendMessage(
	connection IConnection,
	keyPair *bls.KeyPair,
	command string,
	pbMessage proto.Message,
	sign p_common.Sign,
	version string,
) error {
	if connection == nil {
		return errors.New("nil connection")
	}
	body := []byte{}
	if pbMessage != nil {
		var err error
		body, err = proto.Marshal(pbMessage)
		if err != nil {
			return err
		}
	}
	if (sign == p_common.Sign{}) {
		bodyHash := crypto.Keccak256(body)
		sign = bls.Sign(keyPair.GetPrivateKey(), bodyHash)
	}

	messageProto := &pb.Message{
		Header: getHeaderForCommand(keyPair.GetPublicKey(), command, connection.GetAddress(), sign, version),
		Body:   body,
	}
	message := NewMessage(messageProto)
	return connection.SendMessage(message)
}

func SendBytes(
	connection IConnection,
	keyPair *bls.KeyPair,
	command string,
	bytes []byte,
	sign p_common.Sign,
	version string,
) error {
	if (sign == p_common.Sign{}) {
		hash := crypto.Keccak256(bytes)
		sign = bls.Sign(keyPair.GetPrivateKey(), hash)
	}
	messageProto := &pb.Message{
		Header: getHeaderForCommand(keyPair.GetPublicKey(), command, connection.GetAddress(), sign, version),
		Body:   bytes,
	}
	message := NewMessage(messageProto)
	return connection.SendMessage(message)
}
