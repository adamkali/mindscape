package services

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)
type INoteService interface {
	GetAll() ([]repository.Note, error)
	Get(id uuid.UUID) (*repository.Note, error)
	GetMostRecent(user_id uuid.UUID) (*repository.Note, error)
	GetMostRecents(user_id uuid.UUID) ([]repository.Note, error)
	GetByUser(user_id uuid.UUID) ([]repository.Note, error)
	GetByFolder(folder_id uuid.UUID) ([]repository.Note, error)
	GetByDateRange(params *repository.FindNotesByUserIDDateTimeRangeParams) ([]repository.Note, error)


	Create(params *repository.CreateNoteParams) (*repository.Note, error)
	Update(id uuid.UUID, params *repository.UpdateNoteParams) (*repository.Note, error)
	Delete(id uuid.UUID) error
	Move(params *repository.MoveNoteParams) (*repository.Note, error)
}
