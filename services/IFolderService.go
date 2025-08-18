package services

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)


// IFolderService interface
//
// This interface defines the methods that a folder service should implement.
// It is used to define the contract between the the application and the folder database table.
type IFolderService interface {
	GetAll() ([]repository.Folder, error)
	Get(id uuid.UUID) (*repository.Folder, error)
	GetByUser(user_id uuid.UUID) ([]repository.Folder, error)
	GetByParent(parent_id uuid.UUID) ([]repository.Folder, error)
	Create(params *repository.CreateFolderParams) (*repository.Folder, error)
	Update(params *repository.UpdateFolderParams) (*repository.Folder, error)
	Remove(id uuid.UUID) error
	Move(id uuid.UUID, parent_id *uuid.UUID) error
	Delete(id uuid.UUID) error
}

