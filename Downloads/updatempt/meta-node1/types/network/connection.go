package network

import "github.com/ethereum/go-ethereum/common"

type Connection interface {
	// getter
	Address() common.Address
	PublicConnectionAddress() string
	RequestChan() (chan Request, chan error)
	Type() string
	ConnectionAddress() string
	String() string

	// setter
	Init(common.Address, string, string)

	// other
	SendMessage(message Message) error
	Connect() error
	Disconnect() error
	ReadRequest()
	Clone() Connection
}
