package state_channel

import (
	"encoding/hex"
	"fmt"
	"strings"

	e_common "github.com/ethereum/go-ethereum/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type UpdateStateField struct {
	field string
	value []byte
}

func NewUpdateStateField(field string, value []byte) types.UpdateStateField {
	return &UpdateStateField{
		field: field,
		value: value,
	}
}

// general
func (f *UpdateStateField) Marshal() ([]byte, error) {
	return proto.Marshal(f.Proto())
}

func (f *UpdateStateField) Unmarshal(bytes []byte) error {
	pbField := &pb.UpdateStateField{}
	err := proto.Unmarshal(bytes, pbField)
	if err != nil {
		return err
	}
	f.FromProto(pbField)
	return nil
}

func (f *UpdateStateField) Proto() protoreflect.ProtoMessage {
	return &pb.UpdateStateField{
		Field: f.field,
		Value: f.value,
	}
}

func (f *UpdateStateField) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbField := pbMessage.(*pb.UpdateStateField)
	f.field = pbField.Field
	f.value = pbField.Value
}
func (f *UpdateStateField) String() string {
	return fmt.Sprintf("%v: %v", f.field, hex.EncodeToString(f.value))
}

// getter
func (f *UpdateStateField) Field() string {
	return f.field
}
func (f *UpdateStateField) Value() []byte {
	return f.value
}

// UpdateStateData
type StateChannelCommitData struct {
	address e_common.Address
	fields  []types.UpdateStateField
}

func NewStateChannelCommitData(
	address e_common.Address,
	fields []types.UpdateStateField,
) types.StateChannelCommitData {
	return &StateChannelCommitData{
		address: address,
		fields:  fields,
	}
}

// general
func (d *StateChannelCommitData) Marshal() ([]byte, error) {
	return proto.Marshal(d.Proto())
}

func (d *StateChannelCommitData) Unmarshal(bytes []byte) error {
	pbData := &pb.StateChannelCommitData{}
	err := proto.Unmarshal(bytes, pbData)
	if err != nil {
		return err
	}
	d.FromProto(pbData)
	return nil
}

func (d *StateChannelCommitData) Proto() protoreflect.ProtoMessage {
	pbFields := make([]*pb.UpdateStateField, len(d.fields))
	for i, field := range d.fields {
		pbFields[i] = field.Proto().(*pb.UpdateStateField)
	}
	return &pb.StateChannelCommitData{
		Address:      d.address.Bytes(),
		UpdateFields: pbFields,
	}
}

func (d *StateChannelCommitData) FromProto(pbMessage protoreflect.ProtoMessage) {
	pbData := pbMessage.(*pb.StateChannelCommitData)
	d.address = e_common.BytesToAddress(pbData.Address)
	fields := make([]types.UpdateStateField, len(pbData.UpdateFields))
	for i, pbField := range pbData.UpdateFields {
		fields[i] = &UpdateStateField{}
		fields[i].FromProto(pbField)
	}
	d.fields = fields
}

func (d *StateChannelCommitData) String() string {
	fieldStrings := make([]string, len(d.fields))
	for i, field := range d.fields {
		fieldStrings[i] = field.String()
	}
	return fmt.Sprintf("Address: %v\nUpdateFields:\n%s", d.address.Hex(), strings.Join(fieldStrings, "\n"))
}

// getter
func (d *StateChannelCommitData) Address() e_common.Address {
	return d.address
}

func (d *StateChannelCommitData) UpdateStateFields() []types.UpdateStateField {
	return d.fields
}

// setter
func (d *StateChannelCommitData) AddUpdateField(field types.UpdateStateField) {
	d.fields = append(d.fields, field)
}
