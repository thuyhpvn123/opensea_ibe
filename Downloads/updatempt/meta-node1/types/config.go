package types

type Config interface {
	Version() string
	NodeType() string
	PrivateKey() []byte
	PublicConnectionAddress() string
	ConnectionAddress() string
}
