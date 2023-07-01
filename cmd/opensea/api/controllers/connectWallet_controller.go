package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/meta-node/meta-node/cmd/opensea/models"
	// "go.mongodb.org/mongo-driver/mongo"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/network"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/state"
	c_network "gitlab.com/meta-node/meta-node/cmd/opensea/network"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	c_config "gitlab.com/meta-node/meta-node/cmd/opensea/config"
	"gitlab.com/meta-node/meta-node/pkg/logger"

)


type CallData struct {
	server *Server
	client *Client
}
func (caller *CallData) ConnectWallet(ctx *gin.Context) {
	var request struct {
		Address     string    `json:"address,omitempty" `
		PrivateKey string `json:"privatekey,omitempty" `
	}
	// Parse the request body into the request struct
	if err := ctx.ShouldBind(&request); err != nil {
		// Handle error parsing request
		response := models.Response{
			Code: http.StatusBadRequest,
			Data: gin.H{"error": "Invalid request"},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	call:=map[string]interface{}{
		"address":request.Address,
		"priKey":request.PrivateKey,
	}
	caller.ConnectSocket(call)

	response := models.Response{
		Code: 200,
		Data: gin.H{
			"address":request.Address,
		},
	}
	ctx.JSON(200, response)
	return
}

var (
	connectionTypesForClient = []string{p_common.NODE_CONNECTION_TYPE}
)

func (caller *CallData) ConnectSocket( walletKey map[string]interface{} ) map[string]interface{} {
	fmt.Println("ConnectSocket")
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.Config)
	result:=make(map[string]interface{})
	// connect to parent
	connectionsManager := network.NewConnectionsManager(connectionTypesForClient)

	parentConn := network.NewConnection(
		common.HexToAddress(cConfig.ParentAddress),
		cConfig.ParentConnectionType,
		cConfig.ParentConnectionAddress,
	)
	accountStateChan := make(chan state.IAccountState)
	chData :=make(chan interface{})
	handler := c_network.NewHandler(accountStateChan,chData)


	err = parentConn.Connect()
	if err != nil {
		logger.Error(fmt.Sprintf("error when connect to parent %v", err))
	} else {
		// init connection
		connectionsManager.AddParentConnection(parentConn)

		if walletKey["priKey"] == nil {
			logger.Error(fmt.Sprintf("error when GetWalletKeyFromAddress %", err))
			panic(fmt.Sprintf("error when GetWalletKeyFromAddress %v", err))
		}else{
			priKey := common.FromHex(walletKey["priKey"].(string))
			keyPair := bls.NewKeyPair(priKey)

			logger.Info("Running with key pair: " + "\n" + keyPair.String())
			messageSender := network.NewMessageSender(keyPair, config.GetVersion())
			tcpServer := network.NewSockerServer(config, keyPair, connectionsManager, handler)
			tcpServer.OnConnect(parentConn)

			go tcpServer.HandleConnection(parentConn)
	
			// init controller
			transactionCtl := controllers.NewTransactionController(keyPair, messageSender, connectionsManager)
			// init and start client
			// fmt.Println("addressString:",addressString)
			addressString := walletKey["address"].(string)
			caller.client.keyPairMap[addressString] = keyPair
			caller.client.messageSenderMap[addressString] = messageSender
			caller.client.transactionControllerMap[addressString]= transactionCtl
			caller.client.tcpServerMap[addressString]=tcpServer
			caller.client.accountStateChan=accountStateChan

		}
		
		caller.client.connectionsManager= connectionsManager
		caller.client.config = cConfig
	}
	fmt.Println("init connection")
	return result
}
