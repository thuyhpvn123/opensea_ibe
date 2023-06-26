package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

var (
	ErrDisconnected         = errors.New("error connection disconnected")
	ErrInvalidMessageLength = errors.New("invalid message length")
	ErrExceedMessageLength  = errors.New("exceed message length")
	ErrNilConnection        = errors.New("nil connection")
)

type IConnection interface {
	SendMessage(message IMessage) error
	Connect() error
	Disconnect() error
	GetAddress() common.Address
	GetConnectionAddress() string
	GetPublicConnectionAddress() string
	GetType() string
	Init(common.Address, string, string)
	ReadRequest()
	GetRequestChan() (chan IRequest, chan error)
	Clone() IConnection
	String() string
}

func ConnectionFromTcpConnection(tcpConn net.Conn) (IConnection, error) {
	connectionAddress := tcpConn.RemoteAddr().String()
	return &Connection{
		address:           common.Address{},
		cType:             "",
		tcpConn:           tcpConn,
		connectionAddress: connectionAddress,
		requestChan:       make(chan IRequest, 10000),
		errorChan:         make(chan error),
	}, nil
}

func NewConnection(
	address common.Address,
	cType string,
	publicConnectionAddress string,
) IConnection {
	return &Connection{
		address:                 address,
		cType:                   cType,
		publicConnectionAddress: publicConnectionAddress,
		requestChan:             make(chan IRequest, 10000),
		errorChan:               make(chan error),
	}
}

type Connection struct {
	mu      sync.Mutex
	address common.Address
	cType   string

	publicConnectionAddress string
	connectionAddress       string

	requestChan chan IRequest
	errorChan   chan error
	tcpConn     net.Conn
}

func (c *Connection) GetConnectionAddress() string {
	return c.connectionAddress
}

func (c *Connection) GetPublicConnectionAddress() string {
	return c.publicConnectionAddress
}

func (c *Connection) SendMessage(message IMessage) error {
	if c == nil {
		return ErrNilConnection
	}
	b, err := message.Marshal()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(b)))
	_, err = c.tcpConn.Write(length)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Connect() (err error) {

	if c.publicConnectionAddress == "" {
		logger.Info("dialing " + c.connectionAddress)
		c.tcpConn, err = net.Dial("tcp", c.connectionAddress)
	} else {
		logger.Info("dialing " + c.publicConnectionAddress)
		c.tcpConn, err = net.Dial("tcp", c.publicConnectionAddress)
	}

	return err
}

func (c *Connection) Disconnect() error {
	return c.tcpConn.Close()
}

func (c *Connection) GetAddress() common.Address {
	return c.address
}

func (c *Connection) GetType() string {
	return c.cType
}

func (c *Connection) GetRequestChan() (chan IRequest, chan error) {
	return c.requestChan, c.errorChan
}

func (c *Connection) ReadRequest() {
	for {
		bLength := make([]byte, 8)
		_, err := io.ReadFull(c.tcpConn, bLength)
		if err != nil {
			switch err {
			case io.EOF:
				c.errorChan <- ErrDisconnected
			default:
				c.errorChan <- err
			}
			return
		}
		messageLength := binary.LittleEndian.Uint64(bLength)
		start := time.Now()
		maxMsgLength := uint64(1073741824)
		if messageLength > maxMsgLength {
			c.errorChan <- ErrExceedMessageLength
			return
		}

		data := make([]byte, messageLength)
		byteRead, err := io.ReadFull(c.tcpConn, data)
		if err != nil {
			switch err {
			case io.EOF:
				c.errorChan <- ErrDisconnected
			default:
				c.errorChan <- err
			}
			return

		}

		if uint64(byteRead) != messageLength {
			c.errorChan <- ErrExceedMessageLength
			return
		}

		msg := &pb.Message{}
		err = proto.Unmarshal(data[:messageLength], msg)
		if err != nil {
			c.errorChan <- err
			return
		}

		c.requestChan <- NewRequest(c, NewMessage(msg))
		logger.Trace("Process time for read request: " + time.Since(start).String())
	}
}

func (c *Connection) Init(address common.Address, cType string, publicConnectionAddress string) {
	logger.Info(fmt.Sprintf("Init connection with address %v type %v at %v", address, cType, publicConnectionAddress))
	c.address = address
	c.cType = cType
	c.publicConnectionAddress = publicConnectionAddress
	c.connectionAddress = c.tcpConn.RemoteAddr().String()
}

func (c *Connection) String() string {
	return fmt.Sprintf(
		`Address: %v 
		Type %v
		ConnectionAddress %v`,
		c.address, c.cType, c.publicConnectionAddress)
}

func (c *Connection) Clone() IConnection {
	newConn := NewConnection(
		c.address,
		c.cType,
		c.publicConnectionAddress,
	)
	return newConn
}
