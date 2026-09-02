package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type MembersResponse struct {
	Data    []repository.GetMembersByHouseholdIDRow `json:"data"`
	Message string                                  `json:"message"`
	Success bool                                    `json:"success"`
} // @name MembersResponse

func NewMembersResponse() *MembersResponse {
	return &MembersResponse{Data: []repository.GetMembersByHouseholdIDRow{}, Success: false, Message: "Internal Server Error"}
}

func (b *MembersResponse) Successful(ctx echo.Context, data []repository.GetMembersByHouseholdIDRow) error {
	b.Success = true
	b.Message = "OK"
	b.Data = data
	return ctx.JSON(200, b)
}

func (b *MembersResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
