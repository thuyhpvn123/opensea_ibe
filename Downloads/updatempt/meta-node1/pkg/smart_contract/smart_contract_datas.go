package smart_contract

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
	"github.com/meta-node-blockchain/meta-node/types"
)

type SmartContractDatas struct {
	storage            storage.Storage
	smartContractDatas map[common.Address]types.SmartContractData // this hold live smart contract
	dirty              map[common.Address]types.SmartContractData
}

func NewSmartContractDatas(
	storage storage.Storage,
) types.SmartContractDatas {
	return &SmartContractDatas{
		storage:            storage,
		smartContractDatas: make(map[common.Address]types.SmartContractData),
		dirty:              make(map[common.Address]types.SmartContractData),
	}
}

// getter
func (sd *SmartContractDatas) SmartContractData(address common.Address) (types.SmartContractData, error) {
	if sd.smartContractDatas[address] != nil {
		return sd.smartContractDatas[address], nil
	}
	smartContractData := &SmartContractData{}
	sd.smartContractDatas[address] = smartContractData
	b, err := sd.storage.Get(address.Bytes())
	if err != nil {
		return nil, err
	}
	err = smartContractData.Unmarshal(b)
	if err != nil {
		return nil, err
	}
	return smartContractData, nil
}

func (sd *SmartContractDatas) Storage() storage.Storage {
	return sd.storage
}

func (sd *SmartContractDatas) StorageIterator() storage.IIterator {
	return sd.storage.GetIterator()
}

// setter
func (sd *SmartContractDatas) SetSmartContractData(address common.Address, data types.SmartContractData) {
	sd.smartContractDatas[address] = data
	sd.dirty[address] = data
}

func (sd *SmartContractDatas) SetStorages(address common.Address, storages map[string][]byte) {
	ss, _ := sd.SmartContractData(address)
	currentStorage := ss.Storage()
	if currentStorage == nil {
		currentStorage = make(map[string][]byte)
	}

	// add to dirty storage
	if sd.dirty[address] == nil {
		sd.dirty[address] = NewSmartContractData(nil, make(map[string][]byte))
	}

	for k, v := range storages {
		currentStorage[k] = v
		sd.dirty[address].SetStorage(k, v)
	}
	logger.DebugP("Adding storage change for address", address, sd.dirty[address].Storage())

	ss.SetStorages(currentStorage)
}

// other
func (sd *SmartContractDatas) OpenStorage() error {
	return sd.storage.Open()
}

func (sd *SmartContractDatas) CloseStorage() error {
	return sd.storage.Open()
}

func (sd *SmartContractDatas) Dirty() map[common.Address]types.SmartContractData {
	return sd.dirty
}

func (sd *SmartContractDatas) Commit() error {
	batchData := make([][2][]byte, len(sd.smartContractDatas))
	defer func() {
		sd.smartContractDatas = make(map[common.Address]types.SmartContractData)
		sd.dirty = make(map[common.Address]types.SmartContractData)
	}()
	c := 0
	for i, v := range sd.smartContractDatas {
		bData, err := v.Marshal()
		if err != nil {
			return err
		}
		batchData[c] = [2][]byte{
			i.Bytes(),
			bData,
		}
		c++
	}
	sd.storage.BatchPut(batchData)
	sd.storage.Close()
	return nil
}

func (sd *SmartContractDatas) Cancel() {
	sd.dirty = make(map[common.Address]types.SmartContractData)
}

func (sd *SmartContractDatas) CopyToNewPath(newPath string) (types.SmartContractDatas, error) {
	// delete if exists
	_, err := os.Stat(newPath)
	if err == nil {
		os.RemoveAll(newPath)
	}
	newStorage, err := sd.Storage().CopyToNewPath(newPath)
	if err != nil {
		logger.Error("error when create new storage for native smart contract", err)
		return nil, err
	}
	newSd := NewSmartContractDatas(newStorage)
	return newSd, nil
}
