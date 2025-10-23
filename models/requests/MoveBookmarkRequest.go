package requests

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type MoveBookmarkRequest struct {
	UserID      uuid.UUID `json:"userId"`
	BookmarkID  uuid.UUID `json:"bookmarkId"`
	NewParentID uuid.UUID `json:"newParentId"` // The parent folder to move the bookmark to cannot be null
} // @name MoveBookmarkRequest

func (r MoveBookmarkRequest) Into() *repository.MoveBookmarkParams {
	return &repository.MoveBookmarkParams{
		ID:       r.BookmarkID,
		FolderID: r.NewParentID,
	}

}
