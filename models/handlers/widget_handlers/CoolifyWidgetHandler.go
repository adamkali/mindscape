package widget_handlers

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetHandler struct {
	ctx         echo.Context
	code        int
	err         error
	widget      *repository.UserWidget
	applicatons []clients.CoolifyApplication
	services    []clients.CoolifyService
}

func NewCoolifyWidgetHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetHandler {
	return &CoolifyWidgetHandler{
		ctx:         ctx,
		code:        200,
		err:         nil,
		widget:      widget,
		applicatons: []clients.CoolifyApplication{},
		services:    []clients.CoolifyService{},
	}
}

func CoolifyWidgetJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewCoolifyWidgetHandler(ctx, widget)
	return handler.Handle().JSON()
}

type CoolifyWidgetConfig struct {
	BaseUrl             string   `json:"baseUrl"`
	PersonalAccessToken string   `json:"personalAccessToken"`
	ApplicationIds      []string `json:"applicationIds"`
	ServiceIds          []string `json:"serviceIds"`
}

func (h *CoolifyWidgetHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *CoolifyWidgetHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *CoolifyWidgetHandler) Code() int                            { return h.code }
func (h *CoolifyWidgetHandler) Error() error                         { return h.err }

func (h *CoolifyWidgetHandler) Data() any {
	return struct {
		Applicatons []clients.CoolifyApplication
		Services    []clients.CoolifyService
	}{
		Applicatons: h.applicatons,
		Services:    h.services,
	}
}

func (h *CoolifyWidgetHandler) Handle() handlers.IHandler {
	widgetConfig := h.widget.Config
	if widgetConfig == nil {
		h.SetError(nil)
	}

	return h
}

func (h *CoolifyWidgetHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetResponse().Successful(h.ctx, h.applicatons, h.services)
	} else {
		return responses.NewCoolifyWidgetResponse().Fail(h.ctx, h.code, h.err)
	}
}
