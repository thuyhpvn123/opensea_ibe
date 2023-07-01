package common

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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
