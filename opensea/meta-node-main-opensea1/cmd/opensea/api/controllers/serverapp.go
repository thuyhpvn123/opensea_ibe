package controllers

import (
	"fmt"
	// "log"
	// "net/http"
	"sync"
	// "github.com/gin-gonic/gin"

	log "gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/pkg/network"
	"gitlab.com/meta-node/meta-node/cmd/opensea/core"
	c_config "gitlab.com/meta-node/meta-node/cmd/opensea/config"
	controller_client"gitlab.com/meta-node/meta-node/cmd/client/controllers"

)
type Server struct {
	sync.Mutex
	contractABI       map[string]*core.ContractABI
	config            *c_config.Config
}
type Message1 struct {
	Command string       `json:"command"`
	Data interface{}   `json:"data"`
}


func (server *Server) Init(config *c_config.Config) *Server {
	// init subscriber
	server.config = config
	server.contractABI = make(map[string]*core.ContractABI)
	var wg sync.WaitGroup
	for _, contract := range core.Contracts {
		wg.Add(1)
		go server.getABI(&wg, contract)
	}
	wg.Wait()

	fmt.Println("the end")
	return &Server{
		contractABI:       server.contractABI,
		config:            config,
	}
}

func (server *Server) ConnectionHandler()(Client) {
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		log.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.Config)

	// accountStateChan := make(chan state.IAccountState, 1)
	client := Client{
		sendChan: make(chan Message1),
		server: server,
		keyPairMap  : make(map[string]*bls.KeyPair),
		config  :cConfig,
		messageSenderMap : make(map[string]network.IMessageSender),
		// connectionsManager :network.IConnectionsManager{},
		tcpServerMap      : make(map[string]network.ISocketServer),
		transactionControllerMap :make(map[string]controller_client.ITransactionController),
	}
	client.init()

	log.Info("Client Connected successfully") //write on server terminal
	// defer server.clients.Remove(conn)
	return client

}
func (server *Server) getABI(wg *sync.WaitGroup, contract core.Contract) {
	var temp core.ContractABI
	temp.InitContract(contract)
	server.Lock()
	server.contractABI[contract.Name] = &temp
	server.Unlock()
	wg.Done()
}