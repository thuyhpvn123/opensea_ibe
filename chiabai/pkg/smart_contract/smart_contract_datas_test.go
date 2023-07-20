package smart_contract

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/storage"
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
	rs, _ := scd.GetSmartContractData(testAddress1)
	logger.Info("Test get smart contract data result", rs)
	scd.SetStorages(testAddress1, map[string][]byte{
		"0000000000000000000000000000000000000001": []byte{3, 3, 3},
	})

	rs1, _ := scd.GetSmartContractData(testAddress1)
	logger.Info("Test get smart contract data result 2", rs1)
}
