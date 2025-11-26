package services

import (
	"github.com/google/uuid"
	"github.com/adamkali/mindscape/schemas"
)


type IWidgetService interface {
	ReadStorage() []schemas.WidgetSchema 
	GetUserWidgets(userID string) ([]any, error)
	CreateWidget(userID, widgetSchemaID uuid.UUID) error
	UpdateWidget(userID, widgetID uuid.UUID, config schemas.WidgetSchema) error
	DeleteWidget(userID, widgetID uuid.UUID) error
	GetWidgetSchema(widgetID uuid.UUID) (*schemas.WidgetSchema, error)
}
