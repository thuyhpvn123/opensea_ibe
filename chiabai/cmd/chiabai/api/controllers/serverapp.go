package controllers

import (
	"fmt"
	"net/http"
	"log"
	// "net/http"
	"sync"
	// "github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	c_config "gitlab.com/meta-node/meta-node/cmd/chiabai/config"
	"gitlab.com/meta-node/meta-node/cmd/chiabai/core"
	controller_client "gitlab.com/meta-node/meta-node/cmd/client/controllers"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	 "gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
)

type Server struct {
	sync.Mutex
	contractABI map[string]*core.ContractABI
	config      *c_config.Config
}
type Message1 struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
		contractABI: server.contractABI,
		config:      config,
	}
}

// func (server *Server) ConnectionHandler() Client {
// 	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("error when loading config %v", err))
// 		panic(fmt.Sprintf("error when loading config %v", err))
// 	}
// 	cConfig := config.(*c_config.Config)

// 	// accountStateChan := make(chan state.IAccountState, 1)
// 	client := Client{
// 		sendChan:         make(chan Message1),
// 		server:           server,
// 		keyPairMap:       make(map[string]*bls.KeyPair),
// 		config:           cConfig,
// 		messageSenderMap: make(map[string]network.IMessageSender),
// 		// connectionsManager :network.IConnectionsManager{},
// 		tcpServerMap:             make(map[string]network.ISocketServer),
// 		transactionControllerMap: make(map[string]controller_client.ITransactionController),
// 	}
// 	client.init()

// 	logger.Info("Client Connected successfully") //write on server terminal
// 	// defer server.clients.Remove(conn)
// 	return client

// }
func (server *Server) WebsocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.Config)

	// accountStateChan := make(chan state.IAccountState, 1)
	client := Client{
		ws: conn, 
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
	log.Println("Client Connected successfully") //write on server terminal
	// defer server.clients.Remove(conn)

	//listen websocket
	client.handleListen()

}

func (server *Server) getABI(wg *sync.WaitGroup, contract core.Contract) {
	var temp core.ContractABI
	temp.InitContract(contract)
	server.Lock()
	server.contractABI[contract.Name] = &temp
	server.Unlock()
	wg.Done()
}
