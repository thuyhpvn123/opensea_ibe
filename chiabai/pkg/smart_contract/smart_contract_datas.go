package smart_contract

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/storage"
)

type ISmartContractDatas interface {
	GetSmartContractData(common.Address) (ISmartContractData, error)
	SetSmartContractData(common.Address, ISmartContractData)

	GetStorageIterator() storage.IIterator
	SetStorages(address common.Address, storages map[string][]byte)
	OpenStorage() error
	CloseStorage() error
	GetStorage() storage.IStorage

	Dirty() map[common.Address]ISmartContractData
	Commit() error
	Cancel()
	CopyToNewPath(newPath string) (ISmartContractDatas, error)
}

type SmartContractDatas struct {
	storage            storage.IStorage
	smartContractDatas map[common.Address]ISmartContractData // this hold live smart contract
	dirty              map[common.Address]ISmartContractData
}

func NewSmartContractDatas(
	storage storage.IStorage,
) ISmartContractDatas {
	return &SmartContractDatas{
		storage:            storage,
		smartContractDatas: make(map[common.Address]ISmartContractData),
		dirty:              make(map[common.Address]ISmartContractData),
	}
}

func (sd *SmartContractDatas) GetSmartContractData(address common.Address) (ISmartContractData, error) {
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

func (sd *SmartContractDatas) SetSmartContractData(address common.Address, data ISmartContractData) {
	sd.smartContractDatas[address] = data
	sd.dirty[address] = data
}

func (sd *SmartContractDatas) GetStorage() storage.IStorage {
	return sd.storage
}

func (sd *SmartContractDatas) OpenStorage() error {
	return sd.storage.Open()
}

func (sd *SmartContractDatas) CloseStorage() error {
	return sd.storage.Open()
}

func (sd *SmartContractDatas) GetStorageIterator() storage.IIterator {
	return sd.storage.GetIterator()
}

func (sd *SmartContractDatas) SetStorages(address common.Address, storages map[string][]byte) {
	ss, _ := sd.GetSmartContractData(address)
	currentStorage := ss.GetStorage()
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
	logger.DebugP("Adding storage change for address", address, sd.dirty[address].GetStorage())

	ss.SetStorages(currentStorage)
}

func (sd *SmartContractDatas) Dirty() map[common.Address]ISmartContractData {
	return sd.dirty
}

func (sd *SmartContractDatas) Commit() error {
	batchData := make([][2][]byte, len(sd.smartContractDatas))
	defer func() {
		sd.smartContractDatas = make(map[common.Address]ISmartContractData)
		sd.dirty = make(map[common.Address]ISmartContractData)
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
	sd.dirty = make(map[common.Address]ISmartContractData)
}

func (sd *SmartContractDatas) CopyToNewPath(newPath string) (ISmartContractDatas, error) {
	// delete if exists
	_, err := os.Stat(newPath)
	if err == nil {
		os.RemoveAll(newPath)
	}
	newStorage, err := sd.GetStorage().CopyToNewPath(newPath)
	if err != nil {
		logger.Error("error when create new storage for native smart contract", err)
		return nil, err
	}
	newSd := NewSmartContractDatas(newStorage)
	return newSd, nil
}
