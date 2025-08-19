package services

import (
	"context"
	"errors"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteService struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func CreateNoteService(
	ctx context.Context,
	pool *pgxpool.Pool,
) INoteService {
	return &NoteService{
		ctx:  ctx,
		pool: pool,
	}
}

func (noteService NoteService) GetAll() ([]repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	notes, err := repo.FindNotesAll(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	tx.Commit(noteService.ctx)
	return notes, nil
}

func (noteService NoteService) Get(id uuid.UUID) (*repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	note, err := repo.FindNoteById(noteService.ctx, &id)
	if err != nil {
		return nil, err
	}
	tx.Commit(noteService.ctx)
	return &note, nil
}

func (noteService NoteService) GetMostRecent(user_id uuid.UUID) (*repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	note, err := repo.FindNoteByUserIdMostRecent(noteService.ctx, &user_id)
	if err != nil {
		return nil, err
	}
	tx.Commit(noteService.ctx)
	return &note, nil
}

func (noteService NoteService) GetMostRecents(user_id uuid.UUID) ([]repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	notes, err := repo.FindNotesByUserIdMostRecent(noteService.ctx, &user_id)
	if err != nil {
		return []repository.Note{}, err
	}
	tx.Commit(noteService.ctx)
	return notes, nil
}
func (noteService NoteService) GetByUser(user_id uuid.UUID) ([]repository.Note, error) {
	if user_id == uuid.Nil {
		return []repository.Note{}, errors.New("user_id must not be nil")
	}
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	notes, err := repo.FindNotesByUserId(noteService.ctx, &user_id)
	if err != nil {
		return []repository.Note{}, err
	}
	tx.Commit(noteService.ctx)
	return notes, nil
}

func (noteService NoteService) GetByFolder(folder_id uuid.UUID) ([]repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	notes, err := repo.FindNotesByFolderId(noteService.ctx, &folder_id)
	if err != nil {
		return []repository.Note{}, err
	}
	tx.Commit(noteService.ctx)
	return notes, nil
}
func (noteService NoteService) GetByDateRange(params *repository.FindNotesByUserIDDateTimeRangeParams) ([]repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return []repository.Note{}, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	notes, err := repo.FindNotesByUserIDDateTimeRange(noteService.ctx, *params)
	if err != nil {
		return []repository.Note{}, err
	}
	tx.Commit(noteService.ctx)
	return notes, nil
}

func (noteService NoteService) Create(params *repository.CreateNoteParams) (*repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	note, err := repo.CreateNote(noteService.ctx, *params)
	if err != nil {
		return nil, err
	}
	tx.Commit(noteService.ctx)
	return &note, nil
}

func (noteService NoteService) Update(id uuid.UUID, params *repository.UpdateNoteParams) (*repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	err = repo.UpdateNote(noteService.ctx, *params)
	if err != nil {
		return nil, err
	}

	note, err := repo.FindNoteById(noteService.ctx, &id)
	if err != nil {
		return nil, err
	}
	tx.Commit(noteService.ctx)
	return &note, nil
}

func (noteService NoteService) Delete(id uuid.UUID) error {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteNote(noteService.ctx, &id); err != nil {
		return err
	}
	tx.Commit(noteService.ctx)
	return nil
}

func (noteService NoteService) Move(params *repository.MoveNoteParams) (*repository.Note, error) {
	tx, err := noteService.pool.Begin(noteService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(noteService.ctx)
	repo := repository.New(tx)
	err = repo.MoveNote(noteService.ctx, *params)
	if err != nil {
		return nil, err
	}

	note, err := repo.FindNoteById(noteService.ctx, params.ID)
	if err != nil {
		return nil, err
	}
	tx.Commit(noteService.ctx)
	return &note, nil
}
