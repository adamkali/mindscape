package widget_handlers

import (
	"encoding/json"

	"github.com/adamkali/mindscape/clients"
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetApplicationHandler struct {
	ctx    echo.Context
	code   int
	err    error
	widget *repository.UserWidget
	data   []clients.CoolifyApplication
}

func NewCoolifyWidgetApplicationHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) *CoolifyWidgetApplicationHandler {
	return &CoolifyWidgetApplicationHandler{
		ctx:    ctx,
		code:   200,
		err:    nil,
		widget: widget,
		data:   []clients.CoolifyApplication{},
	}
}

func CoolifyWidgetApplicationsJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
) error {
	handler := NewCoolifyWidgetApplicationHandler(ctx, widget)
	return handler.Handle().JSON()
}



func (h *CoolifyWidgetApplicationHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}
func (h *CoolifyWidgetApplicationHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *CoolifyWidgetApplicationHandler) Code() int    { return h.code }
func (h *CoolifyWidgetApplicationHandler) Error() error { return h.err }
func (h *CoolifyWidgetApplicationHandler) Data() any    { return h.data }
func (h *CoolifyWidgetApplicationHandler) Handle() handlers.IHandler {
	var config CoolifyWidgetConfig
	 err := json.Unmarshal(h.widget.Config, &config)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	client := clients.NewCoolifyClient(config.PersonalAccessToken, config.BaseUrl)
	h.data, err = client.GetCoolifyApplications(h.ctx)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	return h
}
func (h *CoolifyWidgetApplicationHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyWidgetApplicationResponse().Successful(h.ctx, h.data)
	} else {
		return responses.NewCoolifyWidgetApplicationResponse().Fail(h.ctx, h.code, h.err)
	}
}
