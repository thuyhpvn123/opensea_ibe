package c_version

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetStorageRoot(t *testing.T) {
	v := common.FromHex("1111111111111111111111111111111111111111111111111111111111111111")
	testStorage := map[string][]byte{
		"1111111111111111111111111111111111111111111111111111111111111111": []byte(v),
		"2222222222222222222222222222222222222222222222222222222222222222": []byte(v),
	}
	root := GetStorageRoot(testStorage)
	fmt.Print(root)
	assert.NotNil(t, root)
}
