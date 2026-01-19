package services

import (
	"context"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/schemas"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WidgetService struct {
	ctx  context.Context
	pool *pgxpool.Pool
	wss  *schemas.WidgetSchemaStorage
}

func CreateWidgetService(
	ctx context.Context,
	pool *pgxpool.Pool,
) *WidgetService {
	storage, err := schemas.EmbeddedScemas()
	if err != nil {
		panic(err)
	}
	return &WidgetService{
		ctx:  ctx,
		pool: pool,
		wss:  storage,
	}
}

func (ws *WidgetService) ReadStorage() []schemas.WidgetSchema {
	return ws.wss.GetAll()
}

func (ws *WidgetService) GetUserWidgets(userID uuid.UUID) ([]repository.UserWidget, error) {
	tx, err := ws.pool.Begin(ws.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ws.ctx)
	var widgets []repository.UserWidget
	repo := repository.New(tx)
	if widgets, err = repo.FindUserWidgetsByUserID(ws.ctx, userID); err != nil {
		return nil, err
	}
	tx.Commit(ws.ctx)
	return widgets, nil
}

func (ws *WidgetService) CreateWidget(
	widgetConfig *repository.CreateUserWidgetParams,
) (*repository.UserWidget, error) {
	// TODO(claude): Add validation
	tx, err := ws.pool.Begin(ws.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ws.ctx)
	var userWidget repository.UserWidget
	repo := repository.New(tx)
	if userWidget, err = repo.CreateUserWidget(ws.ctx, *widgetConfig); err != nil {
		return nil, err
	}
	tx.Commit(ws.ctx)
	return &userWidget, nil
}

func (ws *WidgetService) UpdateWidget(
	widgetConfig *repository.UpdateUserWidgetParams,
) (*repository.UserWidget, error) {
	// TODO(claude): Add validation
	tx, err := ws.pool.Begin(ws.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ws.ctx)
	var userWidget repository.UserWidget
	repo := repository.New(tx)
	if userWidget, err = repo.UpdateUserWidget(ws.ctx, *widgetConfig); err != nil {
		return nil, err
	}
	tx.Commit(ws.ctx)
	return &userWidget, nil
}

func (ws *WidgetService) DeleteWidget(
	widgetID uuid.UUID,
) error {
	tx, err := ws.pool.Begin(ws.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ws.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteUserWidget(ws.ctx, widgetID); err != nil {
		return err
	}
	tx.Commit(ws.ctx)
	return nil

}

func (ws *WidgetService) GetWidgetSchema(
	widgetID uuid.UUID,
) (*schemas.WidgetSchema, error) {
	return ws.wss.Get(widgetID)
}

func (ws *WidgetService) GetUserWidget(
	widgetID uuid.UUID,
) (*repository.UserWidget, error) {
	tx, err := ws.pool.Begin(ws.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ws.ctx)
	var userWidget repository.UserWidget
	repo := repository.New(tx)
	if userWidget, err = repo.FindUserWidgetByID(ws.ctx, widgetID); err != nil {
		return nil, err
	}
	tx.Commit(ws.ctx)
	return &userWidget, nil
}

