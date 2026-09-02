package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

// UsersSearchResponse returns slim user records (id + username only) for
// autocomplete. It deliberately never carries email or other PII.
type UsersSearchResponse struct {
	Data    []repository.SearchUsersByUsernameRow `json:"data"`
	Message string                                `json:"message"`
	Success bool                                  `json:"success"`
} // @name UsersSearchResponse

func NewUsersSearchResponse() *UsersSearchResponse {
	return &UsersSearchResponse{Data: []repository.SearchUsersByUsernameRow{}, Success: false, Message: "Internal Server Error"}
}

func (b *UsersSearchResponse) Successful(ctx echo.Context, data []repository.SearchUsersByUsernameRow) error {
	b.Success = true
	b.Message = "OK"
	b.Data = data
	return ctx.JSON(200, b)
}

func (b *UsersSearchResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
