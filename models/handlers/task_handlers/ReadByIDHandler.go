package task_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReadByIDHandler struct {
	result   services.TaskDTO
	err      error
	code     int
	ctx      echo.Context
	services services.Registrar
}

func NewReadByIDHandler(
	ctx echo.Context,
	services services.Registrar,
) *ReadByIDHandler {
	return &ReadByIDHandler{
		ctx:      ctx,
		services: services,
	}
}

func ReadByIDHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewReadByIDHandler(ctx, services)).Handle().JSON()
}

func (h *ReadByIDHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *ReadByIDHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *ReadByIDHandler) Code() int                            { return h.code }
func (h *ReadByIDHandler) Error() error                         { return h.err }
func (h *ReadByIDHandler) Data() any                            { return h.result }

func (h *ReadByIDHandler) JSON() error {
	if h.err != nil {
		return responses.NewTaskResponse().Fail(h.ctx, h.code, h.err) 
	} else {
		return responses.NewTaskResponse().Successful(h.ctx, &h.result)
	}
}

func (h *ReadByIDHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err = h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	taskId := h.ctx.Param("taskId")
	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	if h.result, err = h.services.TaskService.GetById(parsedTaskId); err != nil {
		return handlers.Lock(h, 500, err)
	}
	if h.result.UserID != userID {
		return handlers.Lock(h, 403, fmt.Errorf("forbidden: task does not belong to user"))
	}
	return h
}
