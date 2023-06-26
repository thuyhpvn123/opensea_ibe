package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/cmd/client/command"
	c_config "gitlab.com/meta-node/meta-node/cmd/client/config"
	"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	c_network "gitlab.com/meta-node/meta-node/cmd/client/network"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/receipt"
	"gitlab.com/meta-node/meta-node/pkg/state"
)

const (
	defaultGenerateKeyPairs = false
	defaultGenerateOutput   = "kp.json"
	defaultGenerateNumber   = 1000
	defaultConfigPath       = "config.json"
	defaultKeyPairsPath     = "kp.json"
	defaultSendAmount       = "3635c9adc5dea00000"
	defaultLogLevel         = logger.FLAG_DEBUG
)

var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
	// flags
	GENERATE_KEYPAIRS        bool
	GENERATE_KEYPAIRS_OUTPUT string
	NUMBER_GENERATE_KEYPAIRS int
	CONFIG_FILE_PATH         string
	KEYPAIRS_PATH            string
	AMOUNT                   string
	LOG_LEVEL                int
)

func main() {
	flag.BoolVar(&GENERATE_KEYPAIRS, "generate", defaultGenerateKeyPairs, "True if want to generate key pairs")
	flag.BoolVar(&GENERATE_KEYPAIRS, "g", defaultGenerateKeyPairs, "True if want to generate key pairs (shorthand)")

	flag.StringVar(&GENERATE_KEYPAIRS_OUTPUT, "generate_output", defaultGenerateOutput, "Output path for generate key pairs")
	flag.StringVar(&GENERATE_KEYPAIRS_OUTPUT, "o", defaultGenerateOutput, "Output path for generate key pairs (shorthand)")

	flag.IntVar(&NUMBER_GENERATE_KEYPAIRS, "generate_number", defaultGenerateNumber, "Number of generate key pairs")
	flag.IntVar(&NUMBER_GENERATE_KEYPAIRS, "n", defaultGenerateNumber, "Number of generate key pairs (shorthand)")

	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.StringVar(&KEYPAIRS_PATH, "keypairs", defaultKeyPairsPath, "Keypairs path")
	flag.StringVar(&KEYPAIRS_PATH, "k", defaultKeyPairsPath, "Keypairs path (shorthand)")

	flag.StringVar(&AMOUNT, "amount", defaultSendAmount, "Amount to send to each account")
	flag.StringVar(&AMOUNT, "a", defaultSendAmount, "Amount to send to each account (shorthand)")

	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.Parse()
	// set logger config
	var loggerConfig = &logger.LoggerConfig{
		Flag:    LOG_LEVEL,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	if GENERATE_KEYPAIRS {
		generateKeypairs()
		return
	}
	sendTransactions()
}

type JKeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

func sendTransactions() {
	config, err := c_config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)
	keyPair := bls.NewKeyPair(config.GetPrivateKey())
	logger.Debug("Running with key pair: " + "\n" + keyPair.String())
	// init message sender
	messageSender := network.NewMessageSender(keyPair, config.GetVersion())
	// connect to parent
	connectionsManager := network.NewConnectionsManager(connectionTypesForClient)
	// connection to parent

	parentConn := network.NewConnection(
		common.HexToAddress(cConfig.ParentAddress),
		cConfig.ParentConnectionType,
		cConfig.ParentConnectionAddress,
	)

	accountStateChan := make(chan state.IAccountState, 1)
	receiptChan := make(chan receipt.IReceipt)
	handler := c_network.NewHandler(accountStateChan, receiptChan)
	go func(chan receipt.IReceipt) {
		count := 0
		for {
			receipt := <-receiptChan
			if receipt.GetStatus() == pb.RECEIPT_STATUS_RETURNED {
				count++
			}
			logger.Info("Success transaction count ", count)
		}
	}(receiptChan)
	tcpServer := network.NewSockerServer(config, keyPair, connectionsManager, handler)
	err = parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		// panic(fmt.Sprintf("error when connect to parent %v", err))
	} else {
		// init connection
		connectionsManager.AddParentConnection(parentConn)
		tcpServer.OnConnect(parentConn)
		go tcpServer.HandleConnection(parentConn)
	}
	// read accounts
	dat, _ := os.ReadFile(KEYPAIRS_PATH)
	jAccounts := []JKeyPair{}
	json.Unmarshal(dat, &jAccounts)

	// get account state
	messageSender.SendBytes(parentConn, command.GetAccountState, keyPair.GetAddress().Bytes(), p_common.Sign{})
	as := <-accountStateChan
	logger.Debug(as)
	// init controller
	transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
	// init and start cli
	amount := uint256.NewInt(0).SetBytes(common.FromHex(AMOUNT))
	maxGas := uint64(p_common.TRANSFER_GAS_COST)
	maxGasPrice := uint64(p_common.MINIMUM_BASE_FEE)
	lastHash := as.GetLastHash()
	pendingBalance := as.GetPendingBalance()
	for _, account := range jAccounts {
		transaction, err := transactionCtl.SendTransaction(
			lastHash,
			common.HexToAddress(account.Address),
			pendingBalance,
			amount,
			maxGas,
			maxGasPrice,
			0,
			0,
			nil,
			nil,
		)
		if err != nil {
			logger.Error(err)
			panic("Transaction error")
		} else {
			logger.Info("Done send transaction to " + account.Address)
		}
		lastHash = transaction.GetHash()
		pendingBalance = uint256.NewInt(0)
	}
	logger.Info("Done")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func generateKeypairs() {
	keyPairList := make([]JKeyPair, NUMBER_GENERATE_KEYPAIRS)
	for i := 0; i < NUMBER_GENERATE_KEYPAIRS; i++ {
		rawKeyPair := bls.GenerateKeyPair()
		keyPairList[i] = JKeyPair{
			PrivateKey: hex.EncodeToString(rawKeyPair.GetPrivateKey().Bytes()),
			PublicKey:  hex.EncodeToString(rawKeyPair.GetPublicKey().Bytes()),
			Address:    hex.EncodeToString(rawKeyPair.GetAddress().Bytes()),
		}
	}
	data, _ := json.Marshal(keyPairList)
	os.WriteFile(GENERATE_KEYPAIRS_OUTPUT, data, 0644)
}
