package apikey_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CreateHandler struct {
	result   *services.ApiKeyDTO
	err      error
	code     int
	ctx      echo.Context
	services services.Registrar
}

func NewCreateHandler(
	ctx echo.Context,
	services services.Registrar,
) *CreateHandler {
	return &CreateHandler{
		ctx:      ctx,
		services: services,
	}
}

func CreateHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewCreateHandler(ctx, services)).Handle().JSON()
}

func (h *CreateHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *CreateHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *CreateHandler) Code() int                            { return h.code }
func (h *CreateHandler) Error() error                         { return h.err }
func (h *CreateHandler) Data() any                            { return h.result }

func (h *CreateHandler) JSON() error {
	if h.err != nil {
		return responses.NewApiKeyResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewApiKeyResponse().Successful(h.ctx, h.result)
}

func (h *CreateHandler) Handle() handlers.IHandler {
	request, err := h.services.ValidatorService.CreateApiKeyRequestValidator(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	jwtToken := h.ctx.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(*services.CustomJwt)
	userID := claims.UserId
	if err = h.services.AuthService.CheckToken(jwtToken.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}

	params := services.CreateApiKeyParams{
		Name:        request.Name,
		NotBefore:   request.NotBefore,
		Expiration:  request.Expiration,
		WriteAccess: request.WriteAccess,
		ReadAccess:  request.ReadAccess,
	}
	h.result, err = h.services.ApiKeyService.Create(userID, params)
	if err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}
