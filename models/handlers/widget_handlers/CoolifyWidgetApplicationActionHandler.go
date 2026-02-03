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

type CoolifyApplicationActionHandler struct {
	ctx      echo.Context
	code     int
	err      error
	widget   *repository.UserWidget
	action   string
	appUUID  string
	result   *responses.CoolifyActionData
}

func NewCoolifyApplicationActionHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
	action string,
	appUUID string,
) *CoolifyApplicationActionHandler {
	return &CoolifyApplicationActionHandler{
		ctx:     ctx,
		code:    200,
		err:     nil,
		widget:  widget,
		action:  action,
		appUUID: appUUID,
		result:  nil,
	}
}

func CoolifyApplicationStartJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
	appUUID string,
) error {
	handler := NewCoolifyApplicationActionHandler(ctx, widget, "start", appUUID)
	return handler.Handle().JSON()
}

func CoolifyApplicationStopJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
	appUUID string,
) error {
	handler := NewCoolifyApplicationActionHandler(ctx, widget, "stop", appUUID)
	return handler.Handle().JSON()
}

func CoolifyApplicationRestartJsonHandler(
	ctx echo.Context,
	widget *repository.UserWidget,
	appUUID string,
) error {
	handler := NewCoolifyApplicationActionHandler(ctx, widget, "restart", appUUID)
	return handler.Handle().JSON()
}

func (h *CoolifyApplicationActionHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *CoolifyApplicationActionHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *CoolifyApplicationActionHandler) Code() int    { return h.code }
func (h *CoolifyApplicationActionHandler) Error() error { return h.err }
func (h *CoolifyApplicationActionHandler) Data() any    { return h.result }

func (h *CoolifyApplicationActionHandler) Handle() handlers.IHandler {
	var config CoolifyWidgetConfig
	err := json.Unmarshal(h.widget.Config, &config)
	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	client := clients.NewCoolifyClient(config.PersonalAccessToken, config.BaseUrl)

	var deploymentId string

	switch h.action {
	case "start":
		deploymentId, err = client.StartCoolifyApplication(h.ctx, h.appUUID)
	case "stop":
		deploymentId, err = client.StopCoolifyApplication(h.ctx, h.appUUID)
	case "restart":
		deploymentId, err = client.RestartCoolifyApplication(h.ctx, h.appUUID)
	default:
		return handlers.Lock(h, 400, fmt.Errorf("invalid action: %s", h.action))
	}

	if err != nil {
		return handlers.Lock(h, 400, err)
	}

	h.result = &responses.CoolifyActionData{
		Action:       h.action,
		AppUUID:      h.appUUID,
		DeploymentId: deploymentId,
	}
	return h
}

func (h *CoolifyApplicationActionHandler) JSON() error {
	if h.err == nil {
		return responses.NewCoolifyActionResponse().Successful(h.ctx, h.result)
	}
	return responses.NewCoolifyActionResponse().Fail(h.ctx, h.code, h.err)
}
