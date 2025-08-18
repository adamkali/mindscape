package services

import (
	"context"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FolderService struct {
	ctx context.Context
	pool *pgxpool.Pool
}

func CreateFolderService(
	ctx context.Context,	
	pool *pgxpool.Pool,
) IFolderService {
	return &FolderService{
		ctx:  ctx,
		pool: pool,
	}
}

func (folderService FolderService) GetAll() ([]repository.Folder, error) {
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return []repository.Folder{}, err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	folders, err := repo.FindFoldersAll(folderService.ctx)
	if err != nil {
		return []repository.Folder{}, err
	}
	tx.Commit(folderService.ctx)
	return folders, nil
}

func (folderService FolderService) Get(id uuid.UUID) (*repository.Folder, error) {
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return nil, err	
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	folder, err := repo.FindFolderById(folderService.ctx, &id)
	if err != nil {
		return nil, err
	}
	tx.Commit(folderService.ctx)
	return &folder, nil
}

func (folderService FolderService) GetByUser(user_id uuid.UUID) ([]repository.Folder, error) {
	if user_id != uuid.Nil {
		return []repository.Folder{}, nil
	}
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return []repository.Folder{}, err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	folders, err := repo.FindFoldersByUserId(folderService.ctx, &user_id)
	if err != nil {
		return []repository.Folder{}, err
	}
	tx.Commit(folderService.ctx)
	return folders, nil
}

// Misnomer
func (folderService FolderService) GetByParent(parent_id uuid.UUID) ([]repository.Folder, error) {
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return []repository.Folder{}, err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	
	var p_id_bytes []byte
	parent_id_type := &pgtype.UUID{}
	if p_id_bytes, err = parent_id.MarshalBinary(); err != nil {
		return []repository.Folder{}, err
	}

	if err = parent_id_type.UnmarshalJSON(p_id_bytes); err != nil {
		return []repository.Folder{}, err
	}

	folders, err := repo.FindFoldersByParentId(folderService.ctx, *parent_id_type)
	if err != nil {
		return []repository.Folder{}, err
	}
	tx.Commit(folderService.ctx)
	return folders, nil
}

func (folderService FolderService) Create(params *repository.CreateFolderParams) (*repository.Folder, error) {
	if params == nil {
		return nil, nil
	}
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)

	folder, err := repo.CreateFolder(folderService.ctx, *params)
	if err != nil {
		return nil, err
	}
	tx.Commit(folderService.ctx)
	return &folder, nil
}

func (folderService FolderService) Update(params *repository.UpdateFolderParams) (*repository.Folder, error) {
	if params == nil {
		return nil, nil
	}
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	err = repo.UpdateFolder(folderService.ctx, *params)
	if err != nil {
		return nil, err
	}

	folder, err := repo.FindFolderById(folderService.ctx, params.ID)
	if err != nil {
		return nil, err
	}
	tx.Commit(folderService.ctx)
	return &folder, nil
}

func (folderService FolderService) Remove(id uuid.UUID) error {
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteFolder(folderService.ctx, &id); err != nil {
		return err
	}
	tx.Commit(folderService.ctx)
	return nil
}

func (folderService FolderService) Move(id uuid.UUID, parent_id *uuid.UUID) error {
	if parent_id == nil {
		return nil
	}
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)

	var p_id_bytes []byte
	parent_id_type := &pgtype.UUID{}
	if p_id_bytes, err = parent_id.MarshalBinary(); err != nil {
		return err
	}
	if err = parent_id_type.UnmarshalJSON(p_id_bytes); err != nil {
		return err
	}
	if err := repo.MoveFolder(folderService.ctx, 
		repository.MoveFolderParams{
			ID: &id,
			ParentID: *parent_id_type,
		}); err != nil {
		return err
	}
	tx.Commit(folderService.ctx)
	return nil
}

func (folderService FolderService) Delete(id uuid.UUID) error {
	tx, err := folderService.pool.Begin(folderService.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(folderService.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteFolder(folderService.ctx, &id); err != nil {
		return err
	}
	tx.Commit(folderService.ctx)
	return nil
}
