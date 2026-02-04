package responses

import (
	"github.com/adamkali/mindscape/schemas"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)


type WidgetData struct {
	ID         uuid.UUID                    `json:"id"`
	Type       string                       `json:"type"`
	Title      string                       `json:"title"`
	Layout     schemas.WidgetLayout         `json:"layout"`
	Properties map[string]schemas.WidgetProperty `json:"properties"`
	Required   []string                     `json:"required"`
}

type WidgetResponse struct {
	Data    *WidgetData `json:"data"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
} // @name WidgetResponse

func (w *WidgetData) From(schema schemas.WidgetSchema) *WidgetData {
	w.ID = schema.ID
	w.Type = schema.Type
	w.Title = schema.Title
	w.Layout = schema.Layout
	w.Properties = schema.Properties
	w.Required = schema.Required
	return w
}

func NewWidgetResponse() *WidgetResponse {
	return &WidgetResponse{
		Data:    &WidgetData{},
		Success: false,
		Message: "",
	}
}

func (w *WidgetResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *WidgetResponse) Successful(ctx echo.Context, schema schemas.WidgetSchema) error {
	w.Data = w.Data.From(schema)
	w.Success = true
	return ctx.JSON(200, w)
}
