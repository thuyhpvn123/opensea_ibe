package transaction

import (
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
)

type CallData struct {
	proto *pb.CallData
}

func NewCallData(input []byte) types.CallData {
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
func (cd *CallData) Input() []byte {
	return cd.proto.Input
}
