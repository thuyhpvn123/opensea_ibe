package common

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ErrorInvalidConnectionAddress = errors.New("invalid connection address")
)

func SplitConnectionAddress(address string) (ip string, port int, err error) {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return "", 0, ErrorInvalidConnectionAddress
	}
	intPort, err := strconv.Atoi(splited[1])
	if err != nil {
		return "", 0, err
	}
	return splited[0], intPort, nil
}

func AddressesToBytes(addresses []common.Address) [][]byte {
	rs := make([][]byte, len(addresses))
	for i, v := range addresses {
		rs[i] = v.Bytes()
	}
	return rs
}

func StringToUint256(str string) (*uint256.Int, bool) {
	bigInt := big.NewInt(0)
	_, success := bigInt.SetString(str, 10)
	return uint256.NewInt(0).SetBytes(bigInt.Bytes()), success
}
