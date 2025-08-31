package responses

import (
	"github.com/adamkali/mindscape/db/repository"
)

type BookmarksResponse struct {
	Data    []repository.Bookmark `json:"data"`
	Message string               `json:"message"`
	Success bool                 `json:"success"`
} // @name BookmarksResponse

func NewBookmarksResponse(data []repository.Bookmark, success bool, message string) *BookmarksResponse {
	return &BookmarksResponse{Data: data, Success: success, Message: message}
}
