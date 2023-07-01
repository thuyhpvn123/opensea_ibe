package config

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	p_common "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/config"
)

type Connection struct {
	Address string `json:address`
	Ip      string `json:ip`
	Port    int    `json:port`
	Type    string `json:type`
}

type Config struct {
	Address                    string       `json:"address"`
	ByteAddress                []byte       `json:"-"`
	BytePrivateKey             []byte       `json:"-"`
	Ip                         string       `json:"ip"`
	Port                       int          `json:"port"`
	NodeType                   string       `json:"node_type"`
	HashPerSecond              int          `json:"hash_per_second"`
	TickPerSecond              int          `json:"tick_per_second"`
	TickPerSlot                int          `json:"tick_per_slot"`
	BlockStackSize             int          `json:"block_stack_size"`
	TimeOutTicks               int          `json:"time_out_ticks"` // how many tick validator should wait before create virture block
	TransactionPerHash         int          `json:"transaction_per_hash"`
	NumberOfValidatePohRoutine int          `json:"number_of_validate_poh_routine"`
	AccountDBPath              string       `json:"account_db_path"`
	SecretKey                  string       `json:"secret_key"`
	TransferFee                int          `json:"transfer_fee"`
	GuaranteeAmount            int          `json:"guarantee_amount"`
	TransactionFee             *uint256.Int `json:"-"`
	TransactionFeeHex          string       `json:"transaction_fee"`

	Version          string     `json:"version"`
	BytePublicKey    []byte     `json:"-"`
	ParentConnection Connection `json:"parent_connection"`
	ServerAddress    string     `json:"server_address"`
	PrivateKey string `json:"private_key"`

	TcpIp   string `json:"tcp_ip"`
	TcpPort int    `json:"tcp_port"`

	ParentAddress           string `json:"parent_address"`
	ParentConnectionAddress string `json:"parent_connection_address"`
	ParentConnectionType    string `json:"parent_connection_type"`
	
	ConnectionAddress       string `json:"connection_address"`
	PublicConnectionAddress string `json:"public_connection_address"`

}

const (
	CONFIG_FILE_PATH = "config/conf.json"
)

func LoadConfig(configPath string) (config.IConfig, error) {
	var config Config
	// raw, err := ioutil.ReadFile("config/conf.json")
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, err
	}

	byteAddress, err := hex.DecodeString(config.Address)
	if err != nil {
		return nil, err
	}
	config.ByteAddress = byteAddress
	// log.Printf("Config loaded: %v\n", config)
	// config.BytePrivateKey, config.BytePublicKey, config.ByteAddress = ccrypto.GenerateKeyPairFromSecretKey(config.SecretKey)
	config.TransactionFee = uint256.NewInt(0).SetBytes(common.FromHex(config.TransactionFeeHex))

	return &config, nil
}

// var AppConfig,err = loadConfig(CONFIG_FILE_PATH)

func (config Config) GetVersion() string {
	return config.Version
}

func (config Config) GetPubkey() []byte {
	return config.BytePublicKey
}

func (config Config) GetPrivateKey() []byte {
	return config.BytePrivateKey
}
func (config Config) GetNodeType() string {
	return p_common.CLIENT_CONNECTION_TYPE
}
func (config Config) GetConnectionAddress() string {
	return config.ConnectionAddress
}

func (config Config) GetPublicConnectionAddress() string {
	return config.PublicConnectionAddress
}
