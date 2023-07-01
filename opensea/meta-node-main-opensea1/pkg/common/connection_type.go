package common

const (
	VALIDATOR_CONNECTION_TYPE     = "validator"
	NODE_CONNECTION_TYPE          = "node"
	CHILD_NODE_CONNECTION_TYPE    = "child_node"
	EXECUTE_MINER_CONNECTION_TYPE = "execute_miner"
	VERIFY_MINER_CONNECTION_TYPE  = "verify_miner"
	CLIENT_CONNECTION_TYPE        = "client"
	STORAGE_CONNECTION_TYPE       = "storage"
	EXPLORER_CONNECTION_TYPE      = "explorer"
)

const (
	NONE_IDX                     = 0
	VALIDATOR_CONNECTION_IDX     = 1
	NODE_CONNECTION_IDX          = 2
	CHILD_NODE_CONNECTION_IDX    = 3
	EXECUTE_MINER_CONNECTION_IDX = 4
	VERIFY_MINER_CONNECTION_IDX  = 5
	CLIENT_CONNECTION_IDX        = 6
	STORAGE_CONNECTION_IDX       = 7
	EXPLORER_CONNECTION_IDX      = 8
)

func MapTypeToIdex(cSType string) int {
	switch cSType {
	case VALIDATOR_CONNECTION_TYPE:
		return VALIDATOR_CONNECTION_IDX
	case NODE_CONNECTION_TYPE:
		return NODE_CONNECTION_IDX
	case CHILD_NODE_CONNECTION_TYPE:
		return CHILD_NODE_CONNECTION_IDX
	case EXECUTE_MINER_CONNECTION_TYPE:
		return EXECUTE_MINER_CONNECTION_IDX
	case VERIFY_MINER_CONNECTION_TYPE:
		return VERIFY_MINER_CONNECTION_IDX
	case CLIENT_CONNECTION_TYPE:
		return CLIENT_CONNECTION_IDX
	case STORAGE_CONNECTION_TYPE:
		return STORAGE_CONNECTION_IDX
	case EXPLORER_CONNECTION_TYPE:
		return EXPLORER_CONNECTION_IDX
	default:
		return NONE_IDX
	}
}
