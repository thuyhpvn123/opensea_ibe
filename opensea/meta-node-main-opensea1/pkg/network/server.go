package network

import (
	"context"
	"fmt"
	"net"
	"time"

	"gitlab.com/meta-node/meta-node/pkg/bls"
	"gitlab.com/meta-node/meta-node/pkg/common"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/config"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
)

type ISocketServer interface {
	Listen(string) error
	Stop()

	OnConnect(IConnection)
	OnDisconnect(IConnection)

	SetKeyPair(*bls.KeyPair)

	HandleConnection(IConnection) error
	GetHandler() chan interface{}

}

type SocketServer struct {
	connectionsManager IConnectionsManager
	listener           net.Listener
	handler            IHandler
	config             config.IConfig
	keyPair            *bls.KeyPair
	ctx                context.Context
	cancelFunc         context.CancelFunc
}

func NewSockerServer(
	config config.IConfig,
	keyPair *bls.KeyPair,
	connectionsManager IConnectionsManager,
	handler IHandler,
) ISocketServer {
	s := &SocketServer{
		config:             config,
		keyPair:            keyPair,
		connectionsManager: connectionsManager,
		handler:            handler,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}
func (s *SocketServer) GetHandler() chan interface{}{
	return s.handler.GetChData()
		
	
}
func (s *SocketServer) Listen(listenAddress string) error {
	publicConnectionAddress := s.config.GetPublicConnectionAddress()
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

func (s *SocketServer) OnConnect(conn IConnection) {
	logger.Info(fmt.Sprintf("On Connect with %v type %v", conn.GetPublicConnectionAddress(), conn.GetType()))
	SendMessage(conn, s.keyPair, common.InitConnection, &pb.InitConnection{
		Address:                 s.keyPair.GetAddress().Bytes(),
		Type:                    s.config.GetNodeType(),
		PublicConnectionAddress: s.config.GetPublicConnectionAddress(),
	}, p_common.Sign{}, s.config.GetVersion())
}

func (s *SocketServer) OnDisconnect(conn IConnection) {
	logger.Warn(
		fmt.Sprintf(
			"On Disconnect with %v - address %v - type %v",
			conn.GetConnectionAddress(),
			conn.GetAddress(),
			conn.GetType(),
		),
	)
	s.connectionsManager.RemoveConnection(conn)
	if conn == s.connectionsManager.GetParentConnection() {
		// stop running if disconnected with parent
		s.Stop()
		// if connection is parent connection then retry connect
		go func(_conn IConnection) {
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

}

func (s *SocketServer) HandleConnection(conn IConnection) error {
	logger.Debug(fmt.Sprintf("handle connection %v", conn.GetAddress()))
	go conn.ReadRequest()
	defer func() {
		conn.Disconnect()
		s.OnDisconnect(conn)
	}()
	requestChan, errorChan := conn.GetRequestChan()
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
