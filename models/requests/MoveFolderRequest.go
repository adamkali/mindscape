package requests

import "github.com/google/uuid"

type MoveFolderRequest struct {
	UserID      uuid.UUID  `json:"userId"`
	FolderID    uuid.UUID  `json:"folderId"`
	NewParentID *uuid.UUID `json:"newParentId"`
} // MoveFolderRequest
