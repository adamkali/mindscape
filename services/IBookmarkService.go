package services

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type IBookmarkService interface {
	GetAll() ([]repository.Bookmark, error)
	Get(id uuid.UUID) (*repository.Bookmark, error)
	GetByFolder(folder_id uuid.UUID) ([]repository.Bookmark, error)
	GetByUser(user_id uuid.UUID) ([]repository.Bookmark, error)
	GetMostRecent(user_id uuid.UUID) (*repository.Bookmark, error)
	GetMostRecents(user_id uuid.UUID) ([]repository.Bookmark, error)
	GetByDateRange(params *repository.FindBookmarksByUserIDDateTimeRangeParams) ([]repository.Bookmark, error)
	Create(params *repository.CreateBookmarkParams) (*repository.Bookmark, error)
	Update(id uuid.UUID, params *repository.UpdateBookmarkParams) (*repository.Bookmark, error)
	Move(params *repository.MoveBookmarkParams) (*repository.Bookmark, error)
	Remove(id uuid.UUID) error
}
