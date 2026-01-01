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

// @Summary Get a Users Widgets
// @Description Get a Users Widgets by their auth token 
// @Description and return the list of widgets associated 
// @Description with the user account in the request params.
//
// @ID          GetUserWidgets
// @Tags        Widgets
// @Produce     json 
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Success     200                 {object}     responses.UserWidgetsResponse
// @Failure     401                 {object}     responses.UserWidgetsResponse
// @Failure     403                 {object}     responses.UserWidgetsResponse
// @Failure     500                 {object}     responses.UserWidgetsResponse
// @Router      /widgets            [get]
func (wc WidgetController) Read(ctx echo.Context) error {
	return handlers.NewReadUserWidgetsHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

// @Summary Get a Users Widget
// @Description Get a Users Widget by their auth token 
// @Description and a path parameter and return the 
// @Description widget from the database.
//
// @ID          GetUserWidget
// @Tags        Widgets
// @Produce     json
// @Param       Authorization       header       string                         true "auth header"     default(Bearer token)
// @Param       user_widget_id      path         string                         true "Widget Id"       default("e38e78a4-2ca3-4c59-a3ea-a2019866e593")
// @Success     200                 {object}     responses.UserWidgetResponse
// @Failure     401                 {object}     responses.UserWidgetResponse
// @Failure     403                 {object}     responses.UserWidgetResponse
// @Failure     500                 {object}     responses.UserWidgetResponse
// @Router      /widgets/{user_widget_id}    [get]
func (wc WidgetController) ReadById(ctx echo.Context) error {
	return handlers.NewReadUserWidgetHandler(
		ctx,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

func (wc WidgetController) AddWidget(ctx echo.Context) error {
	return handlers.NewAddWidgetHandler(
		ctx,
		*wc.ValidatorService,
		wc.WidgetService,
		wc.AuthService,
	).Handle().JSON()
}

func (wc WidgetController) Attatch(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	api := e.Group("/api" + wc.Name)
	api.GET("/schemas", wc.ReadSchemas)
	api.GET("/schemas/:schema_id", wc.GetSchemaByID)

	api.GET("", wc.Read, authMiddleware)
	api.GET("/:user_widget_id", wc.ReadById, authMiddleware)
	api.POST("", wc.AddWidget, authMiddleware)
}
