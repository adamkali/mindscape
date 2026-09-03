package responses

import (
	"github.com/labstack/echo/v4"
)

type BoolResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
} // @name BoolResponse

func NewBoolResponse() *BoolResponse {
	return &BoolResponse{Success: false, Message: "Internal Server Error"}
}

func (b *BoolResponse) Successful(ctx echo.Context) error {
	b.Success = true
	b.Message = "OK"
	return ctx.JSON(200, b)
}

func (b *BoolResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
