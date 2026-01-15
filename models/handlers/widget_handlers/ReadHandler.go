package widget_handlers

import (
	"fmt"

	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/schemas"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

type ReadHandler struct {
	ctx           echo.Context
	code          int
	err           error
	data          []schemas.WidgetSchema
	widgetService services.IWidgetService
}

func NewReadHandler(ctx echo.Context, ws services.IWidgetService) *ReadHandler {
	return &ReadHandler{
		ctx:           ctx,
		widgetService: ws,
	}
}

func (h *ReadHandler) Handle() handlers.IHandler {
	h.data = h.widgetService.ReadStorage()
	h.code = 200
	return h 
}

func (h *ReadHandler) JSON() error {
	if h.err == nil {
		return responses.NewWidgetsResponse().Successful(h.ctx, h.data)
	} else {
		return h.ctx.JSON(h.code, h.err)
	}
}

func (h *ReadHandler) SetError(err error) handlers.IHandler {
	h.err = err
	return h
}

func (h *ReadHandler) SetCode(code int) handlers.IHandler {
	h.code = code
	return h
}

func (h *ReadHandler) Code() int {
	return h.code
}

func (h *ReadHandler) Data() any {
	return h.data
}

func (h *ReadHandler) Error() error {
	return h.err
}


