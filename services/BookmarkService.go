package services

import (
	"context"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookmarkService struct {
	ctx context.Context
	pool *pgxpool.Pool
}

func CreateBookmarkService(
	ctx context.Context,	
	pool *pgxpool.Pool,
) IBookmarkService {
	return &BookmarkService{
		ctx:  ctx,
		pool: pool,
	}
}


func (bookmarkService BookmarkService) GetAll() ([]repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmarks, err := repo.FindBookmarksAll(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	tx.Commit(bookmarkService.ctx)
	return bookmarks, nil

}
func (bookmarkService BookmarkService) Get(id uuid.UUID) (*repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmark, err := repo.FindBookmarkById(bookmarkService.ctx, id)
	if err != nil {
		return nil, err
	}
	tx.Commit(bookmarkService.ctx)
	return &bookmark, nil
}

func (bookmarkService BookmarkService) GetByFolder(folder_id uuid.UUID) ([]repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmarks, err := repo.FindBookmarksByFolderId(bookmarkService.ctx, folder_id)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	tx.Commit(bookmarkService.ctx)
	return bookmarks, nil
}

func (bookmarkService BookmarkService) GetByUser(user_id uuid.UUID) ([]repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmarks, err := repo.FindBookmarksByUserId(bookmarkService.ctx, user_id)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	tx.Commit(bookmarkService.ctx)
	return bookmarks, nil
}
func (bookmarkService BookmarkService) GetMostRecent(user_id uuid.UUID) (*repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmark, err := repo.FindBookmarkByUserIdMostRecent(bookmarkService.ctx, user_id)
	if err != nil {
		return nil, err
	}
	tx.Commit(bookmarkService.ctx)
	return &bookmark, nil
}
func (bookmarkService BookmarkService) GetMostRecents(user_id uuid.UUID) ([]repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmarks, err := repo.FindBookmarksByUserIdMostRecent(bookmarkService.ctx, user_id)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	tx.Commit(bookmarkService.ctx)
	return bookmarks, nil
}
func (bookmarkService BookmarkService) GetByDateRange(params *repository.FindBookmarksByUserIDDateTimeRangeParams) ([]repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmarks, err := repo.FindBookmarksByUserIDDateTimeRange(bookmarkService.ctx, *params)
	if err != nil {
		return []repository.Bookmark{}, err
	}
	tx.Commit(bookmarkService.ctx)
	return bookmarks, nil
}
func (bookmarkService BookmarkService) Create(params *repository.CreateBookmarkParams) (*repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	bookmark, err := repo.CreateBookmark(bookmarkService.ctx, *params)
	if err != nil {
		return nil, err
	}
	tx.Commit(bookmarkService.ctx)
	return &bookmark, nil
}
func (bookmarkService BookmarkService) Update(id uuid.UUID, params *repository.UpdateBookmarkParams) (*repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	err = repo.UpdateBookmark(bookmarkService.ctx, *params)
	if err != nil {
		return nil, err
	}

	bookmark, err := repo.FindBookmarkById(bookmarkService.ctx, id)
	if err != nil {
		return nil, err
	}
	tx.Commit(bookmarkService.ctx)
	return &bookmark, nil
}
func (bookmarkService BookmarkService) Move(params *repository.MoveBookmarkParams) (*repository.Bookmark, error){
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	err = repo.MoveBookmark(bookmarkService.ctx, *params)
	if err != nil {
		return nil, err
	}
	tx.Commit(bookmarkService.ctx)

	bookmark, err := repo.FindBookmarkById(bookmarkService.ctx, params.ID)
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}
func (bookmarkService BookmarkService) Remove(id uuid.UUID) error {
	tx, err := bookmarkService.pool.Begin(bookmarkService.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(bookmarkService.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteBookmark(bookmarkService.ctx, id); err != nil {
		return err
	}
	tx.Commit(bookmarkService.ctx)
	return nil
}
