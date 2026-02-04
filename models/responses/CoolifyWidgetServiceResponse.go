package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetServicesResponse struct {
	Data    []CoolifyWidgetService `json:"data"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
} // @name CoolifyWidgetServiceResponse

func NewCoolifyWidgetServicesResponse() *CoolifyWidgetServicesResponse {
	return &CoolifyWidgetServicesResponse{
		Data:    []CoolifyWidgetService{},
		Success: false,
		Message: "",
	}
}

func (w *CoolifyWidgetServicesResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetServicesResponse) Successful(ctx echo.Context, services []clients.CoolifyService) error {
	w.Success = true
	sers := make([]CoolifyWidgetService, len(services))
	for i, val := range services {
		sers[i] = *newCoolifyWidgetService(val)
	}
	w.Data = sers
	return ctx.JSON(200, w)
}
