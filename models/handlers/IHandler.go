package handlers

type IHandler interface {
	Handle() IHandler
	JSON() error
	SetError(err error) IHandler
	SetCode(code int) IHandler
}

func Lock(h IHandler, code int, err error) IHandler {
	h.SetCode(code)
	h.SetError(err)
	return h
}
