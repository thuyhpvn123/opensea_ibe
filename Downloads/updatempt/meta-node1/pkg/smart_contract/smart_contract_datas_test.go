package smart_contract

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
	"github.com/meta-node-blockchain/meta-node/pkg/storage"
)

func TestGetSmartContractData(t *testing.T) {
	testAddress1 := common.HexToAddress("0000000000000000000000000000000000000002")
	memorydb := storage.NewMemoryDb()
	testSmartContractData := NewSmartContractData([]byte{1, 2, 3}, map[string][]byte{
		"0000000000000000000000000000000000000001": []byte{1, 1, 1},
		"0000000000000000000000000000000000000002": []byte{2, 2, 2},
	})
	bData, _ := testSmartContractData.Marshal()
	memorydb.Put(testAddress1.Bytes(), bData)
	scd := NewSmartContractDatas(
		memorydb,
	)
	rs, _ := scd.SmartContractData(testAddress1)
	logger.Info("Test get smart contract data result", rs)
	scd.SetStorages(testAddress1, map[string][]byte{
		"0000000000000000000000000000000000000001": []byte{3, 3, 3},
	})

	rs1, _ := scd.SmartContractData(testAddress1)
	logger.Info("Test get smart contract data result 2", rs1)
}

func TestNewSmartContractDatas(t *testing.T) {
	sdLevelDb, err := storage.NewLevelDB(
		// "./test/031b68",
		"./test/01ef",
	)
	if err != nil {
		panic(err)
	}
	// TODO: trie.LoadFromStorage()
	nativeSmartContractDatas := NewSmartContractDatas(sdLevelDb)
	sd, err := nativeSmartContractDatas.SmartContractData(common.HexToAddress("0000000000000000000000000000000000000001"))
	logger.Info(sd, err)
	iter := nativeSmartContractDatas.Storage().GetIterator()
	for iter.Next() {
		// need to copy because iter done save value in next call Next
		cValue := make([]byte, len(iter.Value()))
		cKey := make([]byte, len(iter.Key()))
		copy(cValue, iter.Value())
		copy(cKey, iter.Key())
		logger.Debug("Key: value", cKey, cValue)
	}
}

func TestCommit(t *testing.T) {
	sdLevelDb, err := storage.NewLevelDB(
		"./test/testcommit",
	)
	if err != nil {
		panic(err)
	}
	// TODO: trie.LoadFromStorage()
	nativeSmartContractDatas := NewSmartContractDatas(sdLevelDb)
	sd := NewSmartContractData([]byte{1}, make(map[string][]byte))
	nativeSmartContractDatas.SetSmartContractData(common.HexToAddress("0000000000000000000000000000000000000001"), sd)
	nativeSmartContractDatas.Commit()
	nativeSmartContractDatas.CopyToNewPath("./test/testcommit2")
}
