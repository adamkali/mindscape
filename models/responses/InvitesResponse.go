package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type InvitesResponse struct {
	Data    []repository.GetInvitesByUserIDRow `json:"data"`
	Message string                             `json:"message"`
	Success bool                               `json:"success"`
} // @name InvitesResponse

func NewInvitesResponse() *InvitesResponse {
	return &InvitesResponse{Data: []repository.GetInvitesByUserIDRow{}, Success: false, Message: "Internal Server Error"}
}

func (b *InvitesResponse) Successful(ctx echo.Context, data []repository.GetInvitesByUserIDRow) error {
	b.Success = true
	b.Message = "OK"
	b.Data = data
	return ctx.JSON(200, b)
}

func (b *InvitesResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
