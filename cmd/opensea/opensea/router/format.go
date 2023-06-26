package router

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus" 

	. "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type ContractABI struct {
	Name    string
	Address string
	Abi     ABI
}

func (contract *ContractABI) initContract(info Contract) {
	reader, err := os.Open("./abi/" + info.Name + ".json")
	if err != nil {
		log.Fatalf("Error occured while reading %s", "./abi/"+info.Name+".json")
	}
	contract.Abi, err = JSON(reader)
	if err != nil {
		log.Fatalf("Error occured while init abi %s", info.Name)
	}
	contract.Address = info.Address
	contract.Name = info.Name
	fmt.Println("Init contract ", info.Name)
}

func (contract *ContractABI) decode(name, data string) interface{} {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		log.Fatalf("Error occured while convert data to byte[] - Data: %s", data)
	}
	result := make(map[string]interface{})
	err = contract.Abi.UnpackIntoMap(result, name, bytes)
	if err != nil {
		log.Fatalf("Error occured while unpack %s - %s \n %s \n %s", name, err, data, bytes)
	}
	return result
}

func (contract *ContractABI) encode(name string, args ...interface{}) []byte {
	fmt.Println("222222222222")
	formatedData := contract.formatPreEncode(contract.Abi.Methods[name].Inputs, args)
	fmt.Println("5555")
	data, err := contract.Abi.Pack(name, formatedData[:]...)
	if err != nil {
		log.Fatalf("Error occured while pack %s - %s", name, err)
	}
	return data
}

func (contract *ContractABI) formatPreEncode(args Arguments, data []interface{}) []interface{} {
	fmt.Println("33333333")

	i := 0
	temp := make([]interface{}, len(args))
	for _, arg := range args {
		temp[i] = formatData(arg.Type.String(), data[i])
		i++
	}
	fmt.Println("444444")

	return temp
}

// format utils
func formatData(dataType string, data interface{}) interface{} {
	switch dataType {
	case "string":
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
	case "address[]":
		var addressList []common.Address
		for _, item := range data.([]string) {
			addressList = append(addressList, common.HexToAddress(item))
		}
		return addressList
	case "uint256[]":
		var list []interface{}
		for _, item := range data.([]string) {
			list = append(list, formatData("uint256", item))
		}
		return list
	case "uint", "uint256":
		nubmer := big.NewInt(0)
		nubmer, ok := nubmer.SetString(data.(string), 10)
		if !ok {
			log.Warn("Format big int: error")
			return nil
		}
		return nubmer
	default:
		return nil
	}
}

// func (contract *ContractABI) formatToNumber(data interface{}) string {
// 	number := data.(*big.Int).String()
// 	return number
// }
