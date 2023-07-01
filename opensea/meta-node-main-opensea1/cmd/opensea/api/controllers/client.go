package controllers

import (
	"sync"
	"gitlab.com/meta-node/meta-node/cmd/opensea/config"
	// "gitlab.com/meta-node/meta-node/cmd/opensea/api/controllers"
	controller_client"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/pkg/network"
	"gitlab.com/meta-node/meta-node/pkg/state"
	log "github.com/sirupsen/logrus"

)
type Client struct {
	server *Server
	Caller CallData
	sync.Mutex
	sendChan chan Message1
	keyPairMap          map[string]*bls.KeyPair
	config             *config.Config
	messageSenderMap      map[string]network.IMessageSender
	connectionsManager network.IConnectionsManager
	tcpServerMap          map[string]network.ISocketServer
	transactionControllerMap map[string]controller_client.ITransactionController
	accountStateChan      chan state.IAccountState
}
func (client *Client) init() (CallData){
	// send init message
	// client.ws.WriteJSON(
	// 	Message{Type: "message", Msg: "Here is new client"})
	client.Caller = CallData{server: client.server, client: client}
	log.Info("End init client")
	return client.Caller 
}

