package task_handlers

import (
	"fmt"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	handlers "github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type UpdateTaskStatusHandler struct {
	result   services.TaskDTO
	ctx      echo.Context
	services services.Registrar
	err      error
	code     int
}

func NewUpdateTaskStatusHandler(
	ctx echo.Context,
	services services.Registrar,
) *UpdateTaskStatusHandler {
	return &UpdateTaskStatusHandler{
		ctx:      ctx,
		services: services,
	}
}

func UpdateTaskStatusHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewUpdateTaskStatusHandler(ctx, services)).Handle().JSON()
}

func (h *UpdateTaskStatusHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	var err error
	if err = h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	// get task id
	taskId := h.ctx.Param("taskId")
	parsedTaskId, err := uuid.Parse(taskId)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	// get task status
	taskStatus := h.ctx.QueryParam("status")
	taskStatusDate := h.ctx.QueryParam("due")
	dueAt, err := time.Parse(taskStatusDate, "2006-01-02") 
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	switch taskStatus {
	case "a":
	    h.result, err = h.services.TaskService.UpdateAsAmbiguous(parsedTaskId)
		break
	case "c":
	    h.result, err = h.services.TaskService.UpdateAsCancelled(parsedTaskId)
		break
	case "d":
		h.result, err = h.services.TaskService.UpdateAsDone(parsedTaskId)
		break
	case "h":
		h.result, err = h.services.TaskService.UpdateAsHold(parsedTaskId)
		break
	case "p":
		h.result, err = h.services.TaskService.UpdateAsPending(repository.UpdateAsPendingParams{
			ID: parsedTaskId,
			DueAt:    pgtype.Timestamptz{ Time: dueAt, Valid: true },
		})
		break
	case "r":
		h.result, err = h.services.TaskService.UpdateAsRecurring(parsedTaskId)
		break
	case "u":
		h.result, err = h.services.TaskService.UpdateAsUndone(repository.UpdateAsUndoneParams{
			ID: parsedTaskId,
			DueAt:    pgtype.Timestamptz{ Time: dueAt, Valid: true },
		})
		break
	case "i":
		h.result, err = h.services.TaskService.UpdateAsUrgent(parsedTaskId)
		break
	default:
		handlers.Lock(h, 400, fmt.Errorf("invalid status"))
	}
	return h
}

func (h *UpdateTaskStatusHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *UpdateTaskStatusHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *UpdateTaskStatusHandler) Code() int                            { return h.code }
func (h *UpdateTaskStatusHandler) Error() error                         { return h.err }
func (h *UpdateTaskStatusHandler) Data() any                            { return h.result }

func (h *UpdateTaskStatusHandler) JSON() error {
	if h.err != nil {
		return responses.NewTaskResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewTaskResponse().Successful(h.ctx, &h.result)
}

