
package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetServiceResponse struct {
	Data    CoolifyWidgetService `json:"data"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
} // @name CoolifyWidgetServiceResponse

func NewCoolifyWidgetServiceResponse() *CoolifyWidgetServiceResponse {
	return &CoolifyWidgetServiceResponse{
		Data:    CoolifyWidgetService{},
		Success: false,
		Message: "",
	}
}

func (w *CoolifyWidgetServiceResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetServiceResponse) Successful(ctx echo.Context, application clients.CoolifyService) error {
	w.Success = true
	w.Data = *newCoolifyWidgetService(application)
	return ctx.JSON(200, w)
}
