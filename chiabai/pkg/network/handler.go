package network

type IHandler interface {
	HandleRequest(IRequest) error
	GetChData()chan interface{}
}
