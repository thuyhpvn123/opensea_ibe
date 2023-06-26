package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

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
	defaultLogLevel     = logger.FLAG_DEBUG
	defaultConfigPath   = "config.json"
	defaultKeyPairsPath = "kp.json"
	defaultSendAmount   = "01"
	defaultToAddress    = "b31da157e543d158ccc2c46d8557b293844afcec"
	defaultWaitReceipt  = true
	defaultSkipKeyPress = false
	defaultTimeOut      = 10
)

var (
	CONFIG_FILE_PATH string
	KEYPAIRS_PATH    string
	AMOUNT           string
	LOG_LEVEL        int
	WAIT_RECEIPT     bool
	TO_ADDRESS       string
	TIME_OUT         int
	SKIP_KEY_PRESS   bool
	//
)
var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
)

type JKeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Address    string `json:"address"`
}

func main() {
	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.StringVar(&KEYPAIRS_PATH, "keypairs", defaultKeyPairsPath, "Keypairs path")
	flag.StringVar(&KEYPAIRS_PATH, "k", defaultKeyPairsPath, "Keypairs path (shorthand)")

	flag.StringVar(&AMOUNT, "amount", defaultSendAmount, "Amount to send to each account")
	flag.StringVar(&AMOUNT, "a", defaultSendAmount, "Amount to send to each account (shorthand)")

	flag.StringVar(&TO_ADDRESS, "to-address", defaultToAddress, "Eeceiver address")
	flag.StringVar(&TO_ADDRESS, "t", defaultToAddress, "Eeceiver address (shorthand)")

	flag.BoolVar(&WAIT_RECEIPT, "wait-receipt", defaultWaitReceipt, "Wait receipt before send new transaction")
	flag.BoolVar(&WAIT_RECEIPT, "w", defaultWaitReceipt, "Wait receipt before send new transaction (shorthand)")

	flag.IntVar(&TIME_OUT, "time-out", defaultTimeOut, "Wait receipt timeout duration in second")
	flag.IntVar(&TIME_OUT, "to", defaultTimeOut, "Wait receipt timeout duration in second (shorthand)")

	flag.BoolVar(&SKIP_KEY_PRESS, "y", defaultSkipKeyPress, "Wait receipt timeout duration in second (shorthand)")

	flag.Parse()
	// set logger config
	var loggerConfig = &logger.LoggerConfig{
		Flag:    LOG_LEVEL,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	config, err := c_config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)
	jKeyPairs := getKeypairs()
	amount := uint256.NewInt(0).SetBytes(common.FromHex(AMOUNT))
	if !SKIP_KEY_PRESS {
		logger.Debug(fmt.Sprintf("Running for %v accounts. Press any key to continue", len(jKeyPairs)))
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
	for _, v := range jKeyPairs {
		go sendTransactionsLoop(
			v.PrivateKey,
			cConfig.GetVersion(),
			common.HexToAddress(cConfig.ParentAddress),
			cConfig.ParentConnectionType,
			cConfig.ParentConnectionAddress,
			common.HexToAddress(TO_ADDRESS),
			amount,
		)
	}
	<-make(chan bool)
}

func getKeypairs() []JKeyPair {
	dat, _ := os.ReadFile(KEYPAIRS_PATH)
	jKeyPairs := []JKeyPair{}
	json.Unmarshal(dat, &jKeyPairs)
	return jKeyPairs
}

func sendTransactionsLoop(
	hexPrivateKey string,
	version string,
	parentAddess common.Address,
	parentConnectionType string,
	parentConnectionAddress string,
	toAddress common.Address,
	amount *uint256.Int,
) {
	keyPair := bls.NewKeyPair(common.FromHex(hexPrivateKey))
	// init connection
	messageSender := network.NewMessageSender(keyPair, version)
	connectionsManager := network.NewConnectionsManager(connectionTypesForClient)
	parentConn := network.NewConnection(
		parentAddess,
		parentConnectionType,
		parentConnectionAddress,
	)
	accountStateChan := make(chan state.IAccountState, 1)
	var receiptChan chan receipt.IReceipt
	if WAIT_RECEIPT {
		receiptChan = make(chan receipt.IReceipt, 1)
	}
	handler := c_network.NewHandler(accountStateChan, receiptChan)
	tcpServer := network.NewSockerServer(&c_config.ClientConfig{}, keyPair, connectionsManager, handler)
	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		// panic(fmt.Sprintf("error when connect to parent %v", err))
	} else {
		// init connection
		connectionsManager.AddParentConnection(parentConn)
		tcpServer.OnConnect(parentConn)
		go tcpServer.HandleConnection(parentConn)
	}

	// get account state
	messageSender.SendBytes(parentConn, command.GetAccountState, keyPair.GetAddress().Bytes(), p_common.Sign{})
	as := <-accountStateChan
	transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
	lastHash := as.GetLastHash()
	logger.Debug("lasthash", lastHash)
	pendingBalance := as.GetPendingBalance()
	maxGas := uint64(p_common.TRANSFER_GAS_COST)
	maxGasPrice := uint64(p_common.MINIMUM_BASE_FEE)
	for {
		transaction, err := transactionCtl.SendTransaction(
			lastHash,
			toAddress,
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
			<-time.After(5 * time.Millisecond)
			continue
		} else {
			logger.Info("Done send transaction from " + keyPair.GetAddress().String())
		}
		if WAIT_RECEIPT {
			select {
			case receipt := <-receiptChan:
				if receipt.GetStatus() == pb.RECEIPT_STATUS_RETURNED {
					lastHash = receipt.GetTransactionHash()
				} else {
					messageSender.SendBytes(parentConn, command.GetAccountState, keyPair.GetAddress().Bytes(), p_common.Sign{})
					as = <-accountStateChan
					lastHash = as.GetLastHash()
				}
			case <-time.After(time.Duration(TIME_OUT) * time.Second):
				messageSender.SendBytes(parentConn, command.GetAccountState, keyPair.GetAddress().Bytes(), p_common.Sign{})
				as = <-accountStateChan
				lastHash = as.GetLastHash()
			}

		} else {
			lastHash = transaction.GetHash()
			<-time.After(500 * time.Millisecond)
		}
		pendingBalance = uint256.NewInt(0)
	}
}
