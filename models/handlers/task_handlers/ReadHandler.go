package task_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ReadHandler struct {
	result   []services.TaskDTO
	err      error
	code     int
	ctx      echo.Context
	services services.Registrar
}

func NewReadHandler(
	ctx echo.Context,
	services services.Registrar,
) *ReadHandler {
	return &ReadHandler{
		ctx:      ctx,
		services: services,
	}
}

func ReadHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewReadHandler(ctx, services)).Handle().JSON()
}

func (h *ReadHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *ReadHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *ReadHandler) Code() int                            { return h.code }
func (h *ReadHandler) Error() error                         { return h.err }
func (h *ReadHandler) Data() any                            { return h.result }

func (h *ReadHandler) JSON() error {
	if h.err != nil {
		return responses.NewTasksResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewTasksResponse().Successful(h.ctx, h.result)
}

func (h *ReadHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	if h.result, err = h.services.TaskService.GetAll(userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}
