package responses

import (
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type FolderData struct {
	ID              *uuid.UUID            `json:"id"`
	ParentID        *uuid.UUID            `json:"parent_id"`
	UserID          *uuid.UUID            `json:"user_id"`
	Name            string                `json:"name"`
	Description     *string               `json:"description"`
	CreatedDatetime *time.Time            `json:"created_datetime"`
	UpdatedDatetime *time.Time            `json:"updated_datetime"`
	Children        []repository.Folder   `json:"children"`
	Bookmarks       []repository.Bookmark `json:"bookmarks"`
	Notes           []repository.Note     `json:"notes"`
}

type FolderResponse struct {
	Data    FolderData `json:"data"`
	Success bool       `json:"success"`
	Message string     `json:"message"`
}

func NewFolderResponse() *FolderResponse {
	return &FolderResponse{Success: false, Message: ""}
}

func NewFolderResponseWithData(data FolderData, success bool, message string) *FolderResponse {
	folderResponse := NewFolderResponse()
	folderResponse.Data = data
	folderResponse.Success = success
	folderResponse.Message = message
	return folderResponse
}

func NewFolderData(entity repository.Folder) FolderData {
	parentIDUUID := uuid.Nil
	if entity.ParentID.Valid {
		uuidbytes, _ := entity.ParentID.MarshalJSON()
		if err := parentIDUUID.UnmarshalBinary(uuidbytes); err != nil {
			return FolderData{}
		}
	}

	return FolderData{
		ID:              &entity.ID,
		ParentID:        &parentIDUUID,
		UserID:          &entity.UserID,
		Name:            entity.Name,
		Description:     entity.Description,
		CreatedDatetime: entity.CreatedDatetime,
		UpdatedDatetime: entity.UpdatedDatetime,
		Children:        []repository.Folder{},
		Bookmarks:       []repository.Bookmark{},
		Notes:           []repository.Note{},
	}
}
