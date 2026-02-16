package task_handlers

import (
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DeleteHandler struct {
	result   string 
	err      error
	code     int
	ctx      echo.Context
	services services.Registrar
}

func NewDeleteHandler(
	ctx echo.Context,
	services services.Registrar,
) *DeleteHandler {
	return &DeleteHandler{
		ctx:      ctx,
		services: services,
	}
}

func DeleteHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewDeleteHandler(ctx, services)).Handle().JSON()
}

func (h *DeleteHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *DeleteHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *DeleteHandler) Code() int                            { return h.code }
func (h *DeleteHandler) Error() error                         { return h.err }
func (h *DeleteHandler) Data() any                            { return h.result }

func (h *DeleteHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, h.result)
}

func (h *DeleteHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	var err error
	if err = h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}

	taskId := h.ctx.Param("taskId")
	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	var task services.TaskDTO
	if task, err = h.services.TaskService.GetById(parsedTaskId); err != nil {
		return handlers.Lock(h, 500, err)
	}
	if task.UserID != claims.UserId {
		return handlers.Lock(h, 403, err)
	}

	if err = h.services.TaskService.Delete(task.ID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}
