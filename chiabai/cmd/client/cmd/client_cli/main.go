package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/cmd/client/cli"
	c_config "gitlab.com/meta-node/meta-node/cmd/client/config"
	"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	c_network "gitlab.com/meta-node/meta-node/cmd/client/network"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
	"gitlab.com/meta-node/meta-node/pkg/state"
)

const (
	defaultConfigPath ="public/config/config.json"
	defaultLogLevel   = logger.FLAG_DEBUG
)

var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
	// flags
	CONFIG_FILE_PATH string
	LOG_LEVEL        int
)

func main() {
	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")
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
	handler := c_network.NewHandler(accountStateChan, nil)
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
	// init controller
	transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
	// init and start cli
	cli := cli.NewCli(
		keyPair,
		cConfig,
		messageSender,
		connectionsManager,
		transactionCtl,
		accountStateChan,
		tcpServer,
	)
	cli.Start()
}
