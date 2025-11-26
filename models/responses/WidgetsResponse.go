package responses

import (
	"github.com/adamkali/mindscape/schemas"
	"github.com/labstack/echo/v4"
)

type WidgetsResponse struct {
	Data    []WidgetData `json:"data"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

func NewWidgetsResponse() *WidgetsResponse {
	return &WidgetsResponse{
		Data:    []WidgetData{},
		Success: false,
		Message: "",
	}
}

func (w *WidgetsResponse) Fail(
	ctx echo.Context,
	code int,
	err error,
) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func CreateWidgetData(data []schemas.WidgetSchema) []WidgetData {
	widgets := make([]WidgetData, len(data))
	for i, val := range data {
		datum := WidgetData{}
		widgets[i] = *datum.From(val)
	}
	return widgets
}

func (w *WidgetsResponse) Successful(
	ctx echo.Context,
	data []schemas.WidgetSchema,
) error {
	w.Data = CreateWidgetData(data) 
	w.Success = true
	return ctx.JSON(200, w)
}

