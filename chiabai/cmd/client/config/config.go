package config

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/config"
)

type ClientConfig struct {
	PrivateKey string `json:"private_key"`

	ConnectionAddress       string `json:"connection_address"`
	PublicConnectionAddress string `json:"public_connection_address"`

	Version           string       `json:"version"`
	TransactionFeeHex string       `json:"transaction_fee"`
	TransactionFee    *uint256.Int `json:"-"`

	ParentAddress           string `json:"parent_address"`
	ParentConnectionAddress string `json:"parent_connection_address"`
	ParentConnectionType    string `json:"parent_connection_type"`
}

func (c *ClientConfig) GetConnectionAddress() string {
	return c.ConnectionAddress
}

func (c *ClientConfig) GetPublicConnectionAddress() string {
	return c.PublicConnectionAddress
}

func (c *ClientConfig) GetVersion() string {
	return c.Version
}

func (c *ClientConfig) GetPrivateKey() []byte {
	return common.FromHex(c.PrivateKey)
}

func (c *ClientConfig) GetNodeType() string {
	return p_common.CLIENT_CONNECTION_TYPE
}

func LoadConfig(configPath string) (config.IConfig, error) {
	// general config
	config := &ClientConfig{}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}
	config.TransactionFee = uint256.NewInt(0).SetBytes(common.FromHex(config.TransactionFeeHex))
	return config, nil
}
