package state

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "github.com/meta-node-blockchain/meta-node/pkg/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"github.com/meta-node-blockchain/meta-node/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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
	lockingStateChannel []byte,
) types.SmartContractState {
	ssProto := &pb.SmartContractState{
		CreatorPublicKey:    creatorPublicKey,
		StorageHost:         storageHost,
		StorageAddress:      storageAddress,
		CodeHash:            codeHash,
		StorageRoot:         storageRoot,
		LogsHash:            logHash,
		RelatedAddresses:    p_common.AddressesToBytes(relatedAddress),
		LockingStateChannel: lockingStateChannel,
	}
	return SmartContractStateFromProto(ssProto)
}

func SmartContractStateFromProto(ssPb *pb.SmartContractState) types.SmartContractState {
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
func (ss *SmartContractState) Proto() protoreflect.ProtoMessage {
	return ss.proto
}

func (ss *SmartContractState) CreatorPublicKey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(ss.proto.CreatorPublicKey)
}

func (ss *SmartContractState) StorageHost() string {
	return ss.proto.StorageHost
}

func (ss *SmartContractState) StorageAddress() common.Address {
	return common.BytesToAddress(ss.proto.StorageAddress)
}

func (ss *SmartContractState) CodeHash() common.Hash {
	return common.BytesToHash(ss.proto.CodeHash)
}

func (ss *SmartContractState) StorageRoot() common.Hash {
	return common.BytesToHash(ss.proto.StorageRoot)
}

func (ss *SmartContractState) LogsHash() common.Hash {
	return common.BytesToHash(ss.proto.LogsHash)
}

func (ss *SmartContractState) RelatedAddress() []common.Address {
	rs := make([]common.Address, len(ss.proto.RelatedAddresses))
	for i, v := range ss.proto.RelatedAddresses {
		rs[i] = common.BytesToAddress(v)
	}
	return rs
}

func (ss *SmartContractState) LockingStateChannel() common.Address {
	return common.BytesToAddress(ss.proto.LockingStateChannel)
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

func (ss *SmartContractState) SetLockingStateChannel(address common.Address) {
	ss.proto.LockingStateChannel = address.Bytes()
}

type SmartContractStateConfirm struct {
	proto *pb.SmartContractConfirm
}

func NewSmartContractStateConfirm(
	address common.Address,
	smartContractState *pb.SmartContractState,
	blockNumber *uint256.Int,
) types.SmartContractStateConfirm {
	return &SmartContractStateConfirm{
		proto: &pb.SmartContractConfirm{
			Address:            address.Bytes(),
			SmartContractState: smartContractState,
			BlockNumber:        blockNumber.Bytes(),
		},
	}
}

func (ssc *SmartContractStateConfirm) Address() common.Address {
	return common.BytesToAddress(ssc.proto.Address)
}

func (ssc *SmartContractStateConfirm) SmartContractState() types.SmartContractState {
	if ssc.proto.SmartContractState == nil {
		return nil
	}
	return SmartContractStateFromProto(ssc.proto.SmartContractState)
}

func (ssc *SmartContractStateConfirm) BlockNumber() *uint256.Int {
	return uint256.NewInt(0).SetBytes(ssc.proto.BlockNumber)
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
