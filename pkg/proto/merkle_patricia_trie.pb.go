// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: merkle_patricia_trie.proto

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

type MPTNODE_TYPE int32

const (
	MPTNODE_TYPE_FULL  MPTNODE_TYPE = 0
	MPTNODE_TYPE_SHORT MPTNODE_TYPE = 1
	MPTNODE_TYPE_VALUE MPTNODE_TYPE = 2
)

// Enum value maps for MPTNODE_TYPE.
var (
	MPTNODE_TYPE_name = map[int32]string{
		0: "FULL",
		1: "SHORT",
		2: "VALUE",
	}
	MPTNODE_TYPE_value = map[string]int32{
		"FULL":  0,
		"SHORT": 1,
		"VALUE": 2,
	}
)

func (x MPTNODE_TYPE) Enum() *MPTNODE_TYPE {
	p := new(MPTNODE_TYPE)
	*p = x
	return p
}

func (x MPTNODE_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MPTNODE_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_merkle_patricia_trie_proto_enumTypes[0].Descriptor()
}

func (MPTNODE_TYPE) Type() protoreflect.EnumType {
	return &file_merkle_patricia_trie_proto_enumTypes[0]
}

func (x MPTNODE_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MPTNODE_TYPE.Descriptor instead.
func (MPTNODE_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_merkle_patricia_trie_proto_rawDescGZIP(), []int{0}
}

type MPTNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type MPTNODE_TYPE `protobuf:"varint,1,opt,name=type,proto3,enum=merkle_patricia_trie.MPTNODE_TYPE" json:"type,omitempty"`
	Data []byte       `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *MPTNode) Reset() {
	*x = MPTNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_merkle_patricia_trie_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MPTNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MPTNode) ProtoMessage() {}

func (x *MPTNode) ProtoReflect() protoreflect.Message {
	mi := &file_merkle_patricia_trie_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MPTNode.ProtoReflect.Descriptor instead.
func (*MPTNode) Descriptor() ([]byte, []int) {
	return file_merkle_patricia_trie_proto_rawDescGZIP(), []int{0}
}

func (x *MPTNode) GetType() MPTNODE_TYPE {
	if x != nil {
		return x.Type
	}
	return MPTNODE_TYPE_FULL
}

func (x *MPTNode) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type MPTFullNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes [][]byte `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"` // 16 element of 32 bytes hash
	Value []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MPTFullNode) Reset() {
	*x = MPTFullNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_merkle_patricia_trie_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MPTFullNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MPTFullNode) ProtoMessage() {}

func (x *MPTFullNode) ProtoReflect() protoreflect.Message {
	mi := &file_merkle_patricia_trie_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MPTFullNode.ProtoReflect.Descriptor instead.
func (*MPTFullNode) Descriptor() ([]byte, []int) {
	return file_merkle_patricia_trie_proto_rawDescGZIP(), []int{1}
}

func (x *MPTFullNode) GetNodes() [][]byte {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *MPTFullNode) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type MPTShortNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MPTShortNode) Reset() {
	*x = MPTShortNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_merkle_patricia_trie_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MPTShortNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MPTShortNode) ProtoMessage() {}

func (x *MPTShortNode) ProtoReflect() protoreflect.Message {
	mi := &file_merkle_patricia_trie_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MPTShortNode.ProtoReflect.Descriptor instead.
func (*MPTShortNode) Descriptor() ([]byte, []int) {
	return file_merkle_patricia_trie_proto_rawDescGZIP(), []int{2}
}

func (x *MPTShortNode) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *MPTShortNode) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_merkle_patricia_trie_proto protoreflect.FileDescriptor

var file_merkle_patricia_trie_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x6d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x72, 0x69, 0x63, 0x69,
	0x61, 0x5f, 0x74, 0x72, 0x69, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6d, 0x65,
	0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x72, 0x69, 0x63, 0x69, 0x61, 0x5f, 0x74, 0x72,
	0x69, 0x65, 0x22, 0x55, 0x0a, 0x07, 0x4d, 0x50, 0x54, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x36, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x6d, 0x65,
	0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x72, 0x69, 0x63, 0x69, 0x61, 0x5f, 0x74, 0x72,
	0x69, 0x65, 0x2e, 0x4d, 0x50, 0x54, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x39, 0x0a, 0x0b, 0x4d, 0x50, 0x54,
	0x46, 0x75, 0x6c, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x36, 0x0a, 0x0c, 0x4d, 0x50, 0x54, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x2e, 0x0a, 0x0c,
	0x4d, 0x50, 0x54, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x08, 0x0a, 0x04,
	0x46, 0x55, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x48, 0x4f, 0x52, 0x54, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x02, 0x42, 0x3c, 0x0a, 0x32,
	0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x64, 0x2e, 0x6d, 0x65,
	0x72, 0x6b, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x72, 0x69, 0x63, 0x69, 0x61, 0x5f, 0x74, 0x72,
	0x69, 0x65, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_merkle_patricia_trie_proto_rawDescOnce sync.Once
	file_merkle_patricia_trie_proto_rawDescData = file_merkle_patricia_trie_proto_rawDesc
)

func file_merkle_patricia_trie_proto_rawDescGZIP() []byte {
	file_merkle_patricia_trie_proto_rawDescOnce.Do(func() {
		file_merkle_patricia_trie_proto_rawDescData = protoimpl.X.CompressGZIP(file_merkle_patricia_trie_proto_rawDescData)
	})
	return file_merkle_patricia_trie_proto_rawDescData
}

var file_merkle_patricia_trie_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_merkle_patricia_trie_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_merkle_patricia_trie_proto_goTypes = []interface{}{
	(MPTNODE_TYPE)(0),    // 0: merkle_patricia_trie.MPTNODE_TYPE
	(*MPTNode)(nil),      // 1: merkle_patricia_trie.MPTNode
	(*MPTFullNode)(nil),  // 2: merkle_patricia_trie.MPTFullNode
	(*MPTShortNode)(nil), // 3: merkle_patricia_trie.MPTShortNode
}
var file_merkle_patricia_trie_proto_depIdxs = []int32{
	0, // 0: merkle_patricia_trie.MPTNode.type:type_name -> merkle_patricia_trie.MPTNODE_TYPE
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_merkle_patricia_trie_proto_init() }
func file_merkle_patricia_trie_proto_init() {
	if File_merkle_patricia_trie_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_merkle_patricia_trie_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MPTNode); i {
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
		file_merkle_patricia_trie_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MPTFullNode); i {
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
		file_merkle_patricia_trie_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MPTShortNode); i {
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
			RawDescriptor: file_merkle_patricia_trie_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_merkle_patricia_trie_proto_goTypes,
		DependencyIndexes: file_merkle_patricia_trie_proto_depIdxs,
		EnumInfos:         file_merkle_patricia_trie_proto_enumTypes,
		MessageInfos:      file_merkle_patricia_trie_proto_msgTypes,
	}.Build()
	File_merkle_patricia_trie_proto = out.File
	file_merkle_patricia_trie_proto_rawDesc = nil
	file_merkle_patricia_trie_proto_goTypes = nil
	file_merkle_patricia_trie_proto_depIdxs = nil
}
