package state

import (
	"fmt"

	e_common "github.com/ethereum/go-ethereum/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type StateChannelState struct {
	validatorAddresses []e_common.Address
}

func NewStateChannelState(validatorAddresses []e_common.Address) types.StateChannelState {
	return &StateChannelState{
		validatorAddresses: validatorAddresses,
	}
}

// general
func (scs *StateChannelState) Unmarshal(bytes []byte) error {
	pbState := &pb.StateChannelState{}
	err := proto.Unmarshal(bytes, pbState)
	if err != nil {
		return err
	}
	scs.FromProto(pbState)
	return nil
}

func (scs *StateChannelState) Marshal() ([]byte, error) {
	return proto.Marshal(scs.Proto())
}

func (scs *StateChannelState) Proto() protoreflect.ProtoMessage {
	bAddresses := make([][]byte, len(scs.validatorAddresses))
	for i, v := range scs.validatorAddresses {
		bAddresses[i] = v.Bytes()
	}

	return &pb.StateChannelState{
		ValidatorAddress: bAddresses,
	}
}

func (scs *StateChannelState) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbState := pbMessage.(*pb.StateChannelState)
	addresses := make([]e_common.Address, len(pbState.ValidatorAddress))
	for i, bAddress := range pbState.ValidatorAddress {
		addresses[i] = e_common.BytesToAddress(bAddress)
	}
	scs.validatorAddresses = addresses
}

func (scs *StateChannelState) String() string {
	addressStrings := make([]string, len(scs.validatorAddresses))
	for i, address := range scs.validatorAddresses {
		addressStrings[i] = address.Hex()
	}
	return fmt.Sprintf("ValidatorAddresses: %v", addressStrings)
}

func (scs *StateChannelState) ValidatorAddresses() []e_common.Address {
	return scs.validatorAddresses
}
