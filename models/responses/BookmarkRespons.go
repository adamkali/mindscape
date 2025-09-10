package responses

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/labstack/echo/v4"
)

type BookmarkResponse struct {
	Data    *repository.Bookmark `json:"data"`
	Message string               `json:"message"`
	Success bool                 `json:"success"`
} // @name BookmarkResponse

func NewBookmarkResponse() *BookmarkResponse {
	return &BookmarkResponse{Data: nil, Success: false, Message: "Internal Server Error"}
}


func (b *BookmarkResponse) Successful(ctx echo.Context, data repository.Bookmark) error {
	b.Success = true
	b.Message = "OK"
	b.Data = &data
	return ctx.JSON(200, b)
}

func (b *BookmarkResponse) Fail(ctx echo.Context, code int, err error) error {
	b.Success = false
	b.Message = err.Error()
	return ctx.JSON(code, b)
}
