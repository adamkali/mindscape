package task_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GetTasksByTaskTypeHandler struct {
	result   []services.TaskDTO
	ctx      echo.Context
	services services.Registrar
	err      error
	code     int
}

func NewGetTasksByTaskTypeHandler(ctx echo.Context, services services.Registrar) *GetTasksByTaskTypeHandler {
	return &GetTasksByTaskTypeHandler{ctx: ctx, services: services}
}

func GetTasksByTaskTypeJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewGetTasksByTaskTypeHandler(ctx, services)).Handle().JSON()
}

func (h *GetTasksByTaskTypeHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *GetTasksByTaskTypeHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *GetTasksByTaskTypeHandler) Code() int                            { return h.code }
func (h *GetTasksByTaskTypeHandler) Error() error                         { return h.err }
func (h *GetTasksByTaskTypeHandler) Data() any                            { return h.result }

func (h *GetTasksByTaskTypeHandler) JSON() error {
	if h.err != nil {
		return responses.NewTasksResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewTasksResponse().Successful(h.ctx, h.result)
}

func (h *GetTasksByTaskTypeHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	if err := h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	taskTypeId := h.taskId(h.ctx.QueryParam("taskId"))
	if taskTypeId == uuid.Nil {
		return handlers.Lock(h, 400, nil)
	}
	if h.result, h.err = h.services.TaskService.GetTasksByTaskType(repository.GetTasksByTaskTypeParams{
		UserID:     userID,
		TaskTypeID: taskTypeId,
	}); h.err != nil {
		return handlers.Lock(h, 500, h.err)
	}
	return h

}

func (h *GetTasksByTaskTypeHandler) taskId(key string) uuid.UUID {
	taskStatus := h.ctx.QueryParam("status")
	switch taskStatus {
	case "a":
		return uuid.MustParse("e56fd149-24de-4835-9dad-ae861a7c3155")
	case "c":
		return uuid.MustParse("07bae843-7049-449c-a23e-ab78a571d7ca")
	case "d":
		return uuid.MustParse("546d40a2-aebd-4c3e-b1b3-3fd835211c74")
	case "h":
		return uuid.MustParse("56fabcc6-9703-43b5-96fc-2876646a26b9")
	case "p":
		return uuid.MustParse("11360cdc-f811-425f-b565-8b014c45ec25")
	case "r":
		return uuid.MustParse("f1559502-fa64-419b-9b55-57842e1af279")
	case "u":
		return uuid.MustParse("99dee5b2-7ac9-4b02-a3e5-a1c917d90009")
	case "i":
		return uuid.MustParse("106e703a-4dd4-4737-b38b-e4a0000ff158")
	default:
		return uuid.Nil
	}
}
