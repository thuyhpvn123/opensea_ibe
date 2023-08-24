package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/meta-node-blockchain/meta-node/pkg/bls"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"github.com/meta-node-blockchain/meta-node/types/network"
)

type SocketServer struct {
	connectionsManager network.ConnectionsManager
	listener           net.Listener
	handler            network.Handler
	config             types.Config
	keyPair            *bls.KeyPair
	ctx                context.Context
	cancelFunc         context.CancelFunc

	onConnectedCallBack    []func(network.Connection)
	onDisconnectedCallBack []func(network.Connection)
}

func NewSockerServer(
	config types.Config,
	keyPair *bls.KeyPair,
	connectionsManager network.ConnectionsManager,
	handler network.Handler,
) network.SocketServer {
	s := &SocketServer{
		config:             config,
		keyPair:            keyPair,
		connectionsManager: connectionsManager,
		handler:            handler,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}

func (s *SocketServer) AddOnConnectedCallBack(callBack func(network.Connection)) {
	s.onConnectedCallBack = append(s.onConnectedCallBack, callBack)
}

func (s *SocketServer) AddOnDisconnectedCallBack(callBack func(network.Connection)) {
	s.onDisconnectedCallBack = append(s.onDisconnectedCallBack, callBack)
}

func (s *SocketServer) Listen(listenAddress string) error {
	publicConnectionAddress := s.config.PublicConnectionAddress()
	var err error
	s.listener, err = net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}
	defer func() {
		if s.listener != nil {
			s.listener.Close()
			s.listener = nil
		}
	}()
	logger.Debug(fmt.Sprintf("Listening at %v", listenAddress))
	logger.Info(fmt.Sprintf("Public connection address: %v", publicConnectionAddress))
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			tcpConn, err := s.listener.Accept()
			if err != nil {
				logger.Warn(fmt.Sprintf("Error when accept connection %v\n", err))
				continue
			}
			conn, err := ConnectionFromTcpConnection(tcpConn)
			if err != nil {
				logger.Warn(fmt.Sprintf("error when create connection from tcp connection: %v", err))
				continue
			}
			s.OnConnect(conn)
			go s.HandleConnection(conn)
		}
	}
}

func (s *SocketServer) Stop() {
	s.cancelFunc()
}

func (s *SocketServer) OnConnect(conn network.Connection) {
	logger.Info(fmt.Sprintf("On Connect with %v type %v", conn.PublicConnectionAddress(), conn.Type()))
	SendMessage(conn, s.keyPair, p_common.InitConnection, &pb.InitConnection{
		Address:                 s.keyPair.Address().Bytes(),
		Type:                    s.config.NodeType(),
		PublicConnectionAddress: s.config.PublicConnectionAddress(),
	}, p_common.Sign{}, s.config.Version())

	for _, v := range s.onConnectedCallBack {
		v(conn)
	}
}

func (s *SocketServer) OnDisconnect(conn network.Connection) {
	logger.Warn(
		fmt.Sprintf(
			"On Disconnect with %v - address %v - type %v",
			conn.ConnectionAddress(),
			conn.Address(),
			conn.Type(),
		),
	)
	s.connectionsManager.RemoveConnection(conn)
	if conn == s.connectionsManager.ParentConnection() {
		logger.Warn("Disconnected with parent")
		// stop running if disconnected with parent
		s.Stop()
		// if connection is parent connection then retry connect
		go func(_conn network.Connection) {
			for {
				<-time.After(5 * time.Second)
				err := _conn.Connect()
				if err != nil {
					logger.Warn(fmt.Sprintf("error when retry connect to parent %v", err))
				} else {
					s.ctx, s.cancelFunc = context.WithCancel(context.Background())
					s.connectionsManager.AddParentConnection(conn)
					s.OnConnect(conn)
					go s.HandleConnection(conn)
					return
				}
			}
		}(conn)
	}

	for _, v := range s.onDisconnectedCallBack {
		v(conn)
	}

}

func (s *SocketServer) HandleConnection(conn network.Connection) error {
	logger.Debug(fmt.Sprintf("handle connection %v", conn.Address()))
	go conn.ReadRequest()
	defer func() {
		conn.Disconnect()
		s.OnDisconnect(conn)
	}()
	requestChan, errorChan := conn.RequestChan()
	for {
		select {
		case <-s.ctx.Done():
			return nil
		case request := <-requestChan:
			err := s.handler.HandleRequest(request)
			if err != nil {
				logger.Warn(fmt.Sprintf("error when process request %v", err))
				continue
			}
		case err := <-errorChan:
			if err != ErrDisconnected {
				logger.Warn(fmt.Sprintf("error when read request %v", err))
			}
			return err
		}
	}
}

func (s *SocketServer) SetKeyPair(newKeyPair *bls.KeyPair) {
	s.keyPair = newKeyPair
}
