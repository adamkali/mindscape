package handlers

import "fmt"

type IHandler interface {
	Handle() IHandler
	JSON() error
	SetError(err error) IHandler
	SetCode(code int) IHandler
	Code() int
	Data() any
	Error() error
}

func Lock(h IHandler, code int, err error) IHandler {
	h.SetCode(code)
	h.SetError(fmt.Errorf("%d Error: %s", code, err.Error()))
	return h
}
