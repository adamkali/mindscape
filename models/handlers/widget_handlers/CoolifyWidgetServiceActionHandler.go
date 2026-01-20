package widget_handlers

import (
	"encoding/json"
	"fmt"

	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetServiceHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   *clients.CoolifyService
}

func NewCoolifyWidgetServiceHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetServiceHandler {
	return &CoolifyWidgetServiceHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   &clients.CoolifyService{},
	}
}


func CoolifyWidgetServiceJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewCoolifyWidgetServiceHandler(ctx, widget)
	return handler.Handle().JSON()
}


func (h *CoolifyWidgetServiceHandler) SetCode(code int) handlers.IHandler { h.code = code; return h }
func (h *CoolifyWidgetServiceHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *CoolifyWidgetServiceHandler) Code() int    { return h.code }
func (h *CoolifyWidgetServiceHandler) Error() error { return h.err }
func (h *CoolifyWidgetServiceHandler) Data() any    { return h.data }
func (h *CoolifyWidgetServiceHandler) Handle() handlers.IHandler {
	serviceId := h.ctx.Param("service_id")
	if serviceId == "" {
		return handlers.Lock(h, 400, fmt.Errorf("service_id is required"))	
	}
	var config CoolifyWidgetConfig
	err := json.Unmarshal(h.widget.Config, &config)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	client := clients.NewCoolifyClient(config.PersonalAccessToken, config.BaseUrl)
	h.data, err = client.GetCoolifyService(h.ctx, serviceId)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	return h
}

func (h *CoolifyWidgetServiceHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetServiceResponse().Successful(h.ctx, *h.data)
	} else {
		return responses.NewCoolifyWidgetServiceResponse().Fail(h.ctx, h.code, h.err)
	}
}
