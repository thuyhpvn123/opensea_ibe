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
	{Name: "marketplace", Address: "227e5cba2e6f953512a58649cd661643442bf096"},
}

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
