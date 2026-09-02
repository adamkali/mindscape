package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type HouseholdResponse struct {
	Data    repository.Household `json:"data"`
	Message string               `json:"message"`
	Success bool                 `json:"success"`
} // @name HouseholdResponse

func NewHouseholdResponse() *HouseholdResponse {
	return &HouseholdResponse{Data: repository.Household{}, Success: false, Message: "Internal Server Error"}
}

func (b *HouseholdResponse) Successful(ctx echo.Context, data repository.Household) error {
	b.Success = true
	b.Message = "OK"
	b.Data = data
	return ctx.JSON(200, b)
}

func (b *HouseholdResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
