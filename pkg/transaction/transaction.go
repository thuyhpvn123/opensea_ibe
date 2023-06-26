package transaction

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ITransaction interface {
	// general
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
	CalculateHash() common.Hash

	// getter
	GetHash() common.Hash
	GetNewDeviceKey() common.Hash
	GetLastDeviceKey() common.Hash
	GetFromAddress() common.Address
	GetToAddress() common.Address
	GetPubkey() p_common.PublicKey
	GetLastHash() common.Hash
	GetSign() p_common.Sign
	GetAmount() *uint256.Int
	GetPendingUse() *uint256.Int
	GetAction() pb.ACTION
	GetBRelatedAddresses() [][]byte
	GetRelatedAddresses() []common.Address
	GetData() []byte
	GetFee(currentGasPrice uint64) *uint256.Int
	GetProto() protoreflect.ProtoMessage
	GetDeployData() IDeployData
	GetCallData() ICallData
	GetCommissionSign() p_common.Sign
	GetMaxGas() uint64
	GetMaxGasPrice() uint64
	GetMaxTimeUse() uint64
	MaxFee() *uint256.Int
	// setter
	Sign(privateKey p_common.PrivateKey)
	SetHash(common.Hash)

	// changer

	// verifiers
	ValidTransactionHash() bool
	ValidLastHash(fromAccountState state.IAccountState) bool
	ValidDeviceKey(fromAccountState state.IAccountState) bool
	ValidMaxGas() bool
	ValidMaxGasPrice(currentGasPrice uint64) bool
	ValidAmount(fromAccountState state.IAccountState, currentGasPrice uint64) bool
	ValidPendingUse(fromAccountState state.IAccountState) bool
	ValidDeploySmartContractToAccount(fromAccountState state.IAccountState) bool
	ValidCallSmartContractToAccount(toAccountState state.IAccountState) bool
	ValidCallSmartContractData(toAccountState state.IAccountState) bool
	String() string
}

type Transaction struct {
	proto *pb.Transaction
}

func NewTransaction(
	lastHash common.Hash,
	publicKey p_common.PublicKey,
	toAddress common.Address,
	pendingUse *uint256.Int,
	amount *uint256.Int,
	maxGas uint64,
	maxGasPrice uint64,
	maxTimeUse uint64,
	action pb.ACTION,
	data []byte,
	relatedAddresses [][]byte,
	lastDeviceKey common.Hash,
	newDeviceKey common.Hash,
) ITransaction {
	proto := &pb.Transaction{
		LastHash:         lastHash.Bytes(),
		PublicKey:        publicKey.Bytes(),
		ToAddress:        toAddress.Bytes(),
		PendingUse:       pendingUse.Bytes(),
		Amount:           amount.Bytes(),
		MaxGas:           maxGas,
		MaxGasPrice:      maxGasPrice,
		MaxTimeUse:       maxTimeUse,
		Action:           action,
		Data:             data,
		RelatedAddresses: relatedAddresses,
		LastDeviceKey:    lastDeviceKey.Bytes(),
		NewDeviceKey:     newDeviceKey.Bytes(),
	}
	tx := &Transaction{
		proto: proto,
	}
	tx.SetHash(tx.CalculateHash())
	return tx
}

func TransactionsToProto(transactions []ITransaction) []*pb.Transaction {
	rs := make([]*pb.Transaction, len(transactions))
	for i, v := range transactions {
		rs[i] = v.GetProto().(*pb.Transaction)
	}
	return rs
}

func TransactionFromProto(txPb *pb.Transaction) ITransaction {
	return &Transaction{
		proto: txPb,
	}
}

func TransactionsFromProto(pbTxs []*pb.Transaction) []ITransaction {
	rs := make([]ITransaction, len(pbTxs))
	for i, v := range pbTxs {
		rs[i] = TransactionFromProto(v)
	}
	return rs
}

func (t *Transaction) CalculateHash() common.Hash {
	hashPb := &pb.TransactionHashData{
		LastHash:         t.proto.LastHash,
		PublicKey:        t.proto.PublicKey,
		ToAddress:        t.proto.ToAddress,
		PendingUse:       t.proto.PendingUse,
		Amount:           t.proto.Amount,
		MaxGas:           t.proto.MaxGas,
		MaxGasPrice:      t.proto.MaxGasPrice,
		MaxTimeUse:       t.proto.MaxTimeUse,
		Action:           t.proto.Action,
		Data:             t.proto.Data,
		RelatedAddresses: t.proto.RelatedAddresses,
		LastDeviceKey:    t.proto.LastDeviceKey,
		NewDeviceKey:     t.proto.NewDeviceKey,
	}
	bHashPb, _ := proto.Marshal(hashPb)
	return crypto.Keccak256Hash(bHashPb)
}

func (t *Transaction) SetHash(hash common.Hash) {
	t.proto.Hash = hash.Bytes()
}

func (t *Transaction) Sign(privateKey p_common.PrivateKey) {
	t.proto.Sign = bls.Sign(privateKey, t.proto.Hash).Bytes()
}

func (t *Transaction) GetHash() common.Hash {
	return common.BytesToHash(t.proto.Hash)
}

func (t *Transaction) GetLastHash() common.Hash {
	return common.BytesToHash(t.proto.LastHash)
}

func (t *Transaction) Unmarshal(b []byte) error {
	pbTransaction := &pb.Transaction{}
	err := proto.Unmarshal(b, pbTransaction)
	if err != nil {
		return err
	}
	t.proto = pbTransaction
	return nil
}

func (t *Transaction) Marshal() ([]byte, error) {
	return proto.Marshal(t.proto)
}

func (t *Transaction) GetNewDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.NewDeviceKey)
}
func (t *Transaction) GetLastDeviceKey() common.Hash {
	return common.BytesToHash(t.proto.LastDeviceKey)
}
func (t *Transaction) GetFromAddress() common.Address {
	return common.BytesToAddress(
		crypto.Keccak256(t.proto.PublicKey),
	)
}
func (t *Transaction) GetToAddress() common.Address {
	return common.BytesToAddress(t.proto.ToAddress)
}
func (t *Transaction) GetPubkey() p_common.PublicKey {
	return p_common.PubkeyFromBytes(t.proto.PublicKey)
}
func (t *Transaction) GetSign() p_common.Sign {
	return p_common.SignFromBytes(t.proto.Sign)
}
func (t *Transaction) GetAmount() *uint256.Int {
	return uint256.NewInt(0).SetBytes(t.proto.Amount)
}

func (t *Transaction) GetPendingUse() *uint256.Int {
	return uint256.NewInt(0).SetBytes(t.proto.PendingUse)
}

func (t *Transaction) GetAction() pb.ACTION {
	return t.proto.Action
}

func (t *Transaction) GetBRelatedAddresses() [][]byte {
	return t.proto.RelatedAddresses
}

func (t *Transaction) GetRelatedAddresses() []common.Address {
	relatedAddresses := make([]common.Address, len(t.proto.RelatedAddresses)+1)
	for i, v := range t.proto.RelatedAddresses {
		relatedAddresses[i] = common.BytesToAddress(v)
	}
	// append to address
	relatedAddresses[len(t.proto.RelatedAddresses)] = t.GetToAddress()
	return relatedAddresses
}

func (t *Transaction) GetFee(currentGasPrice uint64) *uint256.Int {
	fee := uint256.NewInt(t.proto.MaxGas)
	fee = fee.Mul(fee, uint256.NewInt(currentGasPrice))
	fee = fee.Mul(fee, uint256.NewInt((t.proto.MaxTimeUse/1000)+1.0))
	return fee
}

func (t *Transaction) GetProto() protoreflect.ProtoMessage {
	return t.proto
}

func (t *Transaction) GetData() []byte {
	return t.proto.Data
}

func (t *Transaction) GetDeployData() IDeployData {
	deployData := &DeployData{}
	deployData.Unmarshal(t.GetData())
	return deployData
}

func (t *Transaction) GetCallData() ICallData {
	callData := &CallData{}
	callData.Unmarshal(t.GetData())
	return callData
}

func (t *Transaction) GetCommissionSign() p_common.Sign {
	return p_common.SignFromBytes(t.proto.CommissionSign)
}

func (t *Transaction) GetMaxGas() uint64 {
	return t.proto.MaxGas
}

func (t *Transaction) GetMaxGasPrice() uint64 {
	return t.proto.MaxGasPrice
}

func (t *Transaction) MaxFee() *uint256.Int {
	return uint256.NewInt(0).Mul(
		uint256.NewInt(t.GetMaxGasPrice()),
		uint256.NewInt(t.GetMaxGas()),
	)
}
func (t *Transaction) GetMaxTimeUse() uint64 {
	return t.proto.MaxTimeUse
}

func (t *Transaction) String() string {
	str := fmt.Sprintf(`
	Hash: %v
	From: %v
	To: %v
	Amount: %v
	Action: %v
	Data: %v
	Last Hash: %v
	Max Gas: %v
	Max Gas Price: %v
	Max Time Use: %v
	Sign: %v
`,
		hex.EncodeToString(t.proto.Hash),
		hex.EncodeToString(t.GetFromAddress().Bytes()),
		hex.EncodeToString(t.proto.ToAddress),
		uint256.NewInt(0).SetBytes(t.proto.Amount),
		t.proto.Action,
		hex.EncodeToString(t.proto.Data),
		t.GetLastHash(),
		t.GetMaxGas(),
		t.GetMaxGasPrice(),
		t.GetMaxTimeUse(),
		hex.EncodeToString(t.proto.Sign),
	)
	return str
}

func (t *Transaction) ValidTransactionHash() bool {
	return t.CalculateHash() == t.GetHash()
}

func (t *Transaction) ValidLastHash(fromAccountState state.IAccountState) bool {
	return t.GetLastHash() == fromAccountState.GetLastHash()
}

func (t *Transaction) ValidDeviceKey(fromAccountState state.IAccountState) bool {
	return fromAccountState.GetDeviceKey() == common.Hash{} || // skip check device key if account state doesn't have device key
		crypto.Keccak256Hash(t.GetLastDeviceKey().Bytes()) == fromAccountState.GetDeviceKey()
}

func (t *Transaction) ValidMaxGas() bool {
	return t.GetMaxGas() >= p_common.TRANSFER_GAS_COST
}

func (t *Transaction) ValidMaxGasPrice(currentGasPrice uint64) bool {
	if t.GetToAddress() == p_common.NATIVE_SMART_CONTRACT_REWARD_ADDRESS && t.GetAction() == pb.ACTION_CALL_SMART_CONTRACT {
		// skip check gas price for native smart contract
		return true
	}
	return currentGasPrice <= t.GetMaxGasPrice()
}

func (t *Transaction) ValidAmount(fromAccountState state.IAccountState, currentGasPrice uint64) bool {
	fee := t.GetFee(currentGasPrice)
	if (t.GetCommissionSign() != p_common.Sign{}) {
		fee = uint256.NewInt(0)
	}
	totalBalance := uint256.NewInt(0).Add(fromAccountState.GetBalance(), t.GetPendingUse())
	totalSpend := uint256.NewInt(0).Add(fee, t.GetAmount())
	return !totalBalance.Lt(totalSpend)
}

func (t *Transaction) ValidPendingUse(fromAccountState state.IAccountState) bool {
	pendingBalance := fromAccountState.GetPendingBalance()
	pendingUse := t.GetPendingUse()
	return !pendingUse.Gt(pendingBalance)
}

func (t *Transaction) ValidDeploySmartContractToAccount(fromAccountState state.IAccountState) bool {
	validToAddress := common.BytesToAddress(
		crypto.Keccak256(
			append(
				fromAccountState.GetAddress().Bytes(),
				fromAccountState.GetLastHash().Bytes()...),
		)[12:],
	)
	if validToAddress != t.GetToAddress() {
		logger.Warn("Not match deploy address", validToAddress, t.GetToAddress())
	}
	return validToAddress == t.GetToAddress()
}

func (t *Transaction) ValidCallSmartContractToAccount(toAccountState state.IAccountState) bool {
	// TODO
	return true
}

func (t *Transaction) ValidCallSmartContractData(toAccountState state.IAccountState) bool {
	// TODO
	return true
}
