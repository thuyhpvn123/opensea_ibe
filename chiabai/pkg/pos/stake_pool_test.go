package pos

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func TestGetStakerConnectionAddressAboveThreshold(t *testing.T) {
	stakePool := NewStakePool()
	stakePool.Stake(testAddress1, uint256.NewInt(10))
	stakePool.Stake(testAddress2, uint256.NewInt(10))
	stakePool.Stake(testAddress3, uint256.NewInt(5))

	stakePool.UpdateConnectionAddress(testAddress1, "127.0.0.1:3001")
	stakePool.UpdateConnectionAddress(testAddress2, "127.0.0.1:3002")
	stakePool.UpdateConnectionAddress(testAddress3, "127.0.0.1:3003")

	rs := stakePool.GetStakerConnectionAddressAboveThreshold(uint256.NewInt(10))
	fmt.Print(rs)
	assert.Equal(t, map[common.Address]string{
		testAddress1: "127.0.0.1:3001",
		testAddress2: "127.0.0.1:3002",
	}, rs)
}
