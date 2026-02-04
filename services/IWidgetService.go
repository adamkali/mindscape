package services

import (
	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/schemas"
	"github.com/google/uuid"
)

type IWidgetService interface {
	ReadStorage() []schemas.WidgetSchema
	GetUserWidgets(userID uuid.UUID) ([]repository.UserWidget, error)
	CreateWidget(widgetConfig *repository.CreateUserWidgetParams) (*repository.UserWidget, error)
	UpdateWidget(widgetConfig *repository.UpdateUserWidgetParams) (*repository.UserWidget, error)
	DeleteWidget(widgetID uuid.UUID) error
	GetWidgetSchema(widgetID uuid.UUID) (*schemas.WidgetSchema, error)
	GetUserWidget(widgetID uuid.UUID) (*repository.UserWidget, error)
}
