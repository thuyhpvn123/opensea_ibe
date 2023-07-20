package mvm

/*
#cgo CFLAGS: -w
#cgo CXXFLAGS: -std=c++17 -w
#cgo LDFLAGS: -L./linker/build/lib/static -lmvm_linker -L./c_mvm/build/lib/static -lmvm -lstdc++
#cgo CPPFLAGS: -I./linker/build/include
#include "mvm_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

var apiInstance *MVMApi

type IAccountStateDB interface {
	AddPendingBalance(address common.Address, amount *uint256.Int)
	SubTotalBalance(address common.Address, amount *uint256.Int) error

	GetAccountState(address common.Address) (state.IAccountState, error)
	SetAccountState(state.IAccountState) error
	SetSmartContractState(address common.Address, scState state.ISmartContractState)
	SetStorageRoot(address common.Address, hash common.Hash)
}

type ISmartContractDB interface {
	GetSmartContractData(common.Address) (smart_contract.ISmartContractData, error)
	SetSmartContractData(common.Address, smart_contract.ISmartContractData)
	SetStorages(address common.Address, storages map[string][]byte)
}

type MVMApi struct {
	smartContractDb         ISmartContractDB
	accountStateDb          IAccountStateDB
	currentRelatedAddresses map[common.Address]struct{}
}

func GetMVMApiInstance() *MVMApi {
	return apiInstance
}

func InitMVMApi(
	smartContractDb ISmartContractDB,
	accountStateDb IAccountStateDB,
) {
	if apiInstance == nil {
		apiInstance = &MVMApi{
			smartContractDb,
			accountStateDb,
			make(map[common.Address]struct{}),
		}
	}
}

func (a *MVMApi) SetSmartContractDatas(smartContractDb ISmartContractDB) {
	a.smartContractDb = smartContractDb
}

func (a *MVMApi) GetSmartContractDatas() ISmartContractDB {
	return a.smartContractDb
}

func (a *MVMApi) SetAccountStateDb(accountStateDb IAccountStateDB) {
	a.accountStateDb = accountStateDb
}

func (a *MVMApi) GetAccountStateDb() IAccountStateDB {
	return a.accountStateDb
}

func (a *MVMApi) SetRelatedAddresses(addresses []common.Address) {
	a.currentRelatedAddresses = make(map[common.Address]struct{}, len(addresses))
	for _, v := range addresses {
		a.currentRelatedAddresses[v] = struct{}{}
	}
}

func (a *MVMApi) InRelatedAddress(address common.Address) bool {
	_, ok := a.currentRelatedAddresses[address]
	return ok
}

func (a *MVMApi) Call(
	// transaction data
	bSender []byte,
	bContractAddress []byte,
	bInput []byte,
	amount *uint256.Int,
	gasPrice uint64,
	gasLimit uint64,
	// block context data
	blockPrevrandao uint64,
	blockGasLimit uint64,
	blockTime uint64,
	blockBaseFee uint64,
	blockNumber *uint256.Int,
	blockCoinbase common.Address,
) *MVMExecuteResult {
	// transaction data
	bAmount := amount.Bytes32()
	cBSender := C.CBytes(bSender)
	cBContractAddress := C.CBytes(bContractAddress)
	cBInput := C.CBytes(bInput)
	cBAmount := C.CBytes(bAmount[:])

	// block context data
	bBlockNumber := blockNumber.Bytes32()
	bBlockCoinbase := blockNumber.Bytes()
	cBBlockNumber := C.CBytes(bBlockNumber[:])
	cBBlockCoinbase := C.CBytes(bBlockCoinbase)

	defer C.free(unsafe.Pointer(cBSender))
	defer C.free(unsafe.Pointer(cBContractAddress))
	defer C.free(unsafe.Pointer(cBInput))
	defer C.free(unsafe.Pointer(cBAmount))

	defer C.free(unsafe.Pointer(cBBlockNumber))
	defer C.free(unsafe.Pointer(cBBlockCoinbase))

	cRs := C.call(
		// transaction data
		(*C.uchar)(cBSender),
		(*C.uchar)(cBContractAddress),
		(*C.uchar)(cBInput),
		(C.int)(len(bInput)),
		(*C.uchar)(cBAmount),
		(C.ulonglong)(gasPrice),
		(C.ulonglong)(gasLimit),
		// block context data
		(C.ulonglong)(blockPrevrandao),
		(C.ulonglong)(blockGasLimit),
		(C.ulonglong)(blockTime),
		(C.ulonglong)(blockBaseFee),
		(*C.uchar)(cBBlockNumber),
		(*C.uchar)(cBBlockCoinbase),
	)

	return extractExecuteResult(cRs)
}

func (a *MVMApi) Deploy(
	// transaction data
	bSender []byte,
	bLastHash []byte,
	bContractConstructor []byte,
	amount *uint256.Int,
	gasPrice uint64,
	gasLimit uint64,
	// block context data
	blockPrevrandao uint64,
	blockGasLimit uint64,
	blockTime uint64,
	blockBaseFee uint64,
	blockNumber *uint256.Int,
	blockCoinbase common.Address,
) *MVMExecuteResult {
	// transaction data
	bAmount := amount.Bytes32()
	constructorLength := len(bContractConstructor)
	cBSender := C.CBytes(bSender)
	cBLastHash := C.CBytes(bLastHash)
	cBContractConstructor := C.CBytes(bContractConstructor)
	cBAmount := C.CBytes(bAmount[:])
	// block context data
	bBlockNumber := blockNumber.Bytes32()
	bBlockCoinbase := blockCoinbase.Bytes()

	cBBlockNumber := C.CBytes(bBlockNumber[:])
	cBBlockCoinbase := C.CBytes(bBlockCoinbase)

	defer C.free(unsafe.Pointer(cBSender))
	defer C.free(unsafe.Pointer(cBLastHash))
	defer C.free(unsafe.Pointer(cBContractConstructor))
	defer C.free(unsafe.Pointer(cBAmount))

	defer C.free(unsafe.Pointer(cBBlockNumber))
	defer C.free(unsafe.Pointer(cBBlockCoinbase))

	cRs := C.deploy(
		// transaction data
		(*C.uchar)(cBSender),
		(*C.uchar)(cBLastHash),
		(*C.uchar)(cBContractConstructor),
		(C.int)(constructorLength),
		(*C.uchar)(cBAmount),
		(C.ulonglong)(gasPrice),
		(C.ulonglong)(gasLimit),
		// block context data
		(C.ulonglong)(blockPrevrandao),
		(C.ulonglong)(blockGasLimit),
		(C.ulonglong)(blockTime),
		(C.ulonglong)(blockBaseFee),
		(*C.uchar)(cBBlockNumber),
		(*C.uchar)(cBBlockCoinbase),
	)
	return extractExecuteResult(cRs)
}

func (a *MVMApi) UpdateState(
	transaction transaction.ITransaction,
	mvmRs *MVMExecuteResult,
) {
	// if revert then return amount to sender, and sub receiver
	if mvmRs.Status == pb.RECEIPT_STATUS_THREW {
		a.accountStateDb.AddPendingBalance(
			transaction.GetFromAddress(),
			transaction.GetAmount(),
		)
		a.accountStateDb.SubTotalBalance(
			transaction.GetToAddress(),
			transaction.GetAmount(),
		)
		return
	}

	// update add balance
	for address, addAmount := range mvmRs.MapAddBalance {
		fmtAddress := common.HexToAddress(address)
		a.accountStateDb.AddPendingBalance(
			fmtAddress,
			uint256.NewInt(0).SetBytes(addAmount),
		)
	}
	// update sub balance
	// when execute it's only can sub balance of smart contract, so use sub total balance instead of sub balance
	for address, subAmount := range mvmRs.MapSubBalance {
		fmtAddress := common.HexToAddress(address)
		a.accountStateDb.SubTotalBalance(
			fmtAddress,
			uint256.NewInt(0).SetBytes(subAmount),
		)
	}

	// update deploy contract
	if len(mvmRs.MapCodeHash) > 0 {
		var creatorPublicKey p_common.PublicKey
		var storageHost string
		var storageAddress common.Address
		if transaction.GetAction() == pb.ACTION_DEPLOY_SMART_CONTRACT {
			creatorPublicKey = transaction.GetPubkey()
			storageHost = transaction.GetDeployData().GetStorageHost()
			storageAddress = transaction.GetDeployData().GetStorageAddress()
		} else {
			originSmartContractAs, _ := a.accountStateDb.GetAccountState(
				transaction.GetToAddress(),
			)
			creatorPublicKey = originSmartContractAs.GetSmartContractState().GetCreatorPublicKey()
			storageHost = originSmartContractAs.GetSmartContractState().GetStorageHost()
			storageAddress = originSmartContractAs.GetSmartContractState().GetStorageAddress()
		}
		for address, newCodeHash := range mvmRs.MapCodeHash {
			// create new account
			fmtAddress := common.HexToAddress(address)
			asState := state.NewAccountState(fmtAddress)
			a.accountStateDb.SetAccountState(asState)
			// set smart contract state
			scState := state.NewSmartContractState(
				creatorPublicKey.Bytes(),
				storageHost,
				storageAddress.Bytes(),
				newCodeHash,
				nil, nil, nil,
			)
			a.accountStateDb.SetSmartContractState(
				fmtAddress,
				scState,
			)
		}
	}

	// update storage root
	for address, newStorageRoot := range mvmRs.MapStorageRoot {
		fmtAddress := common.HexToAddress(address)
		a.accountStateDb.SetStorageRoot(
			fmtAddress,
			common.BytesToHash(newStorageRoot),
		)
	}

	// update code
	for address, code := range mvmRs.MapCodeChange {
		fmtAddress := common.HexToAddress(address)
		// create smart contract data
		smartContractData := smart_contract.NewSmartContractData(
			code, nil,
		)
		a.smartContractDb.SetSmartContractData(fmtAddress, smartContractData)
	}

	// update storage
	for address, rawStorages := range mvmRs.MapStorageChange {
		storages := make(map[string][]byte, len(rawStorages))
		for _, kv := range rawStorages {
			storages[hex.EncodeToString(kv[:32])] = kv[32:]
		}
		fmtAddress := common.HexToAddress(address)
		a.smartContractDb.SetStorages(fmtAddress, storages)
	}
}

// GLOBAL STATE Functions

//export GlobalStateGet
func GlobalStateGet(
	address *C.uchar,
) (
	status C.int, // 0 not found, 1 found, 2 not in related
	balance_p *C.uchar,
	code_p *C.uchar,
	code_length C.int,
	storage_p *C.uchar,
	storage_length C.int,
) {

	mvmApi := GetMVMApiInstance()
	bAddress := C.GoBytes(unsafe.Pointer(address), 20)
	fAddress := common.BytesToAddress(bAddress)
	inRelatedAddresses := mvmApi.InRelatedAddress(fAddress)
	if !inRelatedAddresses {
		return C.int(2), (*C.uchar)(C.CBytes([]byte{})), (*C.uchar)(C.CBytes([]byte{})), 0, nil, 0
	}

	accountState, err := mvmApi.accountStateDb.GetAccountState(fAddress)
	if err != nil {
		panic(err)
	}
	if accountState == nil {
		return C.int(0), (*C.uchar)(C.CBytes([]byte{})), (*C.uchar)(C.CBytes([]byte{})), 0, nil, 0
	}

	u256Balance := uint256.NewInt(0).Add(
		accountState.GetBalance(),
		accountState.GetPendingBalance(),
	)
	b32Balance := u256Balance.Bytes32()
	bCode := []byte{}
	gStorage := []byte{}
	lenStorage := 0

	smartContractState := accountState.GetSmartContractState()
	if smartContractState != nil {
		scData, err := mvmApi.smartContractDb.GetSmartContractData(fAddress)
		if err != nil {
			panic(err)
		}
		bCode = scData.GetCode()
		for k, v := range scData.GetStorage() {
			b := common.FromHex(k)
			b = append(b, v...)
			gStorage = append(gStorage, b...)
			lenStorage++
		}
	}
	cStorage := C.CBytes(gStorage[:])
	lenCode := len(bCode)

	cBBalance := C.CBytes(b32Balance[:])
	cBCode := C.CBytes(bCode)

	return C.int(1), (*C.uchar)(cBBalance), (*C.uchar)(cBCode), (C.int)(lenCode), (*C.uchar)(cStorage), (C.int)(lenStorage)
}

// go functions
