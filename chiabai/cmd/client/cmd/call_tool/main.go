package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

const (
	defaultLogLevel     = logger.FLAG_DEBUG
	defaultConfigPath   = "config.json"
	defaultDataFile     = "data.json"
	defaultSkipKeyPress = false
)

var (
	CONFIG_FILE_PATH string
	DATA_FILE_PATH   string
	LOG_LEVEL        int
	SKIP_KEY_PRESS   bool
	//
)
var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
)

type SCData struct {
	Action         string   `json:"action"`
	Input          string   `json:"input"`
	Address        string   `json:"address"`
	RelatedAddress []string `json:"related_address"`
	StorageHost    string   `json:"storage_host"`
	StorageAddress string   `json:"storage_address"`
	Amount string   `json:"amount"`

}

func main() {
	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.StringVar(&DATA_FILE_PATH, "data", defaultDataFile, "Data file path")
	flag.StringVar(&DATA_FILE_PATH, "d", defaultDataFile, "Data file path (shorthand)")

	flag.BoolVar(&SKIP_KEY_PRESS, "skip", defaultSkipKeyPress, "Skip press to run new transaction")
	flag.BoolVar(&SKIP_KEY_PRESS, "s", defaultSkipKeyPress, "Skip press to run new transaction (shorthand)")

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

	sendTransactionsLoop(
		cConfig,
	)

	logger.Debug("Done. Press any key to exist")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func sendTransactionsLoop(
	config *c_config.ClientConfig,
) {
	keyPair := bls.NewKeyPair(common.FromHex(config.PrivateKey))
	// init connection
	messageSender := network.NewMessageSender(keyPair, config.Version)
	connectionsManager := network.NewConnectionsManager(connectionTypesForClient)
	parentConn := network.NewConnection(
		common.HexToAddress(config.ParentAddress),
		config.ParentConnectionType,
		config.ParentConnectionAddress,
	)
	accountStateChan := make(chan state.IAccountState, 1)
	receiptChan := make(chan receipt.IReceipt, 1)

	handler := c_network.NewHandler(accountStateChan, receiptChan)
	tcpServer := network.NewSockerServer(config, keyPair, connectionsManager, handler)
	err := parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
		panic(fmt.Sprintf("error when connect to parent %v", err))
	} else {
		// init connection
		connectionsManager.AddParentConnection(parentConn)
		tcpServer.OnConnect(parentConn)
		go tcpServer.HandleConnection(parentConn)
	}
	transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
	datas := getDatas()
	logger.Info("Total request", len(datas))
	for _, v := range datas {
		// get account state
		messageSender.SendBytes(parentConn, command.GetAccountState, keyPair.GetAddress().Bytes(), p_common.Sign{})
		as := <-accountStateChan
		lastHash := as.GetLastHash()
		logger.Debug("lasthash", lastHash)
		pendingBalance := as.GetPendingBalance()
		maxGas := uint64(2000000)
		maxGasPrice := uint64(p_common.MINIMUM_BASE_FEE)
		// amount := uint256.NewInt(0)
		amount := uint256.NewInt(0).SetBytes(common.FromHex(v.Amount))
		var action pb.ACTION
		var toAddress common.Address
		var bData []byte
		if v.Action == "deploy" {
			action = pb.ACTION_DEPLOY_SMART_CONTRACT
			toAddress = common.BytesToAddress(
				crypto.Keccak256(
					append(
						as.GetAddress().Bytes(),
						as.GetLastHash().Bytes()...),
				)[12:],
			)
			deployData := transaction.NewDeployData(
				common.FromHex(v.Input),
				v.StorageHost,
				common.HexToAddress(v.StorageAddress),
			)
			bData, err = deployData.Marshal()
			if err != nil {
				panic(err)
			}
		} else {
			action = pb.ACTION_CALL_SMART_CONTRACT
			toAddress = common.HexToAddress(v.Address)
			callData := transaction.NewCallData(common.FromHex(v.Input))
			bData, err = callData.Marshal()
			if err != nil {
				panic(err)
			}
		}
		relatedADdress := make([][]byte, len(v.RelatedAddress))
		for i, v := range v.RelatedAddress {
			relatedADdress[i] = common.FromHex(v)
		}
		_, err := transactionCtl.SendTransaction(
			lastHash,
			toAddress,
			pendingBalance,
			amount,
			maxGas,
			maxGasPrice,
			1000,
			action,
			bData,
			relatedADdress,
		)
		if err != nil {
			logger.Error(err)
			<-time.After(5 * time.Millisecond)
			continue
		} else {
			logger.Info("Done send transaction from " + keyPair.GetAddress().String())
		}

		receipt := <-receiptChan
		logger.Info("Receive receipt", receipt)
		if receipt.GetStatus() != pb.RECEIPT_STATUS_RETURNED && receipt.GetStatus() != pb.RECEIPT_STATUS_HALTED {
			panic("Fail transaction")
		}
		if !SKIP_KEY_PRESS {
			logger.Debug("Press any key to continue")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
		}
	}

}

func getDatas() []SCData {
	dat, _ := os.ReadFile(DATA_FILE_PATH)
	fmt.Println("getData:",string(dat))
	scDatas := []SCData{}
	err := json.Unmarshal(dat, &scDatas)
	if err != nil {
		panic(err)
	}
	return scDatas
}
