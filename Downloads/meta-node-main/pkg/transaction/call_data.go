package transaction

import (
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

type ICallData interface {
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	// geter
	GetInput() []byte
}

type CallData struct {
	proto *pb.CallData
}

func NewCallData(input []byte) ICallData {
	return &CallData{
		proto: &pb.CallData{
			Input: input,
		},
	}
}

func (cd *CallData) Unmarshal(b []byte) error {
	cdPb := &pb.CallData{}
	err := proto.Unmarshal(b, cdPb)
	if err != nil {
		return err
	}
	cd.proto = cdPb
	return nil
}

func (cd *CallData) Marshal() ([]byte, error) {
	return proto.Marshal(cd.proto)
}

// geter
func (cd *CallData) GetInput() []byte {
	return cd.proto.Input
}
