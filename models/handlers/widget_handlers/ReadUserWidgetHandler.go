package widget_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReadUserWidgetHandler struct {
	ctx           echo.Context
	code          int
	err           error
	data          *repository.UserWidget
	widgetService services.IWidgetService
	authService   services.IAuthService
}

func NewReadUserWidgetHandler(
	ctx echo.Context,
	ws services.IWidgetService,
	as services.IAuthService,
) *ReadUserWidgetHandler {
	return &ReadUserWidgetHandler{
		ctx:           ctx,
		widgetService: ws,
		authService:   as,
	}
}

func (h *ReadUserWidgetHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	err := h.authService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	userWidgetID, err := uuid.Parse(h.ctx.Param("user_widget_id"))
	h.data, err = h.widgetService.GetUserWidget(userWidgetID)
	if err != nil {
		return handlers.Lock(h, 404, err)
	}
	return h
}

func (h *ReadUserWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewUserWidgetResponse().Successful(h.ctx, h.data)
	} else {
		return responses.NewUserWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}

func (h *ReadUserWidgetHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *ReadUserWidgetHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *ReadUserWidgetHandler) Code() int {
	return h.code
}

func (h *ReadUserWidgetHandler) Error() error {
	return h.err
}

func (h *ReadUserWidgetHandler) Data() any {
	return h.data
}
