package widget_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ReadUserWidgetsHandler struct {
	ctx           echo.Context
	code          int
	err           error
	data          []repository.UserWidget
	widgetService services.IWidgetService
	authService   services.IAuthService
}

func NewReadUserWidgetsHandler(
	ctx echo.Context,
	ws services.IWidgetService,
	as services.IAuthService,
) *ReadUserWidgetsHandler {
	return &ReadUserWidgetsHandler{
		ctx:           ctx,
		widgetService: ws,
		authService:   as,
	}
}

func (h *ReadUserWidgetsHandler) Handle() handlers.IHandler {
	jwt_token := h.ctx.Get("user").(*jwt.Token)
	claims := jwt_token.Claims.(*services.CustomJwt)
	err := h.authService.CheckToken(jwt_token.Raw)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	h.data, err = h.widgetService.GetUserWidgets(claims.UserId)
	if err != nil {
		return handlers.Lock(h, 404, err)
	}
	return h
}

func (h *ReadUserWidgetsHandler) JSON() error {
	if h.err == nil {
		return responses.NewUserWidgetsResponse().Successful(h.ctx, h.data) 
	} else {
		return responses.NewUserWidgetsResponse().Fail(h.ctx, h.code, h.err) 
	}
}

func (h *ReadUserWidgetsHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *ReadUserWidgetsHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *ReadUserWidgetsHandler) Code() int {
	return h.code
}

func (h *ReadUserWidgetsHandler) Error() error {
	return h.err
}

func (h *ReadUserWidgetsHandler) Data() any {
	return h.data
}
