package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type HouseholdsResponse struct {
	Data    []repository.Household `json:"data"`
	Message string                 `json:"message"`
	Success bool                   `json:"success"`
} // @name HouseholdsResponse

func NewHouseholdsResponse() *HouseholdsResponse {
	return &HouseholdsResponse{Data: []repository.Household{}, Success: false, Message: "Internal Server Error"}
}

func (b *HouseholdsResponse) Successful(ctx echo.Context, data []repository.Household) error {
	b.Success = true
	b.Message = "OK"
	b.Data = data
	return ctx.JSON(200, b)
}

func (b *HouseholdsResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
