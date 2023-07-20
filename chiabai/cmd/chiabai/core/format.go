package core

import (
	"encoding/hex"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	. "github.com/ethereum/go-ethereum/accounts/abi"
)

type ContractABI struct {
	Name    string
	Address string
	Abi     ABI
}

func (contract *ContractABI) InitContract(info Contract) {
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
func (contract *ContractABI) Decode(name, data string) interface{} {
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
