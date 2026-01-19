package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetApplicationResponse struct {
	Data    CoolifyWidgetApplication `json:"data"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
} // @name CoolifyWidgetApplicationResponse

func NewCoolifyWidgetApplicationResponse() *CoolifyWidgetApplicationResponse {
	return &CoolifyWidgetApplicationResponse{
		Data:    CoolifyWidgetApplication{},
		Success: false,
		Message: "",
	}
}

func (w *CoolifyWidgetApplicationResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetApplicationResponse) Successful(ctx echo.Context, application clients.CoolifyApplication) error {
	w.Success = true
	w.Data = *newCoolifyWidgetApplication(application)
	return ctx.JSON(200, w)
}
