package controllers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
	handlers "github.com/adamkali/mindscape/models/handlers/widget_handlers"
)

func (uc WidgetController) ControllerName() string {
	return uc.Name
}

func BuildWidgetController(p *Registrar) WidgetController {
	return WidgetController{
		Name:             "/widgets",
		Config:           p.Config,
		AuthService:      p.AuthService,
		MinioService:     p.MinioService,
		WidgetService:    p.WidgetService,
		RedisService:     p.RedisService,
		ValidatorService: p.ValidatorService,
	}
}

type WidgetController struct {
	Name             string
	Config           *configuration.Configuration
	AuthService      services.IAuthService
	WidgetService    services.IWidgetService
	MinioService     services.IMinioService
	RedisService     services.IRedisService
	ValidatorService *services.ValidatorService
}

// @Summary Get Widget Schemas
// @Description Get all embeded widget schemas from the Schema 
// @Description Storage 
//
// @ID          GetWidgetSchemas
// @Tags        Widgets
// @Produce     json
// @Success     200                 {object}     responses.WidgetsResponse
// @Failure     404                 {object}     responses.WidgetsResponse
// @Router      /widgets/schemas    [get]
func (wc WidgetController) ReadSchemas(
	ctx echo.Context,
) error {
	return handlers.NewReadHandler(
		ctx,
		wc.WidgetService,
	).Handle().JSON()
}

// @Summary Get Widget Schema by ID
// @Description Get embeded widget schema from the Schema
// @Description Storage by its identifier
//
// @ID          GetWidgetSchemaByID
// @Tags        Widgets
// @Produce     json
// @Param       schema_id           path         string                         true "Widget Schema Id"
// @Success     200                 {object}     responses.WidgetResponse
// @Failure     400                 {object}     responses.WidgetResponse
// @Failure     404                 {object}     responses.WidgetResponse
// @Router      /widgets/schemas/{schema_id}    [get]
func (wc WidgetController) GetSchemaByID(
	ctx echo.Context,
) error {
	return handlers.NewGetByIdHandler(
		ctx,
		wc.WidgetService,
	).Handle().JSON()
}

func (wc WidgetController) Attatch(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	api := e.Group("/api" + wc.Name)
	api.GET("/schemas", wc.ReadSchemas)
	api.GET("/schemas/:schema_id", wc.GetSchemaByID)
}
