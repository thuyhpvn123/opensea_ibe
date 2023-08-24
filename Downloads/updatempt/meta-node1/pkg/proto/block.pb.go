// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: block.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BLOCK_TYPE int32

const (
	BLOCK_TYPE_SUCCESS BLOCK_TYPE = 0
	BLOCK_TYPE_FAIL    BLOCK_TYPE = 1
)

// Enum value maps for BLOCK_TYPE.
var (
	BLOCK_TYPE_name = map[int32]string{
		0: "SUCCESS",
		1: "FAIL",
	}
	BLOCK_TYPE_value = map[string]int32{
		"SUCCESS": 0,
		"FAIL":    1,
	}
)

func (x BLOCK_TYPE) Enum() *BLOCK_TYPE {
	p := new(BLOCK_TYPE)
	*p = x
	return p
}

func (x BLOCK_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BLOCK_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_block_proto_enumTypes[0].Descriptor()
}

func (BLOCK_TYPE) Type() protoreflect.EnumType {
	return &file_block_proto_enumTypes[0]
}

func (x BLOCK_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BLOCK_TYPE.Descriptor instead.
func (BLOCK_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{0}
}

type BlockHashData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number            []byte     `protobuf:"bytes,1,opt,name=Number,proto3" json:"Number,omitempty"`
	Type              BLOCK_TYPE `protobuf:"varint,2,opt,name=Type,proto3,enum=block.BLOCK_TYPE" json:"Type,omitempty"`
	LastEntryHash     []byte     `protobuf:"bytes,3,opt,name=LastEntryHash,proto3" json:"LastEntryHash,omitempty"`
	LeaderAddress     []byte     `protobuf:"bytes,4,opt,name=LeaderAddress,proto3" json:"LeaderAddress,omitempty"`
	AccountStatesRoot []byte     `protobuf:"bytes,5,opt,name=AccountStatesRoot,proto3" json:"AccountStatesRoot,omitempty"`
	ReceiptRoot       []byte     `protobuf:"bytes,6,opt,name=ReceiptRoot,proto3" json:"ReceiptRoot,omitempty"`
	BaseFee           uint64     `protobuf:"varint,7,opt,name=BaseFee,proto3" json:"BaseFee,omitempty"`
	GasLimit          uint64     `protobuf:"varint,8,opt,name=GasLimit,proto3" json:"GasLimit,omitempty"`
	TimeStamp         uint64     `protobuf:"varint,9,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	StakeStatesRoot   []byte     `protobuf:"bytes,10,opt,name=StakeStatesRoot,proto3" json:"StakeStatesRoot,omitempty"`
}

func (x *BlockHashData) Reset() {
	*x = BlockHashData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHashData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHashData) ProtoMessage() {}

func (x *BlockHashData) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHashData.ProtoReflect.Descriptor instead.
func (*BlockHashData) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{0}
}

func (x *BlockHashData) GetNumber() []byte {
	if x != nil {
		return x.Number
	}
	return nil
}

func (x *BlockHashData) GetType() BLOCK_TYPE {
	if x != nil {
		return x.Type
	}
	return BLOCK_TYPE_SUCCESS
}

func (x *BlockHashData) GetLastEntryHash() []byte {
	if x != nil {
		return x.LastEntryHash
	}
	return nil
}

func (x *BlockHashData) GetLeaderAddress() []byte {
	if x != nil {
		return x.LeaderAddress
	}
	return nil
}

func (x *BlockHashData) GetAccountStatesRoot() []byte {
	if x != nil {
		return x.AccountStatesRoot
	}
	return nil
}

func (x *BlockHashData) GetReceiptRoot() []byte {
	if x != nil {
		return x.ReceiptRoot
	}
	return nil
}

func (x *BlockHashData) GetBaseFee() uint64 {
	if x != nil {
		return x.BaseFee
	}
	return 0
}

func (x *BlockHashData) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *BlockHashData) GetTimeStamp() uint64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *BlockHashData) GetStakeStatesRoot() []byte {
	if x != nil {
		return x.StakeStatesRoot
	}
	return nil
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash              []byte     `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Number            []byte     `protobuf:"bytes,2,opt,name=Number,proto3" json:"Number,omitempty"`
	Type              BLOCK_TYPE `protobuf:"varint,3,opt,name=Type,proto3,enum=block.BLOCK_TYPE" json:"Type,omitempty"`
	LastEntryHash     []byte     `protobuf:"bytes,4,opt,name=LastEntryHash,proto3" json:"LastEntryHash,omitempty"`
	LeaderAddress     []byte     `protobuf:"bytes,5,opt,name=LeaderAddress,proto3" json:"LeaderAddress,omitempty"`
	AccountStatesRoot []byte     `protobuf:"bytes,6,opt,name=AccountStatesRoot,proto3" json:"AccountStatesRoot,omitempty"`
	ReceiptRoot       []byte     `protobuf:"bytes,7,opt,name=ReceiptRoot,proto3" json:"ReceiptRoot,omitempty"`
	BaseFee           uint64     `protobuf:"varint,8,opt,name=BaseFee,proto3" json:"BaseFee,omitempty"`
	GasLimit          uint64     `protobuf:"varint,9,opt,name=GasLimit,proto3" json:"GasLimit,omitempty"`
	TimeStamp         uint64     `protobuf:"varint,10,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	StakeStatesRoot   []byte     `protobuf:"bytes,11,opt,name=StakeStatesRoot,proto3" json:"StakeStatesRoot,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{1}
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetNumber() []byte {
	if x != nil {
		return x.Number
	}
	return nil
}

func (x *Block) GetType() BLOCK_TYPE {
	if x != nil {
		return x.Type
	}
	return BLOCK_TYPE_SUCCESS
}

func (x *Block) GetLastEntryHash() []byte {
	if x != nil {
		return x.LastEntryHash
	}
	return nil
}

func (x *Block) GetLeaderAddress() []byte {
	if x != nil {
		return x.LeaderAddress
	}
	return nil
}

func (x *Block) GetAccountStatesRoot() []byte {
	if x != nil {
		return x.AccountStatesRoot
	}
	return nil
}

func (x *Block) GetReceiptRoot() []byte {
	if x != nil {
		return x.ReceiptRoot
	}
	return nil
}

func (x *Block) GetBaseFee() uint64 {
	if x != nil {
		return x.BaseFee
	}
	return 0
}

func (x *Block) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *Block) GetTimeStamp() uint64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *Block) GetStakeStatesRoot() []byte {
	if x != nil {
		return x.StakeStatesRoot
	}
	return nil
}

type FullBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Block               *Block                  `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	AccountStateChanges []*AccountState         `protobuf:"bytes,2,rep,name=AccountStateChanges,proto3" json:"AccountStateChanges,omitempty"`
	Receipts            []*Receipt              `protobuf:"bytes,3,rep,name=Receipts,proto3" json:"Receipts,omitempty"`
	ValidatorSigns      map[string][]byte       `protobuf:"bytes,4,rep,name=ValidatorSigns,proto3" json:"ValidatorSigns,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NextLeaderAddress   []byte                  `protobuf:"bytes,5,opt,name=NextLeaderAddress,proto3" json:"NextLeaderAddress,omitempty"`
	StakeStateChanges   map[string]*StakeStates `protobuf:"bytes,6,rep,name=StakeStateChanges,proto3" json:"StakeStateChanges,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FullBlock) Reset() {
	*x = FullBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FullBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FullBlock) ProtoMessage() {}

func (x *FullBlock) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FullBlock.ProtoReflect.Descriptor instead.
func (*FullBlock) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{2}
}

func (x *FullBlock) GetBlock() *Block {
	if x != nil {
		return x.Block
	}
	return nil
}

func (x *FullBlock) GetAccountStateChanges() []*AccountState {
	if x != nil {
		return x.AccountStateChanges
	}
	return nil
}

func (x *FullBlock) GetReceipts() []*Receipt {
	if x != nil {
		return x.Receipts
	}
	return nil
}

func (x *FullBlock) GetValidatorSigns() map[string][]byte {
	if x != nil {
		return x.ValidatorSigns
	}
	return nil
}

func (x *FullBlock) GetNextLeaderAddress() []byte {
	if x != nil {
		return x.NextLeaderAddress
	}
	return nil
}

func (x *FullBlock) GetStakeStateChanges() map[string]*StakeStates {
	if x != nil {
		return x.StakeStateChanges
	}
	return nil
}

type ConfirmBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash              []byte            `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Number            []byte            `protobuf:"bytes,2,opt,name=Number,proto3" json:"Number,omitempty"`
	AccountStatesRoot []byte            `protobuf:"bytes,3,opt,name=AccountStatesRoot,proto3" json:"AccountStatesRoot,omitempty"`
	TimeStamp         uint64            `protobuf:"varint,4,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	ValidatorSigns    map[string][]byte `protobuf:"bytes,5,rep,name=ValidatorSigns,proto3" json:"ValidatorSigns,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	NextLeaderAddress []byte            `protobuf:"bytes,6,opt,name=NextLeaderAddress,proto3" json:"NextLeaderAddress,omitempty"`
}

func (x *ConfirmBlock) Reset() {
	*x = ConfirmBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmBlock) ProtoMessage() {}

func (x *ConfirmBlock) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmBlock.ProtoReflect.Descriptor instead.
func (*ConfirmBlock) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{3}
}

func (x *ConfirmBlock) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *ConfirmBlock) GetNumber() []byte {
	if x != nil {
		return x.Number
	}
	return nil
}

func (x *ConfirmBlock) GetAccountStatesRoot() []byte {
	if x != nil {
		return x.AccountStatesRoot
	}
	return nil
}

func (x *ConfirmBlock) GetTimeStamp() uint64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *ConfirmBlock) GetValidatorSigns() map[string][]byte {
	if x != nil {
		return x.ValidatorSigns
	}
	return nil
}

func (x *ConfirmBlock) GetNextLeaderAddress() []byte {
	if x != nil {
		return x.NextLeaderAddress
	}
	return nil
}

var File_block_proto protoreflect.FileDescriptor

var file_block_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x1a, 0x13, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x72, 0x65, 0x63, 0x65, 0x69,
	0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe8, 0x02, 0x0a, 0x0d, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x25, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x11, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x61, 0x73,
	0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x24, 0x0a, 0x0d, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x6f, 0x6f, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f,
	0x6f, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70,
	0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x42, 0x61, 0x73, 0x65, 0x46, 0x65, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x42, 0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x47, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x47, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54,
	0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x28, 0x0a, 0x0f, 0x53, 0x74, 0x61,
	0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0f, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x6f, 0x6f, 0x74, 0x22, 0xf4, 0x02, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e,
	0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x4c, 0x61, 0x73, 0x74, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x48, 0x61, 0x73, 0x68, 0x12, 0x24, 0x0a, 0x0d, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x4c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x11,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0b, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x42, 0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x42,
	0x61, 0x73, 0x65, 0x46, 0x65, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x47, 0x61, 0x73, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x47, 0x61, 0x73, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x28, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x6f, 0x6f, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x53, 0x74, 0x61, 0x6b, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x22, 0xa4, 0x04, 0x0a, 0x09, 0x46,
	0x75, 0x6c, 0x6c, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x22, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x4d, 0x0a, 0x13,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x13, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x08, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x52,
	0x08, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x73, 0x12, 0x4c, 0x0a, 0x0e, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67,
	0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x4e, 0x65, 0x78, 0x74, 0x4c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x11, 0x4e, 0x65, 0x78, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x55, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x27, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x2e, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x11, 0x53, 0x74, 0x61, 0x6b, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x1a, 0x41, 0x0a, 0x13,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a,
	0x60, 0x0a, 0x16, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x6b, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0xc8, 0x02, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x2c,
	0x0a, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x6f, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x73, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x4f, 0x0a, 0x0e, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x72, 0x6d, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f,
	0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x4e,
	0x65, 0x78, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x4e, 0x65, 0x78, 0x74, 0x4c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x1a, 0x41, 0x0a, 0x13, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x23, 0x0a, 0x0a,
	0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41, 0x49, 0x4c, 0x10,
	0x01, 0x42, 0x2d, 0x0a, 0x23, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x6e, 0x6f,
	0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c,
	0x65, 0x64, 0x2e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_proto_rawDescOnce sync.Once
	file_block_proto_rawDescData = file_block_proto_rawDesc
)

func file_block_proto_rawDescGZIP() []byte {
	file_block_proto_rawDescOnce.Do(func() {
		file_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_proto_rawDescData)
	})
	return file_block_proto_rawDescData
}

var file_block_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_block_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_block_proto_goTypes = []interface{}{
	(BLOCK_TYPE)(0),       // 0: block.BLOCK_TYPE
	(*BlockHashData)(nil), // 1: block.BlockHashData
	(*Block)(nil),         // 2: block.Block
	(*FullBlock)(nil),     // 3: block.FullBlock
	(*ConfirmBlock)(nil),  // 4: block.ConfirmBlock
	nil,                   // 5: block.FullBlock.ValidatorSignsEntry
	nil,                   // 6: block.FullBlock.StakeStateChangesEntry
	nil,                   // 7: block.ConfirmBlock.ValidatorSignsEntry
	(*AccountState)(nil),  // 8: account_state.AccountState
	(*Receipt)(nil),       // 9: receipt.Receipt
	(*StakeStates)(nil),   // 10: account_state.StakeStates
}
var file_block_proto_depIdxs = []int32{
	0,  // 0: block.BlockHashData.Type:type_name -> block.BLOCK_TYPE
	0,  // 1: block.Block.Type:type_name -> block.BLOCK_TYPE
	2,  // 2: block.FullBlock.Block:type_name -> block.Block
	8,  // 3: block.FullBlock.AccountStateChanges:type_name -> account_state.AccountState
	9,  // 4: block.FullBlock.Receipts:type_name -> receipt.Receipt
	5,  // 5: block.FullBlock.ValidatorSigns:type_name -> block.FullBlock.ValidatorSignsEntry
	6,  // 6: block.FullBlock.StakeStateChanges:type_name -> block.FullBlock.StakeStateChangesEntry
	7,  // 7: block.ConfirmBlock.ValidatorSigns:type_name -> block.ConfirmBlock.ValidatorSignsEntry
	10, // 8: block.FullBlock.StakeStateChangesEntry.value:type_name -> account_state.StakeStates
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_block_proto_init() }
func file_block_proto_init() {
	if File_block_proto != nil {
		return
	}
	file_account_state_proto_init()
	file_receipt_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHashData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FullBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_block_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_block_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_proto_goTypes,
		DependencyIndexes: file_block_proto_depIdxs,
		EnumInfos:         file_block_proto_enumTypes,
		MessageInfos:      file_block_proto_msgTypes,
	}.Build()
	File_block_proto = out.File
	file_block_proto_rawDesc = nil
	file_block_proto_goTypes = nil
	file_block_proto_depIdxs = nil
}
