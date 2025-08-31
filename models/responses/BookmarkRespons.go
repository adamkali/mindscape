package responses

import (
	"github.com/adamkali/mindscape/db/repository"
)

type BookmarkResponse struct {
	Data    *repository.Bookmark `json:"data"`
	Message string               `json:"message"`
	Success bool                 `json:"success"`
} // @name BookmarkResponse

func NewBookmarkResponse(
	data *repository.Bookmark,
	success bool,
	message string,
) *BookmarkResponse {
	return &BookmarkResponse{Data: data, Success: success, Message: message}
}
