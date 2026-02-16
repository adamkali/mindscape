package task_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/requests"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UpdateTaskContentHandler struct {
	result services.TaskDTO
	ctx      echo.Context
	services services.Registrar
	err      error
	code     int
}

func NewUpdateTaskContentHandler(ctx echo.Context, services services.Registrar) *UpdateTaskContentHandler {
	return &UpdateTaskContentHandler{ctx: ctx, services: services}
}

func UpdateTaskContentHandlerJsonHandler(
	ctx echo.Context,
	services services.Registrar,
) error {
	return (NewUpdateTaskContentHandler(ctx, services)).Handle().JSON()
}

func (h *UpdateTaskContentHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *UpdateTaskContentHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *UpdateTaskContentHandler) Code() int                            { return h.code }
func (h *UpdateTaskContentHandler) Error() error                         { return h.err }
func (h *UpdateTaskContentHandler) Data() any                            { return h.result }

func (h *UpdateTaskContentHandler) JSON() error {
	if h.err != nil {
		return responses.NewTaskResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewTaskResponse().Successful(h.ctx, &h.result)
}

func (h *UpdateTaskContentHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	userID := claims.UserId
	var err error
	if err := h.services.AuthService.CheckToken(jwt_token.Raw); err != nil {
		return handlers.Lock(h, 401, err)
	}
	var request *requests.UpdateTaskContentRequest
	if request, err = h.services.ValidatorService.UpdateTaskContentRequest(h.ctx); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if h.result, err = h.services.TaskService.GetById(request.ID); err != nil {
		return handlers.Lock(h, 404, err)
	}
	if h.result.UserID != userID {
		return handlers.Lock(h, 403, err)
	}
	if h.result, err = h.services.TaskService.UpdateTaskContent(repository.UpdateTaskContentParams{
		ID: request.ID,
		Name: &request.Name,
		Description: &request.Description,
	}); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}
