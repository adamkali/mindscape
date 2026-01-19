package widget_handlers

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/labstack/echo/v4"
)

type CoolifyData struct {
}

type CoolifyWidgetHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   CoolifyData
}

func NewCoolifyWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetHandler {
	return &CoolifyWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   CoolifyData{},
	}
}

func CoolifyWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := &CoolifyWidgetHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
	}
	return handler.Handle().JSON()
}

type CoolifyWidgetConfig struct {
	Username            string   `json:"username"`
	PersonalAccessToken string   `json:"personalAccessToken"`
	ProjectIds          []string `json:"projectIds"`
}

func (h *CoolifyWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *CoolifyWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *CoolifyWidgetHandler) Code() int                            { return h.code }
func (h *CoolifyWidgetHandler) Error() error                         { return h.err }
func (h *CoolifyWidgetHandler) Data() any                            { return h.data }

func (h *CoolifyWidgetHandler) Handle() handlers.IHandler {
	return h
}

func (h *CoolifyWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetResponse().Successful(h.ctx, h.data)
	} else {
		return responses.NewCoolifyWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
