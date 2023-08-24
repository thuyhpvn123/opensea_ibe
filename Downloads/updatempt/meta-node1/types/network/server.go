package network

import "github.com/meta-node-blockchain/meta-node/pkg/bls"

type SocketServer interface {
	Listen(string) error
	Stop()

	OnConnect(Connection)
	OnDisconnect(Connection)

	SetKeyPair(*bls.KeyPair)

	HandleConnection(Connection) error

	AddOnConnectedCallBack(callBack func(Connection))
	AddOnDisconnectedCallBack(callBack func(Connection))
}
