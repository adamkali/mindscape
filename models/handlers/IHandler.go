package handlers

type IHandler interface {
	JSON() error
	Handle(fun any) *IHandler
	Lock(code int) *IHandler
}

