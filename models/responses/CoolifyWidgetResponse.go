package responses

import "github.com/labstack/echo/v4"

type CoolifyWidgetData struct {
}

type CoolifyWidgetResponse struct {
	Data    *CoolifyWidgetData `json:"data"`
	Message string                    `json:"message"`
	Success bool                      `json:"success"`
} // @name CoolifyWidgetResponse


func NewCoolifyWidgetResponse(
) *CoolifyWidgetResponse {
	return &CoolifyWidgetResponse{
		Data:    &CoolifyWidgetData{},
		Success: true,
		Message: "Ok",
	}
}

func (w *CoolifyWidgetResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetResponse) Successful(
	ctx echo.Context,
) error {
	w.Success = true
	w.Data = nil 
	return ctx.JSON(200, w)
}

