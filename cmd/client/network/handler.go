package network

import (
	"errors"
	"fmt"

	"gitlab.com/meta-node/meta-node/pkg/smart_contract"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/cmd/client/command"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/receipt"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

var (
	ErrorCommandNotFound = errors.New("command not found")
)

type Handler struct {
	accountStateChan chan state.IAccountState
	receiptChan      chan receipt.IReceipt
}

func NewHandler(
	accountStateChan chan state.IAccountState,
	receiptChan chan receipt.IReceipt,
) *Handler {
	return &Handler{
		accountStateChan: accountStateChan,
		receiptChan:      receiptChan,
	}
}

func (h *Handler) HandleRequest(request network.IRequest) (err error) {
	cmd := request.GetMessage().GetCommand()
	logger.Debug("handling command: " + cmd)
	switch cmd {
	case command.InitConnection:
		return h.handleInitConnection(request)
	case command.AccountState:
		return h.handleAccountState(request)
	case command.Receipt:
		return h.handleReceipt(request)
	case command.TransactionError:
		return h.handleTransactionError(request)
	case command.EventLogs:
		return h.handleEventLogs(request)
	}
	return ErrorCommandNotFound
}

/*
handleInitConnection will receive request from connection
then init that connection with data in request then
add it to connection manager
*/
func (h *Handler) handleInitConnection(request network.IRequest) (err error) {
	conn := request.GetConnection()
	initData := &pb.InitConnection{}
	err = request.GetMessage().Unmarshal(initData)
	if err != nil {
		return err
	}
	address := common.BytesToAddress(initData.Address)
	logger.Debug(fmt.Sprintf(
		"init connection from %v type %v", address, initData.Type,
	))
	conn.Init(address, initData.Type, initData.PublicConnectionAddress)
	return nil
}

/*
handleAccountState will receive account state from connection
then push it to account state chan
*/
func (h *Handler) handleAccountState(request network.IRequest) (err error) {
	accountState := &state.AccountState{}
	logger.Info(request.GetMessage().GetBody())
	err = accountState.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Receive Account state: \n%v", accountState))
	h.accountStateChan <- accountState
	return nil
}

/*
handleAccountState will receive receipt from connection
then print it out
*/
func (h *Handler) handleReceipt(request network.IRequest) (err error) {
	receipt := &receipt.Receipt{}
	err = receipt.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	if h.receiptChan != nil {
		h.receiptChan <- receipt
	} else {
		logger.Debug(fmt.Sprintf("Receive receipt: %v", receipt))
		logger.Debug(fmt.Sprintf("Receive To address: %v", request.GetMessage().GetToAddress()))
	}
	return nil
}

/*
handleTransactionError will receive transaction error from parent node connection
then print it out
*/
func (h *Handler) handleTransactionError(request network.IRequest) (err error) {
	transactionErr := &transaction.TransactionError{}
	err = transactionErr.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	logger.Debug("Receive transaction error:", transactionErr)

	return nil
}

func (h *Handler) handleEventLogs(request network.IRequest) error {
	eventLogs := smart_contract.EventLogs{}
	err := eventLogs.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		logger.Error("Handle Event Logs Error", err)
		return err
	}
	eventLogList := eventLogs.GetEventLogList()
	for _, eventLog := range eventLogList {
		logger.Debug("EventLogs: ", eventLog.String())
	}
	return nil
}
