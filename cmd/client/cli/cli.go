package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/cmd/client/command"
	"gitlab.com/meta-node/meta-node/cmd/client/config"
	"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

var (
	ErrorGetAccountStateTimedOut = errors.New("get account state timed out")
	ErrorInvalidAction           = errors.New("invalid action")
)

type ICli interface {
	Start()
	Stop()
	PrintCommands()
	PrintMessage(string, string)
	SendTransaction() error
	CreateAccount()
	ChangeAccount()
	GetAccountState(common.Address) (state.IAccountState, error)
	ReadInput() string
	ReadInputAddress() common.Address
}

type Cli struct {
	keyPair            *bls.KeyPair
	config             *config.ClientConfig
	messageSender      network.IMessageSender
	connectionsManager network.IConnectionsManager
	tcpServer          network.ISocketServer
	stop               bool
	commands           map[int]string
	reader             *bufio.Reader

	transactionController controllers.ITransactionController
	accountStateChan      chan state.IAccountState
	defaultRelatedAddress map[common.Address][][]byte
}

func NewCli(
	keyPair *bls.KeyPair,
	config *config.ClientConfig,
	messageSender network.IMessageSender,
	connectionsManager network.IConnectionsManager,
	transactionController controllers.ITransactionController,
	accountStateChan chan state.IAccountState,
	tcpServer network.ISocketServer,
) ICli {
	commands := map[int]string{
		0: "Exit",
		1: "Send transaction",
		2: "Change account",
		3: "Create account",
		4: "Get account state",
		5: "Subscribe",
	}
	return &Cli{
		keyPair:               keyPair,
		config:                config,
		messageSender:         messageSender,
		connectionsManager:    connectionsManager,
		stop:                  false,
		commands:              commands,
		transactionController: transactionController,
		accountStateChan:      accountStateChan,
		tcpServer:             tcpServer,
		defaultRelatedAddress: make(map[common.Address][][]byte),
	}
}

func (cli *Cli) Start() {
	cli.reader = bufio.NewReader(os.Stdin)
	for {
		if cli.stop {
			return
		}
		cli.PrintCommands()

		command := cli.ReadInput()
		switch command {
		case "0":
		//TODO
		case "1":
			cli.SendTransaction()
		case "2":
			cli.ChangeAccount()
		case "3":
			cli.CreateAccount()
		case "4":
			cli.PrintMessage("Enter address: ", "")
			cli.GetAccountState(cli.ReadInputAddress())
			//TODO4
		case "5":
			cli.Subscribe()
		}
	}
}

func (cli *Cli) Subscribe() {
	cli.PrintMessage("Enter smart contract storage host: ", "")
	storageHost := cli.ReadInput()
	cli.PrintMessage("Enter smart contract address: ", "")
	contractAddress := cli.ReadInput()

	storageConnection := network.NewConnection(cli.keyPair.GetAddress(), p_common.STORAGE_CONNECTION_TYPE, storageHost)
	err := storageConnection.Connect()
	if err != nil {
		logger.Error("Subscribe fail", err)
		return
	}
	go cli.tcpServer.HandleConnection(storageConnection)

	err = cli.messageSender.SendBytes(storageConnection, command.SubscribeToAddress, common.HexToAddress(contractAddress).Bytes(), p_common.Sign{})
	if err != nil {
		logger.Error("Subscribe fail", err)
	}
	logger.Debug("Subscribe address: ", contractAddress)
}

// TODE Cli stop
func (cli *Cli) Stop() {

}

func (cli *Cli) PrintCommands() {
	str := p_common.Cyan + "======= Commands =======\n" + p_common.Purple
	for i := 0; i < len(cli.commands); i++ {
		str += fmt.Sprintf("%v: %v\n", i, cli.commands[i])
	}
	str += p_common.Reset
	fmt.Print(str)
}

func (cli *Cli) SendTransaction() error {
	cli.PrintMessage("Enter to address: ", "")
	toAddress := cli.ReadInputAddress()
	cli.PrintMessage("Enter to amount (hex): ", "")
	hexAmount := cli.ReadInput()
	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
	cli.PrintMessage(`Enter action (default 0):
	0: None
	1: Stake
	2: Unstake
	3: Deploy smart contract
	4: Call smart contract`, "")

	actionStr := cli.ReadInput()
	var action pb.ACTION
	if actionStr == "" {
		action = 0
	} else {
		actionI, _ := strconv.Atoi(actionStr)
		action = pb.ACTION(int32(actionI))
	}
	if action < 0 || action > 4 {
		return ErrorInvalidAction
	}

	var data []byte
	if action == pb.ACTION_STAKE {
		cli.PrintMessage("Enter connection address: ", "")
		connectionAddress := cli.ReadInput()
		data = []byte(connectionAddress)
	}
	if action == pb.ACTION_UNSTAKE {
		cli.PrintMessage("Enter unstake amount(hex): ", "")
		data = common.FromHex(cli.ReadInput())
	}

	var err error
	as, err := cli.GetAccountState(cli.keyPair.GetAddress())
	if err != nil {
		return err
	}

	if action == pb.ACTION_DEPLOY_SMART_CONTRACT {
		data, err = cli.getDataForDeploySmartContract()
		if err != nil {
			panic(err)
		}
		toAddress = common.BytesToAddress(
			crypto.Keccak256(
				append(
					as.GetAddress().Bytes(),
					as.GetLastHash().Bytes()...),
			)[12:],
		)
	}

	if action == pb.ACTION_CALL_SMART_CONTRACT {
		data, err = cli.getDataForCallSmartContract()
		if err != nil {
			panic(err)
		}
	}

	var relatedAddresses [][]byte
	if action == pb.ACTION_CALL_SMART_CONTRACT || action == pb.ACTION_DEPLOY_SMART_CONTRACT {
		relatedAddresses = cli.ReadRelatedAddress(toAddress)
	}

	cli.PrintMessage("Enter max gas (default 500000): ", "")
	maxGas, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		maxGas = 500000
	}

	cli.PrintMessage("Enter max gas price in gwei (default 10 gwei): ", "")
	maxGasPriceGwei, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		maxGasPriceGwei = 10
	}
	maxGasPrice := 1000000000 * maxGasPriceGwei

	cli.PrintMessage("Enter max time use in milli second (default 1000 milli second): ", "")
	maxTimeUse, err := strconv.ParseUint(cli.ReadInput(), 10, 64)
	if err != nil {
		maxTimeUse = 1000
	}

	transaction, err := cli.transactionController.SendTransaction(
		as.GetLastHash(),
		toAddress,
		as.GetPendingBalance(),
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		action,
		data,
		relatedAddresses,
	)
	logger.Debug("Sending transaction", transaction)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func (cli *Cli) ChangeAccount() {
	cli.PrintMessage("Enter to private key(hex): ", "")
	hexPrivateKey := cli.ReadInput()
	keyPair := bls.NewKeyPair(common.FromHex(hexPrivateKey))
	cli.keyPair = keyPair
	cli.messageSender.SetKeyPair(keyPair)
	cli.transactionController.SetKeyPair(keyPair)
	cli.tcpServer.SetKeyPair(keyPair)
	// disconnect parent connection and init new conneciton

	parentConn := cli.connectionsManager.GetParentConnection()
	if parentConn != nil {
		newParentConn := parentConn.Clone()
		cli.connectionsManager.AddParentConnection(newParentConn)
		parentConn.Disconnect()
		err := newParentConn.Connect()
		if err != nil {
			logger.Warn("error when connect to parent", err)
		} else {
			cli.tcpServer.OnConnect(newParentConn)
			go cli.tcpServer.HandleConnection(newParentConn)
		}
	}
	logger.Info("Running with key pair:", keyPair)
}

func (cli *Cli) CreateAccount() {
	keyPair := bls.GenerateKeyPair()
	logger.Info(fmt.Sprintf("Key pair:\n%v", keyPair))
}

func (cli *Cli) ReadInput() string {
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Replace(input, "\n", "", -1)
	return input
}

func (cli *Cli) ReadInputAddress() common.Address {
	input := cli.ReadInput()
	address := common.HexToAddress(input)
	return address
}

func (cli *Cli) PrintMessage(message string, color string) {
	if color == "" {
		color = p_common.Purple
	}
	fmt.Printf(color+"%v\n"+p_common.Reset, message)
}

func (cli *Cli) GetAccountState(address common.Address) (state.IAccountState, error) {
	parentConn := cli.connectionsManager.GetParentConnection()
	cli.messageSender.SendBytes(parentConn, command.GetAccountState, address.Bytes(), p_common.Sign{})
	select {
	case accountState := <-cli.accountStateChan:
		return accountState, nil
	case <-time.After(2 * time.Second):
		return nil, ErrorGetAccountStateTimedOut
	}
}

func (cli *Cli) getDataForDeploySmartContract() ([]byte, error) {
	cli.PrintMessage("Enter to smart contract file name (in contracts folder): ", "")
	contractFileName := cli.ReadInput()
	b, _ := os.ReadFile("./contracts/" + contractFileName)
	cli.PrintMessage("Enter smart contract storage host: ", "")
	contractStorageHost := cli.ReadInput()
	if contractStorageHost == "" {
		contractStorageHost = "127.0.0.1:3051"
	}
	cli.PrintMessage("Enter smart contract storage address: ", "")
	contractStorageAddress := cli.ReadInput()
	if contractStorageAddress == "" {
		contractStorageAddress = "da7284fac5e804f8b9d71aa39310f0f86776b51d"
	}
	deployData := transaction.NewDeployData(common.FromHex(string(b)), contractStorageHost, common.HexToAddress(contractStorageAddress))
	return deployData.Marshal()
}

func (cli *Cli) getDataForCallSmartContract() ([]byte, error) {
	cli.PrintMessage("Enter to input for call smart contract (hex): ", "")
	input := cli.ReadInput()
	callData := transaction.NewCallData(common.FromHex(input))
	return callData.Marshal()
}

func (cli *Cli) ReadRelatedAddress(smartcontractAddress common.Address) [][]byte {
	cli.PrintMessage("Enter Related Address: ", "")
	stringRelatedAddresses := cli.ReadInput()
	if stringRelatedAddresses == "" {
		if cli.defaultRelatedAddress[smartcontractAddress] == nil {
			return [][]byte{}
		}
		return cli.defaultRelatedAddress[smartcontractAddress]
	}
	hexRelatedAddresses := strings.Split(stringRelatedAddresses, ",")
	relatedAddresses := make([][]byte, len(hexRelatedAddresses))
	logger.Debug("Temp Related Address")
	for idx, hexAddress := range hexRelatedAddresses {
		address := common.HexToAddress(hexAddress)
		logger.Debug(address)
		relatedAddresses[idx] = address.Bytes()
	}
	cli.defaultRelatedAddress[smartcontractAddress] = append(cli.defaultRelatedAddress[smartcontractAddress], relatedAddresses...)
	return relatedAddresses
}
