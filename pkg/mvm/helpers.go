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
	"encoding/json"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

func extractExecuteResult(cExecuteResult C.struct_ExecuteResult) *MVMExecuteResult {
	status := pb.RECEIPT_STATUS(cExecuteResult.b_exitReason)
	var exception pb.EXCEPTION
	if status == pb.RECEIPT_STATUS_THREW {
		exception = pb.EXCEPTION(cExecuteResult.b_exception)
	} else {
		exception = pb.EXCEPTION_NONE
	}

	// extract add balance
	mapAddBalance := extractAddBalance(cExecuteResult)
	mapSubBalance := extractSubBalance(cExecuteResult)
	mapCodeChange, mapCodeHash := extractCodeChange(cExecuteResult)
	mapStorageChange, mapStorageRoot := extractStorageChange(cExecuteResult)
	jEventLogs := extractEventLogs(cExecuteResult)

	uptr := unsafe.Pointer(cExecuteResult.b_exmsg)
	exmsg := string(C.GoBytes(uptr, cExecuteResult.length_exmsg))
	C.free(uptr)

	uptr = unsafe.Pointer(cExecuteResult.b_output)
	rt := C.GoBytes(uptr, cExecuteResult.length_output)
	C.free(uptr)

	gasUsed := (uint64)(cExecuteResult.gas_used)

	return &MVMExecuteResult{
		mapAddBalance,
		mapSubBalance,
		mapCodeChange, mapCodeHash,
		mapStorageChange, mapStorageRoot,
		jEventLogs,
		status,
		exception,
		exmsg,
		rt,
		gasUsed,
	}
}

// extract funcs
func extractAddBalance(
	cExecuteResult C.struct_ExecuteResult,
) (
	mapAddBalance map[string][]byte,
) {
	// extract add balance
	bAddBalanceChange := unsafe.Slice(cExecuteResult.b_add_balance_change, cExecuteResult.length_add_balance_change)
	mapAddBalance = make(map[string][]byte, len(bAddBalanceChange))
	for _, v := range bAddBalanceChange {
		uptr := unsafe.Pointer(v)
		addrWithAddBalanceChange := C.GoBytes(uptr, (C.int)(64))
		C.free(uptr)
		mapAddBalance[hex.EncodeToString(addrWithAddBalanceChange[12:32])] = addrWithAddBalanceChange[32:]
	}
	C.free(unsafe.Pointer(cExecuteResult.b_add_balance_change))
	return
}

func extractSubBalance(
	cExecuteResult C.struct_ExecuteResult,
) (
	mapSubBalance map[string][]byte,
) {
	bSubBalanceChange := unsafe.Slice(cExecuteResult.b_sub_balance_change, cExecuteResult.length_sub_balance_change)
	mapSubBalance = make(map[string][]byte, len(bSubBalanceChange))
	for _, v := range bSubBalanceChange {
		uptr := unsafe.Pointer(v)
		addrWithSubBalanceChange := C.GoBytes(uptr, (C.int)(64))
		C.free(uptr)
		mapSubBalance[hex.EncodeToString(addrWithSubBalanceChange[12:32])] = addrWithSubBalanceChange[32:]
	}
	C.free(unsafe.Pointer(cExecuteResult.b_sub_balance_change))
	return
}

func extractCodeChange(
	cExecuteResult C.struct_ExecuteResult,
) (
	mapCodeChange map[string][]byte,
	mapCodeHash map[string][]byte,
) {
	mapCodeChange = make(map[string][]byte, cExecuteResult.length_code_change)
	mapCodeHash = make(map[string][]byte, cExecuteResult.length_code_change)

	bCodeChange := unsafe.Slice(cExecuteResult.b_code_change, cExecuteResult.length_code_change)
	cLengthCodes := unsafe.Slice(cExecuteResult.length_codes, cExecuteResult.length_code_change)
	lengthCodes := make([]int, cExecuteResult.length_code_change)
	for i, v := range cLengthCodes {
		lengthCodes[i] = int(v)
	}

	for i, v := range lengthCodes {
		uptr := unsafe.Pointer(unsafe.Pointer(bCodeChange[i]))
		addrWithCode := C.GoBytes(uptr, (C.int)(v+32))
		C.free(uptr)
		address := hex.EncodeToString(addrWithCode[12:32])
		code := addrWithCode[32:]
		mapCodeChange[address] = code
		mapCodeHash[address] = crypto.Keccak256(code)
	}
	C.free(unsafe.Pointer(cExecuteResult.b_code_change))
	C.free(unsafe.Pointer(cExecuteResult.length_codes))

	return
}

func extractStorageChange(
	cExecuteResult C.struct_ExecuteResult,
) (
	mapStorageChange map[string][][]byte,
	mapStorageRoot map[string][]byte,
) {
	// extract storage changes
	mapStorageChange = make(map[string][][]byte, cExecuteResult.length_storage_change)
	mapStorageRoot = make(map[string][]byte, cExecuteResult.length_storage_change)

	bStorageChange := unsafe.Slice(cExecuteResult.b_storage_change, cExecuteResult.length_storage_change)
	cLengthStorages := unsafe.Slice(cExecuteResult.length_storages, cExecuteResult.length_storage_change)
	lengthStorages := make([]int, cExecuteResult.length_storage_change)
	for i, v := range cLengthStorages {
		lengthStorages[i] = int(v)
	}

	for i, v := range lengthStorages {
		uptr := unsafe.Pointer(unsafe.Pointer(bStorageChange[i]))
		addrWithStorageChanges := C.GoBytes(uptr, (C.int)(v+32))
		C.free(uptr)
		address := hex.EncodeToString(addrWithStorageChanges[12:32])
		addrWithStorageChanges = addrWithStorageChanges[32:]
		storageCount := v / 64
		mapStorageChange[address] = make([][]byte, storageCount)
		for j := 0; j < storageCount; j++ {
			mapStorageChange[address][j] = addrWithStorageChanges[j*64 : (j+1)*64]
		}
	}

	// extract storage root
	bStorageRoots := unsafe.Slice(cExecuteResult.b_storage_roots, cExecuteResult.length_storage_change)
	for _, v := range bStorageRoots {
		uptr := unsafe.Pointer(v)
		addrWithStorageRoot := C.GoBytes(uptr, (C.int)(64))
		C.free(uptr)
		mapStorageRoot[hex.EncodeToString(addrWithStorageRoot[12:32])] = addrWithStorageRoot[32:]
	}

	C.free(unsafe.Pointer(cExecuteResult.b_storage_change))
	C.free(unsafe.Pointer(cExecuteResult.length_storages))
	C.free(unsafe.Pointer(cExecuteResult.b_storage_roots))

	return
}

func extractEventLogs(
	cExecuteResult C.struct_ExecuteResult,
) (
	logJson LogsJson,
) {
	uptr := unsafe.Pointer(unsafe.Pointer(cExecuteResult.b_logs))
	rawLogs := C.GoBytes(uptr, cExecuteResult.length_logs)
	C.free(uptr)
	json.Unmarshal(rawLogs, &logJson.Logs)
	return
}

// accountStates will be modified log hash when call this function if emit logs
func MVMResultToExecuteResult(
	transaction transaction.ITransaction,
	mvmRs *MVMExecuteResult,
	blockNumber *uint256.Int,
	accountStates state.IAccountStates,
) smart_contract.IExecuteResult {
	transactionHash := transaction.GetHash()
	action := transaction.GetAction()
	// if revert then return amount to sender, and sub receiver
	if mvmRs.Status == pb.RECEIPT_STATUS_THREW {
		amount := transaction.GetAmount()
		var mapAddBalance map[string][]byte
		var mapSubBalance map[string][]byte
		if !amount.IsZero() {
			fromAddress := transaction.GetFromAddress().Bytes()
			toAddress := transaction.GetToAddress().Bytes()
			mapAddBalance = map[string][]byte{
				hex.EncodeToString(fromAddress): amount.Bytes(),
			}
			mapSubBalance = map[string][]byte{
				hex.EncodeToString(toAddress): amount.Bytes(),
			}
		}
		return smart_contract.NewExecuteResult(
			transactionHash,
			action,
			mapAddBalance,
			mapSubBalance,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			mvmRs.Status,
			mvmRs.Exception,
			mvmRs.Return,
			mvmRs.GasUsed,
		)
	}

	eventLogs := mvmRs.JEventLogs.GetCompleteEventLogs(
		blockNumber,
		transactionHash,
	)
	mapAddressAddLogs := make(map[common.Address][]smart_contract.IEventLog)
	for _, v := range eventLogs {
		address := v.GetAddress()
		mapAddressAddLogs[address] = append(mapAddressAddLogs[address], v)
	}
	logHash := make(map[string][]byte, len(mapAddressAddLogs))
	for address, logs := range mapAddressAddLogs {
		as, _ := accountStates.GetAccountState(address)
		lastLogsHash := as.GetSmartContractState().GetLogsHash()
		if len(lastLogsHash) == 0 {
			lastLogsHash = common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
		}
		newLogsHash := smart_contract.GetNewLogHash(lastLogsHash, logs)
		accountStates.SetLogsHash(address, newLogsHash)
		logHash[hex.EncodeToString(address.Bytes())] = newLogsHash.Bytes()
	}

	return smart_contract.NewExecuteResult(
		transactionHash,
		action,
		mvmRs.MapAddBalance,
		mvmRs.MapSubBalance,
		mvmRs.MapCodeChange,
		mvmRs.MapCodeHash,
		mvmRs.MapStorageChange,
		mvmRs.MapStorageRoot,
		eventLogs,
		logHash,
		mvmRs.Status,
		mvmRs.Exception,
		mvmRs.Return,
		mvmRs.GasUsed,
	)
}
