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

type GetTasksByQueueTypeHandler struct {
	result []services.TaskDTO
	err    error
	code   int
	ctx    echo.Context
	services services.Registrar
}

func NewGetTasksByQueueTypeHandler(
	ctx echo.Context,
	services services.Registrar,
) *GetTasksByQueueTypeHandler {
	return &GetTasksByQueueTypeHandler{
		ctx:      ctx,
		services: services,
	}
}

func GetTasksByQueueTypeHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewGetTasksByQueueTypeHandler(ctx, services)).Handle().JSON()
}


func (h *GetTasksByQueueTypeHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GetTasksByQueueTypeHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GetTasksByQueueTypeHandler) Code() int                            { return h.code }
func (h *GetTasksByQueueTypeHandler) Error() error                         { return h.err }
func (h *GetTasksByQueueTypeHandler) Data() any                            { return h.result }

func (h *GetTasksByQueueTypeHandler) JSON() error {
	if h.err != nil {
		return responses.NewTasksResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewTasksResponse().Successful(h.ctx, h.result)
}

func (h *GetTasksByQueueTypeHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	if err := h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	queue := h.GetQueue()
	if queue == nil {
		return handlers.Lock(h, 400, fmt.Errorf("invalid queue type"))
	}
	if h.result, h.err = queue(userID); h.err != nil {
		return handlers.Lock(h, 500, h.err)
	}
	return h
}

func (h *GetTasksByQueueTypeHandler) GetQueue() func(uuid.UUID) ([]services.TaskDTO, error) {
	switch h.ctx.QueryParam("queueType") {
	case "":
		return h.services.TaskService.GetTasksByAvailableTaskType
	case "a":
		return h.services.TaskService.GetTasksByAvailableTaskType
	case "c":
		return h.services.TaskService.GetTasksByCompletedTaskType
	case "s":
		return h.services.TaskService.GetTasksByScheduledTaskType
	case "x":
		return h.services.TaskService.GetTasksByCancelledTaskType
	default:
		return nil
	}
}
