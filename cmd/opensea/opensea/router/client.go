package router

import (
	"sync"
	"gitlab.com/meta-node/meta-node/cmd/client/config"
	"github.com/gorilla/websocket"
	"gitlab.com/meta-node/meta-node/cmd/client/controllers"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/pkg/network"
	"gitlab.com/meta-node/meta-node/pkg/state"
	log "github.com/sirupsen/logrus"

)
type Client struct {
	ws     *websocket.Conn
	server *Server
	caller CallData
	sync.Mutex
	sendChan chan Message1
	keyPairMap          map[string]*bls.KeyPair
	config             *config.ClientConfig
	messageSenderMap      map[string]network.IMessageSender
	connectionsManager network.IConnectionsManager
	tcpServerMap          map[string]network.ISocketServer
	transactionControllerMap map[string]controllers.ITransactionController
	accountStateChan      chan state.IAccountState
}
func (client *Client) init() {
	// send init message
	// client.ws.WriteJSON(
	// 	Message{Type: "message", Msg: "Here is new client"})
	client.caller = CallData{server: client.server, client: client}
	go client.handleMessage()
	log.Info("End init client")
}
func (client *Client) handleListen() {
	for {
		// Read in a new message as JSON and map it to a Message object
		var msg map[string]interface{}
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
			break
		}
		// log.Info("Message from client: ", msg)
		client.handleCallChain(msg)
	}
}


// handle message struct tu chain tra ve va chuyen qua dang JSON gui toi cac client
func (client *Client) handleMessage() {
	for {
		msg := <-client.sendChan
		// msg1 := <-sendDataC
		log.Info(msg)
		err := client.ws.WriteJSON(msg)

		if err != nil {
			log.Printf("error: %v", err)
			client.ws.Close()
		}
	}
}
func (client *Client) handleCallChain(msg map[string]interface{}) {
	// handle call
	switch msg["command"] {
	case "get-all-wallet":
	default:
		log.Warn("Require call not match: ", msg)
	}
}