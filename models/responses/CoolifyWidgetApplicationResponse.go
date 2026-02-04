package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetApplicationsResponse struct {
	Data    []CoolifyWidgetApplication `json:"data"`
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
} // @name CoolifyWidgetApplicationResponse

func NewCoolifyWidgetApplicationResponse() *CoolifyWidgetApplicationsResponse {
	return &CoolifyWidgetApplicationsResponse{
		Data:    []CoolifyWidgetApplication{},
		Success: false,
		Message: "",
	}
}

func (w *CoolifyWidgetApplicationsResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetApplicationsResponse) Successful(ctx echo.Context, application []clients.CoolifyApplication) error {
	w.Success = true
	for _, val := range application {
		w.Data = append(w.Data, *newCoolifyWidgetApplication(val))
	}
	return ctx.JSON(200, w)
}
