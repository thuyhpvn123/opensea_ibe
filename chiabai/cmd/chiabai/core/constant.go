package core

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	// . "github.com/ethereum/go-ethereum/accounts/abi"
)

type Account struct {
	Address string
	Private string
}

var PORT int
var Contracts = [...]Contract{
	{Name: "chiabai", Address: "d0b586c617341ddc60d14d04c397a1236fa985fd"},
}

//f8eaba3eb679f6defbe78ce8dd5229ec3622f2a7
func GetPORT() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// log.Info("PORT: ", os.Getenv("PORT"))
	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
	return PORT
}

type Contract struct {
	Name    string
	Address string
}
