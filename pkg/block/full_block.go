package block

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/receipt"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"google.golang.org/protobuf/proto"
)

type IFullBlock interface {
	Unmarshal(b []byte) error
	LoadFromProto(fbProto *pb.FullBlock)
	GetBlock() IBlock
	SetBlock(IBlock)
	GetValidatorSigns() map[cm.PublicKey]cm.Sign
	AddValidatorSign(cm.PublicKey, cm.Sign)
	SetValidatorSigns(map[cm.PublicKey]cm.Sign)
	GetAccountStateChanges() map[common.Address]state.IAccountState
	GetReceipts() receipt.IReceipts
	GetGasUsed() uint64
	GetFee() *uint256.Int
	Marshal() ([]byte, error)
	SetAccountStates(state.IAccountStates)
	GetAccountStates() state.IAccountStates
	SetNativeSmartContractDatas(smart_contract.ISmartContractDatas)
	GetNativeSmartContractDatas() smart_contract.ISmartContractDatas
	NativeSmartContractDataChanges() map[common.Address]smart_contract.ISmartContractData
	GetNextLeaderAddress() common.Address

	SetTimeStamp(uint64)
	SetNextLeaderAddress(common.Address)
}

type FullBlock struct {
	block    IBlock
	receipts receipt.IReceipts
	// account state data
	accountStateChanges map[common.Address]state.IAccountState
	accountStates       state.IAccountStates

	nativeSmartContractDataChanges map[common.Address]smart_contract.ISmartContractData
	nativeSmartContractDatas       smart_contract.ISmartContractDatas

	//
	validatorSigns map[cm.PublicKey]cm.Sign

	nextLeaderAddress common.Address
}

func NewFullBlock(
	block IBlock,
	receipts receipt.IReceipts,
	validatorSigns map[cm.PublicKey]cm.Sign,
	accountStates state.IAccountStates,
	accountStateChanges map[common.Address]state.IAccountState,
	nativeSmartContractDatas smart_contract.ISmartContractDatas,
	nativeSmartContractDataChanges map[common.Address]smart_contract.ISmartContractData,
	nextLeaderAddress common.Address,
) IFullBlock {
	return &FullBlock{
		block:                          block,
		receipts:                       receipts,
		accountStates:                  accountStates,
		accountStateChanges:            accountStateChanges,
		nativeSmartContractDatas:       nativeSmartContractDatas,
		nativeSmartContractDataChanges: nativeSmartContractDataChanges,
		nextLeaderAddress:              nextLeaderAddress,
	}
}

func (fb *FullBlock) Unmarshal(b []byte) error {
	fbProto := &pb.FullBlock{}
	err := proto.Unmarshal(b, fbProto)
	if err != nil {
		return err
	}
	fb.LoadFromProto(fbProto)
	return nil
}

func (fb *FullBlock) LoadFromProto(fbProto *pb.FullBlock) {
	block := NewBlock(fbProto.Block)

	receipts := receipt.NewReceipts()
	for _, v := range fbProto.Receipts {
		receipts.AddReceipt(receipt.ReceiptFromProto(v))
	}

	accountStateChanges := map[common.Address]state.IAccountState{}
	for _, v := range fbProto.AccountStateChanges {
		accountStateChanges[common.BytesToAddress(v.Address)] = state.AccountStateFromProto(v)
	}

	validatorSigns := make(map[cm.PublicKey]cm.Sign)
	for i, v := range fbProto.ValidatorSigns {
		validatorSigns[cm.PubkeyFromBytes(common.FromHex(i))] = cm.SignFromBytes(v)
	}

	// native smart contract data:
	smartContractDataChanges := make(map[common.Address]smart_contract.ISmartContractData, len(fbProto.NativeSmartContractDataChanges))
	for a, v := range fbProto.NativeSmartContractDataChanges {
		smartContractDataChanges[common.HexToAddress(a)] = smart_contract.SmartContractDataFromProto(v)
	}

	*fb = FullBlock{
		block:                          block,
		receipts:                       receipts,
		accountStateChanges:            accountStateChanges,
		validatorSigns:                 validatorSigns,
		nativeSmartContractDataChanges: smartContractDataChanges,
	}
}

func (fb *FullBlock) GetBlock() IBlock {
	return fb.block
}

func (fb *FullBlock) SetBlock(b IBlock) {
	fb.block = b
}

func (fb *FullBlock) GetValidatorSigns() map[cm.PublicKey]cm.Sign {
	return fb.validatorSigns
}

func (fb *FullBlock) AddValidatorSign(pk cm.PublicKey, sign cm.Sign) {
	fb.validatorSigns[pk] = sign
}

func (fb *FullBlock) SetValidatorSigns(signs map[cm.PublicKey]cm.Sign) {
	fb.validatorSigns = signs
}

func (fb *FullBlock) GetAccountStateChanges() map[common.Address]state.IAccountState {
	return fb.accountStateChanges
}

func (fb *FullBlock) GetReceipts() receipt.IReceipts {
	return fb.receipts
}

func (fb *FullBlock) SetAccountStates(as state.IAccountStates) {
	fb.accountStates = as.Copy()
}

func (fb *FullBlock) GetAccountStates() state.IAccountStates {
	return fb.accountStates
}

func (fb *FullBlock) SetNativeSmartContractDatas(sd smart_contract.ISmartContractDatas) {
	fb.nativeSmartContractDatas = sd
}

func (fb *FullBlock) GetNativeSmartContractDatas() smart_contract.ISmartContractDatas {
	return fb.nativeSmartContractDatas
}

func (fb *FullBlock) NativeSmartContractDataChanges() map[common.Address]smart_contract.ISmartContractData {
	return fb.nativeSmartContractDataChanges
}
func (fb *FullBlock) Marshal() ([]byte, error) {
	validatorSigns := make(map[string][]byte, len(fb.validatorSigns))
	for i, v := range fb.validatorSigns {
		validatorSigns[hex.EncodeToString(i.Bytes())] = v.Bytes()
	}
	var receipts []*pb.Receipt
	if fb.receipts != nil {
		receiptsMap := fb.receipts.GetReceiptsMap()
		receipts = make([]*pb.Receipt, len(receiptsMap))
		i := 0
		for _, v := range receiptsMap {
			receipts[i] = v.GetProto().(*pb.Receipt)
			i++
		}
	}

	accountStateChanges := make([]*pb.AccountState, len(fb.accountStateChanges))
	i := 0
	for _, v := range fb.accountStateChanges {
		accountStateChanges[i] = v.GetProto().(*pb.AccountState)
		i++
	}

	// native smart contract data
	// storage
	var dirtyNativeSC map[string]*pb.SmartContractData
	if fb.nativeSmartContractDatas != nil {
		dirtyNativeSC = make(map[string]*pb.SmartContractData, len(fb.nativeSmartContractDatas.Dirty()))
		for i, v := range fb.nativeSmartContractDatas.Dirty() {
			logger.DebugP("Enter dirty native smart contract storage", i, v)
			dirtyNativeSC[hex.EncodeToString(i.Bytes())] = v.GetProto().(*pb.SmartContractData)
		}
	}

	fbProto := &pb.FullBlock{
		Block:                          fb.block.GetProto().(*pb.Block),
		Receipts:                       receipts,
		AccountStateChanges:            accountStateChanges,
		ValidatorSigns:                 validatorSigns,
		NativeSmartContractDataChanges: dirtyNativeSC,
	}
	return proto.Marshal(fbProto)
}

func (fb *FullBlock) SetTimeStamp(timestamp uint64) {
	fb.block.SetTimeStamp(timestamp)
}

func (fb *FullBlock) SetNextLeaderAddress(address common.Address) {
	fb.nextLeaderAddress = address
}

func (fb *FullBlock) GetNextLeaderAddress() common.Address {
	return fb.nextLeaderAddress
}

func (fb *FullBlock) GetGasUsed() uint64 {
	gasUsed := uint64(0)
	if fb.receipts == nil {
		return gasUsed
	} else {
		for _, v := range fb.receipts.GetReceiptsMap() {
			gasUsed += v.GetGasUsed()
		}
	}
	return gasUsed
}

func (fb *FullBlock) GetFee() *uint256.Int {
	return uint256.NewInt(0).Mul(
		uint256.NewInt(fb.GetBlock().GetBaseFee()),
		uint256.NewInt(fb.GetGasUsed()),
	)
}
