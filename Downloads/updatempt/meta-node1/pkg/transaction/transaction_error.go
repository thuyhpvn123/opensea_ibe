package transaction

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	pb "github.com/meta-node-blockchain/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
)

var (
	InvalidTransactionHash              = "InvalidTransactionHash"
	InvalidNewDeviceKey                 = "InvalidNewDeviceKey"
	NotMatchLastHash                    = "NotMatchLastHash"
	InvalidLastDeviceKey                = "InvalidLastDeviceKey"
	InvalidAmount                       = "InvalidAmount"
	InvalidPendingUse                   = "InvalidPendingUse"
	InvalidDeploySmartContractToAccount = "InvalidDeploySmartContractToAccount"
	InvalidCallSmartContractToAccount   = "InvalidCallSmartContractToAccount"
	InvalidCallSmartContractData        = "InvalidCallSmartContractData"
	InvalidStakeAddress                 = "InvalidStakeAddress"
	InvalidUnstakeAddress               = "InvalidUnstakeAddress"
	InvalidUnstakeAmount                = "InvalidUnstakeAmount"
	InvalidMaxGas                       = "InvalidMaxGas"
	InvalidMaxGasPrice                  = "InvalidMaxGasPrice"
	InvalidCommissionSign               = "InvalidCommissionSign"
	NotEnoughBalanceForCommissionFee    = "NotEnoughBalanceForCommissionFee"
	InvalidOpenChannelToAccount         = "InvalidOpenChannelToAccount"
	InvalidSign                         = "InvalidSign"
	InvalidCommitAddress                = "InvalidCommitAddress"
	InvalidOpenAccountAmount            = "InvalidOpenAccoutnAmount"

	transactionErrorCodes = map[string]*TransactionErrorCode{
		InvalidTransactionHash:              {1, "invalid transaction hash"},
		InvalidNewDeviceKey:                 {2, "invalid new device key"},
		NotMatchLastHash:                    {3, "not match last hash"},
		InvalidLastDeviceKey:                {4, "invalid last device key"},
		InvalidAmount:                       {5, "invalid amount"},
		InvalidPendingUse:                   {6, "invalid pending use"},
		InvalidDeploySmartContractToAccount: {7, "invalid deploy smart contract to account"},
		InvalidCallSmartContractToAccount:   {8, "invalid call smart contract to account"},
		InvalidCallSmartContractData:        {9, "invalid call smart contract data"},
		InvalidStakeAddress:                 {10, "invalid stake address"},
		InvalidUnstakeAddress:               {11, "invalid unstake address"},
		InvalidUnstakeAmount:                {12, "invalid unstake amount"},
		InvalidMaxGas:                       {13, "invalid max gas"},
		InvalidMaxGasPrice:                  {14, "invalid max gas price"},
		InvalidCommissionSign:               {15, "invalid commission sign"},
		NotEnoughBalanceForCommissionFee:    {16, "smart contract not enough balance for commission fee"},
		InvalidOpenChannelToAccount:         {17, "invalid open channel to account"},
		InvalidSign:                         {18, "invalid sign"},
		InvalidCommitAddress:                {19, "invalid commit address"},
		InvalidOpenAccountAmount:            {20, "invalid open account amount"},
	}
)

type TransactionErrorCode struct {
	code        int64
	description string
}
type TransactionError struct {
	proto *pb.TransactionError
}

func NewTransactionError(
	err string,
	transactionHash common.Hash,
) *TransactionError {
	transactionError := transactionErrorCodes[err]
	return &TransactionError{
		proto: &pb.TransactionError{
			TransactionHash: transactionHash.Bytes(),
			Description:     transactionError.description,
			Code:            transactionError.code,
		},
	}
}

func (transactionErr *TransactionError) Marshal() ([]byte, error) {
	return proto.Marshal(transactionErr.proto)
}

func (transactionErr *TransactionError) Unmarshal(b []byte) error {
	errPb := &pb.TransactionError{}
	err := proto.Unmarshal(b, errPb)
	if err != nil {
		return err
	}
	transactionErr.proto = errPb
	return nil
}

func (transactionErr *TransactionError) String() string {
	str := fmt.Sprintf(
		"Transaction Hash: %v\nCode: %v\nDescription: %v",
		common.BytesToHash(transactionErr.proto.TransactionHash),
		transactionErr.proto.Code,
		transactionErr.proto.Description)
	return str
}
