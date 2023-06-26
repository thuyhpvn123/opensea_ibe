package command

const (
	//General
	InitConnection = "InitConnection"

	// Send messages
	SendTransaction    = "SendTransaction"
	GetAccountState    = "GetAccountState"
	SubscribeToAddress = "SubscribeToAddress"

	// Receive message
	AccountState     = "AccountState"
	Receipt          = "Receipt"
	TransactionError = "TransactionError"
	EventLogs        = "EventLogs"
)
