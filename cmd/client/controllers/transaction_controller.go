package controllers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/cmd/client/command"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/pkg/network"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	t "gitlab.com/meta-node/meta-node/pkg/transaction"
)

type ITransactionController interface {
	SendTransaction(
		lastHash common.Hash,
		toAddress common.Address,
		pendingUse *uint256.Int,
		amount *uint256.Int,
		maxGas uint64,
		maxGasFee uint64,
		maxTimeUse uint64,
		action pb.ACTION,
		data []byte,
		relatedAddress [][]byte,
	) (t.ITransaction, error)
	SetKeyPair(keyPair *bls.KeyPair)
}

type TransactionController struct {
	keyPair           *bls.KeyPair
	messageSender     network.IMessageSender
	connectionManager network.IConnectionsManager
}

func NewTransactionController(
	keyPair *bls.KeyPair,
	messageSender network.IMessageSender,
	connectionManager network.IConnectionsManager,
) ITransactionController {
	return &TransactionController{
		keyPair:           keyPair,
		messageSender:     messageSender,
		connectionManager: connectionManager,
	}
}

func (tc *TransactionController) SendTransaction(
	lastHash common.Hash,
	toAddress common.Address,
	pendingUse *uint256.Int,
	amount *uint256.Int,
	maxGas uint64,
	maxGasFee uint64,
	maxTimeUse uint64,
	action pb.ACTION,
	data []byte,
	relatedAddress [][]byte,
) (t.ITransaction, error) {

	lastDeviceKey := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	newDeviceKey := common.HexToHash("290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")
	transaction := t.NewTransaction(
		lastHash,
		tc.keyPair.GetPublicKey(),
		toAddress,
		pendingUse,
		amount,
		maxGas,
		maxGasFee,
		maxTimeUse,
		action,
		data,
		relatedAddress,
		lastDeviceKey,
		newDeviceKey,
	)
	transaction.Sign(tc.keyPair.GetPrivateKey())
	bTransaction, err := transaction.Marshal()
	if err != nil {
		return nil, err
	}
	parentConnection := tc.connectionManager.GetParentConnection()
	err = tc.messageSender.SendBytes(parentConnection, command.SendTransaction, bTransaction, transaction.GetSign())
	return transaction, err
}

func (tc *TransactionController) SetKeyPair(keyPair *bls.KeyPair) {
	tc.keyPair = keyPair
}
