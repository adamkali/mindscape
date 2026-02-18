package apikey_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ListHandler struct {
	result   []services.ApiKeyDTO
	err      error
	code     int
	ctx      echo.Context
	services services.Registrar
}

func NewListHandler(
	ctx echo.Context,
	services services.Registrar,
) *ListHandler {
	return &ListHandler{
		ctx:      ctx,
		services: services,
	}
}

func ListHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewListHandler(ctx, services)).Handle().JSON()
}

func (h *ListHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *ListHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *ListHandler) Code() int                            { return h.code }
func (h *ListHandler) Error() error                         { return h.err }
func (h *ListHandler) Data() any                            { return h.result }

func (h *ListHandler) JSON() error {
	if h.err != nil {
		return responses.NewApiKeysResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewApiKeysResponse().Successful(h.ctx, h.result)
}

func (h *ListHandler) Handle() handlers.IHandler {
	jwtToken := h.ctx.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.services.AuthService.CheckToken(jwtToken.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	if h.result, err = h.services.ApiKeyService.List(userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}
