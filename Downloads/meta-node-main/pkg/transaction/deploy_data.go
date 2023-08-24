package transaction

import (
	"github.com/ethereum/go-ethereum/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type IDeployData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// geter
	GetCode() []byte
	GetStorageHost() string
	GetStorageAddress() common.Address
}

type DeployData struct {
	proto *pb.DeployData
}

func NewDeployData(code []byte, storageHost string, storageAddress common.Address) IDeployData {
	return &DeployData{
		proto: &pb.DeployData{
			Code:           code,
			StorageHost:    storageHost,
			StorageAddress: storageAddress.Bytes(),
		},
	}
}

func (dd *DeployData) GetStorageAddress() common.Address {
	return common.BytesToAddress(dd.proto.StorageAddress)
}
func (dd *DeployData) Unmarshal(b []byte) error {
	ddPb := &pb.DeployData{}
	err := proto.Unmarshal(b, ddPb)
	if err != nil {
		return err
	}
	dd.proto = ddPb
	return nil
}

func (dd *DeployData) Marshal() ([]byte, error) {
	return proto.Marshal(dd.proto)
}

// geter
func (dd *DeployData) GetCode() []byte {
	return dd.proto.Code
}

func (dd *DeployData) GetStorageHost() string {
	return dd.proto.StorageHost
}
