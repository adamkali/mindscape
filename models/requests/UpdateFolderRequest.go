package requests

import "github.com/google/uuid"

type UpdateFolderRequest struct {
	UserID      uuid.UUID `json:"userId"`
	FolderID    uuid.UUID `json:"folderId"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}
