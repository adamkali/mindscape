package requests

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type UpdateBookmarkRequest struct {
	UserID     uuid.UUID `json:"userId"`
	BookmarkID uuid.UUID `json:"bookmarkId"`
	Name       string    `json:"name"`
	Link       string    `json:"link"`
} // @name UpdateBookmarkRequest

func (r UpdateBookmarkRequest) Into(folderID uuid.UUID) *repository.UpdateBookmarkParams {
	return &repository.UpdateBookmarkParams{
		ID:       r.BookmarkID,
		FolderID: folderID,
		Name:     r.Name,
		Link:     r.Link,
	}
}
