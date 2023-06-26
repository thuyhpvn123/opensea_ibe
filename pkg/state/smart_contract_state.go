package state

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ISmartContractState interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string

	// getter
	GetProto() protoreflect.ProtoMessage
	GetCreatorPublicKey() p_common.PublicKey
	GetStorageHost() string
	GetStorageAddress() common.Address
	GetCodeHash() common.Hash
	GetStorageRoot() common.Hash
	GetLogsHash() common.Hash
	GetRelatedAddress() []common.Address

	// setter
	SetCreatorPublicKey(p_common.PublicKey)
	SetStorageHost(string)
	SetCodeHash(common.Hash)
	SetStorageRoot(common.Hash)
	SetLogsHash(common.Hash)
	SetRelatedAddress([]common.Address)
}

type SmartContractState struct {
	proto *pb.SmartContractState
}

func NewSmartContractState(
	creatorPublicKey []byte,
	storageHost string,
	storageAddress []byte,
	codeHash []byte,
	storageRoot []byte,
	logHash []byte,
	relatedAddress []common.Address,
) ISmartContractState {
	ssProto := &pb.SmartContractState{
		CreatorPublicKey: creatorPublicKey,
		StorageHost:      storageHost,
		StorageAddress:   storageAddress,
		CodeHash:         codeHash,
		StorageRoot:      storageRoot,
		LogsHash:         logHash,
		RelatedAddresses: p_common.AddressesToBytes(relatedAddress),
	}
	return SmartContractStateFromProto(ssProto)
}

func SmartContractStateFromProto(ssPb *pb.SmartContractState) ISmartContractState {
	if proto.Equal(ssPb, &pb.SmartContractState{}) {
		return nil
	}
	return &SmartContractState{
		proto: ssPb,
	}
}

// general
func (ss *SmartContractState) Unmarshal(b []byte) error {
	ssProto := &pb.SmartContractState{}
	err := proto.Unmarshal(b, ssProto)
	if err != nil {
		return err
	}
	ss.proto = ssProto
	return nil
}

func (ss *SmartContractState) Marshal() ([]byte, error) {
	return proto.Marshal(ss.proto)
}

func (ss *SmartContractState) String() string {
	str := fmt.Sprintf(`
	CreatorPublicKey: %v
	StorageHost: %v
	CodeHash: %v
	StorageRoot: %v
	LogsHash: %v
	RelatedAddresses: 
`,
		hex.EncodeToString(ss.proto.CreatorPublicKey),
		ss.proto.StorageHost,
		hex.EncodeToString(ss.proto.CodeHash),
		hex.EncodeToString(ss.proto.StorageRoot),
		hex.EncodeToString(ss.proto.LogsHash),
	)
	for _, v := range ss.proto.RelatedAddresses {
		str += fmt.Sprintf("\t%v\n", hex.EncodeToString(v))
	}
	return str
}

// getter
func (ss *SmartContractState) GetProto() protoreflect.ProtoMessage {
	return ss.proto
}

func (ss *SmartContractState) GetCreatorPublicKey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(ss.proto.CreatorPublicKey)
}

func (ss *SmartContractState) GetStorageHost() string {
	return ss.proto.StorageHost
}

func (ss *SmartContractState) GetStorageAddress() common.Address {
	return common.BytesToAddress(ss.proto.StorageAddress)
}

func (ss *SmartContractState) GetCodeHash() common.Hash {
	return common.BytesToHash(ss.proto.CodeHash)
}

func (ss *SmartContractState) GetStorageRoot() common.Hash {
	return common.BytesToHash(ss.proto.StorageRoot)
}

func (ss *SmartContractState) GetLogsHash() common.Hash {
	return common.BytesToHash(ss.proto.LogsHash)
}

func (ss *SmartContractState) GetRelatedAddress() []common.Address {
	rs := make([]common.Address, len(ss.proto.RelatedAddresses))
	for i, v := range ss.proto.RelatedAddresses {
		rs[i] = common.BytesToAddress(v)
	}
	return rs
}

// setter
func (ss *SmartContractState) SetCreatorPublicKey(pk p_common.PublicKey) {
	ss.proto.CreatorPublicKey = pk.Bytes()
}

func (ss *SmartContractState) SetStorageHost(host string) {
	ss.proto.StorageHost = host
}

func (ss *SmartContractState) SetCodeHash(hash common.Hash) {
	ss.proto.CodeHash = hash.Bytes()
}

func (ss *SmartContractState) SetStorageRoot(hash common.Hash) {
	ss.proto.StorageRoot = hash.Bytes()
}

func (ss *SmartContractState) SetLogsHash(hash common.Hash) {
	ss.proto.LogsHash = hash.Bytes()
}

func (ss *SmartContractState) SetRelatedAddress(addresses []common.Address) {
	ss.proto.RelatedAddresses = make([][]byte, len(addresses))
	for i, v := range addresses {
		ss.proto.RelatedAddresses[i] = v.Bytes()
	}
}

type ISmartContractStateConfirm interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	String() string
	GetAddress() common.Address
	GetSmartContractState() ISmartContractState
}

type SmartContractStateConfirm struct {
	proto *pb.SmartContractConfirm
}

func NewSmartContractStateConfirm(address common.Address, smartContractState *pb.SmartContractState) ISmartContractStateConfirm {
	return &SmartContractStateConfirm{
		proto: &pb.SmartContractConfirm{
			Address:            address.Bytes(),
			SmartContractState: smartContractState,
		},
	}
}

func (ssc *SmartContractStateConfirm) GetAddress() common.Address {
	return common.BytesToAddress(ssc.proto.Address)
}

func (ssc *SmartContractStateConfirm) GetSmartContractState() ISmartContractState {
	if ssc.proto.SmartContractState == nil {
		return nil
	}
	return SmartContractStateFromProto(ssc.proto.SmartContractState)
}

// general
func (ssc *SmartContractStateConfirm) Unmarshal(b []byte) error {
	sscProto := &pb.SmartContractConfirm{}
	err := proto.Unmarshal(b, sscProto)
	if err != nil {
		return err
	}
	ssc.proto = sscProto
	return nil
}

func (ssc *SmartContractStateConfirm) Marshal() ([]byte, error) {
	return proto.Marshal(ssc.proto)
}

func (ssc *SmartContractStateConfirm) String() string {
	str := fmt.Sprintf(`
	Address: %v
	CreatorPublicKey: %v
	StorageHost: %v
	CodeHash: %v
	StorageRoot: %v
	LogsHash: %v
	RelatedAddresses: 
`,
		hex.EncodeToString(ssc.proto.Address),
		hex.EncodeToString(ssc.proto.SmartContractState.CreatorPublicKey),
		ssc.proto.SmartContractState.StorageHost,
		hex.EncodeToString(ssc.proto.SmartContractState.CodeHash),
		hex.EncodeToString(ssc.proto.SmartContractState.StorageRoot),
		hex.EncodeToString(ssc.proto.SmartContractState.LogsHash),
	)
	for _, v := range ssc.proto.SmartContractState.RelatedAddresses {
		str += fmt.Sprintf("\t%v\n", hex.EncodeToString(v))
	}
	return str
}
