package controllers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	// "gitlab.com/meta-node/meta-node/cmd/chiabai/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/cmd/client/command"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"

	. "github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/meta-node/meta-node/cmd/chiabai/network"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

type ActionListenerCallback map[string]interface{}

var defaultRelatedAddress [][]byte
var (
	ErrorGetAccountStateTimedOut = errors.New("get account state timed out")
	ErrorInvalidAction           = errors.New("invalid action")
)

type WalletKey struct {
	PriKey []byte
	PubKey []byte
}

func (caller *CallData) TryCall(callMap map[string]interface{}) interface{} {
	i := 0
	var result interface{}
	result = "TimeOut"

	for {
		if i >= 3 {
			break
		}
		if i != 0 {
			time.Sleep(time.Second)
		}
		result = caller.call(callMap)

		if result != "TimeOut" {
			log.Info("Success time - ", i)
			log.Info(" - Result: ", result)
			return result
		}
		i++
	}

	return result
}

func (caller *CallData) call(callMap map[string]interface{}) interface{} {
	fromAddress, _ := callMap["from"].(string)
	caller.SendTransaction1(callMap)

	for {

		select {
		case receiver := <-caller.client.tcpServerMap[fromAddress].GetHandler():
			// log.Info("Hash on server", common.BytesToHash(hash.([]byte)))
			// log.Info("Hash from chain", (receiver).(network.Receipt).Hash)
			// if (receiver).(network.Receipt).Hash != common.BytesToHash(hash.([]byte)) {
			// 	continue
			// }
			return (receiver).(network.Receipt1).Value
		case <-time.After(5 * time.Second):
			return "TimeOut"
		}
	}

}

var transferMu sync.Mutex

// func (caller *CallData) Transfer(ctx *gin.Context) interface{} {
// 	transferMu.Lock()
// 	defer transferMu.Unlock()
// 	request := make(map[string]interface{})
// 	ctx.BindJSON(request)
// 	fmt.Println("request là:",request)
// 	err:=caller.SendTransaction1(request)
// 		// wait receipt
// 	// select {
// 	// case receipt := <-caller.client.sendChan:
// 	// 	ctx.JSON(http.StatusOK, gin.H{
// 	// 		"transaction_hash": common.Bytes2Hex(receipt.TransactionHash),
// 	// 		"from_address":     common.Bytes2Hex(receipt.FromAddress),
// 	// 		"to_address":       common.Bytes2Hex(receipt.ToAddress),
// 	// 		"amount":           common.Bytes2Hex(receipt.Amount),
// 	// 		"status":           receipt.Status,
// 	// 		"return_value":     common.Bytes2Hex(receipt.ReturnValue),
// 	// 	})
// 	// case <-time.After(5 * time.Second):
// 	// 	ctx.JSON(http.StatusRequestTimeout, gin.H{})
// 	// }

// 	return err
// }
func (caller *CallData) SendTransaction1(call map[string]interface{}) error {
	fmt.Println("call là:", call)
	relatedAddress := caller.EnterRelatedAddress(call)
	fromAddress, _ := call["from"].(string)
	toAddressStr, _ := call["to"].(string)
	toAddress := common.HexToAddress(toAddressStr)
	hexAmount, _ := call["amount"].(string)
	if hexAmount == "" {
		hexAmount = "0"
	}
	amount := uint256.NewInt(0).SetBytes(common.FromHex(hexAmount))
	var maxGas uint64
	maxGaskq, ok := call["gas"].(float64)
	if !ok {
		maxGas = 500000
	}
	maxGas = uint64(maxGaskq)

	var maxGasPriceGwei uint64
	maxGasPriceGweikq, ok := call["gasPrice"].(float64)
	if !ok {
		maxGasPriceGwei = 10
	}
	maxGasPriceGwei = uint64(maxGasPriceGweikq)
	maxGasPrice := 1000000000 * maxGasPriceGwei

	var maxTimeUse uint64
	maxTimeUsekq, ok := call["timeUse"].(float64)
	if !ok {
		maxTimeUse = 60000
	}
	maxTimeUse = uint64(maxTimeUsekq)
	var action pb.ACTION
	action = pb.ACTION_CALL_SMART_CONTRACT

	sign := GetSignGetAccountState(call)

	as, err := caller.GetAccountState(fromAddress, sign)
	if err != nil {
		return err
	}
	data, err := caller.GetDataForCallSmartContract(call)
	if err != nil {
		panic(err)
	}
	transaction, err := caller.client.transactionControllerMap[fromAddress].SendTransaction(
		as.GetLastHash(),
		toAddress,
		as.GetPendingBalance(),
		amount,
		maxGas,
		maxGasPrice,
		maxTimeUse,
		action,
		data,
		relatedAddress,
	)
	logger.Debug("Sending transaction", transaction)
	if err != nil {
		logger.Warn(err)
	}
	fmt.Printf("Send transaction %v", transaction)

	return err
}

func GetSignGetAccountState(call map[string]interface{}) cm.Sign {
	hash := crypto.Keccak256(common.FromHex(call["from"].(string)))

	if call["priKey"] == nil {
		logger.Error(fmt.Sprintf("error when get wallet key "))
	}

	keyPair := bls.NewKeyPair(common.FromHex(call["priKey"].(string)))
	prikey := keyPair.GetPrivateKey()
	sign := bls.Sign(prikey, hash)
	return sign
}

func (caller *CallData) GetAccountState(address string, sign cm.Sign) (state.IAccountState, error) {
	parentConn := caller.client.connectionsManager.GetParentConnection()
	caller.client.messageSenderMap[address].SendBytes(parentConn, command.GetAccountState, common.FromHex(address), sign)

	select {
	case accountState := <-caller.client.accountStateChan:
		return accountState, nil
	case <-time.After(5 * time.Second):
		return nil, ErrorGetAccountStateTimedOut
	}

}

func (caller *CallData) GetWalletInfo(call map[string]interface{}) {

	sign := GetSignGetAccountState(call)
	as, err := caller.GetAccountState(call["from"].(string), sign)
	if err != nil {
		logger.Error(fmt.Sprintf("error when GetAccountState %", err))
		panic(fmt.Sprintf("error when GetAccountState %v", err))
	}
	result := map[string]interface{}{
		"address":         as.GetAddress(),
		"last_hash":       as.GetLastHash(),
		"balance":         as.GetBalance(),
		"pending_balance": as.GetPendingBalance(),
	}
	// header := models.Header{ Success:true,Data: result}
	// kq := utils.NewResultTransformer(header)
	fmt.Println("result:", result)
	// caller.sentToClient("desktop","get-wallet-info", false,kq)
}

func (caller *CallData) sentToClient(command string, data interface{}) {
	caller.client.sendChan <- Message1{command, data}
	// sendQueue[caller.client.ws] <- Message{msgType, value}
}

func (caller *CallData) EnterRelatedAddress(call map[string]interface{}) [][]byte {
	var arrmap []map[string]interface{}
	arr, _ := call["relatedAddresses"].([]interface{})
	if call["relatedAddresses"] == nil || len(arr) == 0 {
		return defaultRelatedAddress
	} else {
		for _, v := range arr {
			arrmap = append(arrmap, v.(map[string]interface{}))
		}

		var relatedAddStr []string

		for _, v := range arrmap {
			relatedAddStr = append(relatedAddStr, v["address"].(string))
		}
		var relatedAddress [][]byte

		// temp := strings.Split(relatedAddStr, ",")
		logger.Info("Temp Related Address")
		for _, addr := range relatedAddStr {
			addressHex := common.HexToAddress(addr)
			logger.Info(addressHex)
			relatedAddress = append(relatedAddress, addressHex.Bytes())
		}
		defaultRelatedAddress = append(defaultRelatedAddress, relatedAddress...)
		return relatedAddress

	}
}
func (caller *CallData) GetDataForCallSmartContract(call map[string]interface{}) ([]byte, error) {
	kq := caller.EncodeAbi(call)
	callData := transaction.NewCallData(kq)
	return callData.Marshal()
}

func (caller *CallData) EncodeAbi(call map[string]interface{}) []byte {
	var inputArray []interface{}
	if call["inputArray"] == nil {
		inputArray = []interface{}{}
	} else {
		inputArray, _ = call["inputArray"].([]interface{})
	}
	functionName, _ := call["function-name"].(string)

	// abiData, ok := call["abiData"].(string)
	// if !ok {
	// 	logger.Error(fmt.Sprintf("error when get abiData %"))
	// 	panic(fmt.Sprintf("error when get abiData "))
	// }
	// abiJson, err := JSON(strings.NewReader(abiData))
	// if err != nil {
	// 	panic(err)
	// }
	reader, err := os.Open("./abi/chiabai.json")
	fmt.Println("1111111111111111111")
	if err != nil {
		log.Fatalf("Error occured while reading %s", "./abi/chiabai.json")
	}
	abiJson, err := JSON(reader)
	if err != nil {
		panic(err)
	}

	var abiTypes []interface{}
	for _, item := range inputArray {
		itemArr := encodeAbiItem(item)
		for _, v := range itemArr {
			abiTypes = append(abiTypes, v)
		}

	}

	fmt.Println("kkkkkkkkkk")
	fmt.Printf("type: %T",abiTypes)
	out, err := abiJson.Pack(functionName, abiTypes[:]...)

	if err != nil {
		panic(err)
	}
	fmt.Println("out:", hex.EncodeToString(out))
	return out
}

func encodeAbiItem(item interface{}) []interface{} {
	var result []interface{}
	var itemMap map[string]interface{}
	fmt.Println("222222222222222")

	if err := json.Unmarshal([]byte(item.(string)), &itemMap); err != nil {
		log.Fatal(err)
	}
	itemType, _ := itemMap["type"].(string)
	fmt.Println("itemType:",itemType)
	switch itemType {
	case "tuple":
		fmt.Println("3333333333")

		var value []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
			log.Fatal(err)
		}

		var components []interface{}
		fmt.Println("444444444444")

		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
			log.Fatal(err)
		}

		var abiTypes []interface{}
		for i, component := range components {
			componentBytes, _ := json.Marshal(component)
			componentType, _ := component.(map[string]interface{})["type"].(string)
			if componentType == "tuple" || componentType == "tuple[]" {
				components[i].(map[string]interface{})["value"] = value[i]
				abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
			} else {
				abiTypes = append(abiTypes, getAbiType(componentType, value[i]))
			}
		}
		result = abiTypes
	case "tuple[]":
		var value []interface{}
		fmt.Println("555555555555")

		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["value"])), &value); err != nil {
			log.Fatal(err)
		}


		fmt.Println("66666666666666")
		var components []interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", itemMap["components"])), &components); err != nil {
			log.Fatal(err)
		}

		var tuples []interface{}
		for _, v := range value {
			vArray := v.([]interface{})
			var abiTypes []interface{}
			for j, component := range components {
				componentBytes, _ := json.Marshal(component)
				componentType, _ := component.(map[string]interface{})["type"].(string)
				components[j].(map[string]interface{})["value"] = vArray[j]
				if componentType == "tuple" || componentType == "tuple[]" {
					abiTypes = append(abiTypes, encodeAbiItem(componentBytes))
				} else {
					abiTypes = append(abiTypes, getAbiType(componentType, vArray[j]))
				}
			}
			tuples = append(tuples, abiTypes...)
		}
		result = tuples
	default:
		fmt.Println("77777777777777")

		value := itemMap["value"]

		var arr []interface{}

		result1 := getAbiType(itemType, value)
		result = append(arr, result1)
		fmt.Println("jjjjjjjjjjjjjj")
	}
	return result
}
func getAbiType(dataType string, data interface{}) interface{} {
	fmt.Println("888888888888")

	if strings.Contains(dataType, "int") {
		params := big.NewInt(0)
		params, ok := params.SetString(fmt.Sprintf("%v", int64(data.(float64))), 10)

		if !ok {
			log.Warn("Format big int: error")
			return nil
		}
		return params

	} else {
		fmt.Println("dataType:",dataType)
		switch dataType {
		
		case "string":
			fmt.Println("aaaaaaaaaaa")
			return data.(string)
		case "bool":
			return data.(bool)
		case "address":
			return common.HexToAddress(data.(string))
		case "uint8":
			intVar, err := strconv.Atoi(data.(string))
			if err != nil {
				log.Warn("Conver Uint8 fail", err)
				return nil
			}
			return uint8(intVar)
		// case "uint", "uint256":
		// 	nubmer := big.NewInt(0)
		// 	nubmer, ok := nubmer.SetString(data.(string), 10)
		// 	if !ok {
		// 		log.Warn("Format big int: error")
		// 		return nil
		// 	}
		// 	return nubmer
		case "array", "slice":
			fmt.Println("999999999")

			fmt.Println("array nè")
			fmt.Println("data:",data)
			rv := reflect.ValueOf(data)
			var out []interface{}
			for i := 0; i < rv.Len(); i++ {
				out = append(out, rv.Index(i).Interface())
			}

			return out
		case "string[]":
			fmt.Println("ppppppp")
			var out []string
			for i := 0; i < len(data.([]interface{})); i++ {
				out = append(out, data.([]interface{})[i].(string))
			}

			return out
		case "address[]":
			fmt.Println("kkkkkk")
			var out []common.Address
			for i := 0; i < len(data.([]interface{})); i++ {
				out = append(out,common.HexToAddress( data.([]interface{})[i].(string)))
			}

			return out

		default:
			fmt.Println("1000000000")
			return data
		}
	}
}
