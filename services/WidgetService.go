package services

import (
	"context"

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

func (ws *WidgetService) GetUserWidgets(userID string) ([]any, error) {
	return []any{}, nil
}

func (ws *WidgetService) CreateWidget(
	userID,
	widgetSchemaID uuid.UUID,
) error {
	return nil
}

func (ws *WidgetService) UpdateWidget(
	userID,
	widgetID uuid.UUID,
	config schemas.WidgetSchema,
) error {
	return nil
}

func (ws *WidgetService) DeleteWidget(
	userID,
	widgetID uuid.UUID,
) error {
	return nil
}

func (ws *WidgetService) GetWidgetSchema(
	widgetID uuid.UUID,
) (*schemas.WidgetSchema, error) {
	return ws.wss.Get(widgetID)
}
